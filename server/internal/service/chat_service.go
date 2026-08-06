package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ai-gateway-server/internal/model"
)

// ChatService 处理大模型调用 + 日志采集
type ChatService struct {
	DB            *gorm.DB
	RDB           *redis.Client
	DefaultQuota  int
	ModelProfiles map[string]ModelProfile
}

// ModelProfile 模型行为特征
type ModelProfile struct {
	MinCompletionTokens int     // 最小输出 token
	MaxCompletionTokens int     // 最大输出 token
	MinLatencyMs        int     // 最小延迟
	MaxLatencyMs        int     // 最大延迟
	SuccessRate         float64 // 成功率 0-1
}

// ChatRequest 调用请求
type ChatRequest struct {
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required,min=1"`
}

// Message 对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 调用响应（兼容 OpenAI 格式）
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 回复选项
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage Token 用量
type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	Cost             float64 `json:"cost"`
}

// ChatResult 完整调用结果（用于日志写入）
type ChatResult struct {
	Response         ChatResponse
	ProjectID        uint
	ApiKeyID         uint
	ModelID          uint
	ModelName        string
	InputPrice       float64
	OutputPrice      float64
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Cost             float64
	LatencyMs        int
	Status           string
	ErrorMessage     string
}

// NewChatService 创建 ChatService，初始化模型特征
func NewChatService(db *gorm.DB, rdb *redis.Client, quota int) *ChatService {
	return &ChatService{
		DB:           db,
		RDB:          rdb,
		DefaultQuota: quota,
		ModelProfiles: map[string]ModelProfile{
			"qwen2.5-72b":     {80, 500, 200, 800, 0.990},
			"deepseek-v3":     {100, 450, 150, 500, 0.985},
			"minimax-text-01": {60, 400, 250, 1200, 0.970},
			"glm-4":           {50, 350, 80, 400, 0.995},
			"qwen2.5-coder":   {80, 450, 150, 700, 0.988},
		},
	}
}

// Chat 处理一次大模型调用请求
func (s *ChatService) Chat(ctx context.Context, projectID, apiKeyID uint, req ChatRequest) (*ChatResult, error) {
	// 1. 查找模型
	var m model.Model
	if err := s.DB.Where("name = ?", req.Model).First(&m).Error; err != nil {
		return nil, fmt.Errorf("模型 %s 不存在或已下线", req.Model)
	}
	if m.Status == "offline" {
		return nil, fmt.Errorf("模型 %s 已下线", req.Model)
	}

	// 2. 估算 prompt_tokens
	promptTokens := estimatePromptTokens(req.Messages)

	// 3. 按模型特征生成数据
	profile, ok := s.ModelProfiles[m.Name]
	if !ok {
		// 默认特征
		profile = ModelProfile{80, 400, 100, 600, 0.98}
	}

	completionTokens := profile.MinCompletionTokens + rand.Intn(profile.MaxCompletionTokens-profile.MinCompletionTokens+1)
	latencyMs := profile.MinLatencyMs + rand.Intn(profile.MaxLatencyMs-profile.MinLatencyMs+1)

	// 4. 状态判定
	status := "success"
	errMsg := ""
	if rand.Float64() > profile.SuccessRate {
		r := rand.Float64()
		if r < 0.5 {
			status = "fail"
			errMsg = "internal server error"
		} else {
			status = "timeout"
			errMsg = "request timeout after 30s"
		}
	}

	// 5. 费用计算
	cost := (float64(promptTokens)/1000.0)*m.InputPrice + (float64(completionTokens)/1000.0)*m.OutputPrice
	cost = math.Round(cost*1000000) / 1000000 // 保留 6 位小数

	// 6. 生成回复
	content := generateMockContent(req.Messages)

	result := &ChatResult{
		Response: ChatResponse{
			ID:      fmt.Sprintf("chatcmpl-mock-%s", randomID(12)),
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   m.Name,
			Choices: []Choice{{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			}},
			Usage: Usage{
				PromptTokens:     promptTokens,
				CompletionTokens: completionTokens,
				TotalTokens:      promptTokens + completionTokens,
				Cost:             cost,
			},
		},
		ProjectID:        projectID,
		ApiKeyID:         apiKeyID,
		ModelID:          m.ID,
		ModelName:        m.Name,
		InputPrice:       m.InputPrice,
		OutputPrice:      m.OutputPrice,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Cost:             cost,
		LatencyMs:        latencyMs,
		Status:           status,
		ErrorMessage:     errMsg,
	}

	// 7. 同步写入 MySQL
	if err := s.collectMySQL(result); err != nil {
		log.Printf("[Collector] MySQL write error: %v", err)
	}

	// 8. 异步更新 Redis（失败不影响主流程）
	go s.collectRedis(result)

	return result, nil
}

// collectMySQL 写入 api_usage_logs
func (s *ChatService) collectMySQL(r *ChatResult) error {
	log := model.ApiUsageLog{
		ProjectID:        r.ProjectID,
		ApiKeyID:         r.ApiKeyID,
		ModelID:          r.ModelID,
		PromptTokens:     r.PromptTokens,
		CompletionTokens: r.CompletionTokens,
		TotalTokens:      r.TotalTokens,
		Cost:             r.Cost,
		LatencyMs:        r.LatencyMs,
		Status:           r.Status,
		ErrorMessage:     r.ErrorMessage,
	}
	return s.DB.Create(&log).Error
}

// collectRedis 异步更新 Redis 实时指标
func (s *ChatService) collectRedis(r *ChatResult) {
	ctx := context.Background()
	dateKey := time.Now().Format("20060102")
	ttl := 48 * time.Hour

	pipe := s.RDB.Pipeline()

	// 全局当日指标
	pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:total_tokens", dateKey), int64(r.TotalTokens))
	pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:total_requests", dateKey), 1)

	costCents := int64(math.Round(r.Cost * 1000000)) // 用整数存费用（微元）
	if r.Cost > 0 {
		pipe.IncrByFloat(ctx, fmt.Sprintf("stats:%s:total_cost", dateKey), r.Cost)
		_ = costCents // 避免未使用警告
	}
	pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:total_latency_ms", dateKey), int64(r.LatencyMs))

	switch r.Status {
	case "success":
		pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:success_count", dateKey), 1)
	case "fail":
		pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:fail_count", dateKey), 1)
	case "timeout":
		pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:timeout_count", dateKey), 1)
	}

	// 各模型当日指标
	pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:model:%d:tokens", dateKey, r.ModelID), int64(r.TotalTokens))
	pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:model:%d:requests", dateKey, r.ModelID), 1)
	if r.Cost > 0 {
		pipe.IncrByFloat(ctx, fmt.Sprintf("stats:%s:model:%d:cost", dateKey, r.ModelID), r.Cost)
	}
	pipe.IncrBy(ctx, fmt.Sprintf("stats:%s:model:%d:latency_ms", dateKey, r.ModelID), int64(r.LatencyMs))

	// 额度计数
	pipe.IncrBy(ctx, fmt.Sprintf("quota:key:%d:%s:tokens", r.ApiKeyID, dateKey), int64(r.TotalTokens))
	pipe.IncrBy(ctx, fmt.Sprintf("quota:key:%d:%s:requests", r.ApiKeyID, dateKey), 1)

	// 全局累计
	pipe.IncrBy(ctx, "stats:total:total_tokens", int64(r.TotalTokens))
	pipe.IncrBy(ctx, "stats:total:total_requests", 1)
	if r.Cost > 0 {
		pipe.IncrByFloat(ctx, "stats:total:total_cost", r.Cost)
	}

	// 设置过期
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:total_tokens", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:total_requests", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:total_cost", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:success_count", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:fail_count", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:timeout_count", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:total_latency_ms", dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:model:%d:tokens", dateKey, r.ModelID), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:model:%d:requests", dateKey, r.ModelID), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:model:%d:cost", dateKey, r.ModelID), ttl)
	pipe.Expire(ctx, fmt.Sprintf("stats:%s:model:%d:latency_ms", dateKey, r.ModelID), ttl)
	pipe.Expire(ctx, fmt.Sprintf("quota:key:%d:%s:tokens", r.ApiKeyID, dateKey), ttl)
	pipe.Expire(ctx, fmt.Sprintf("quota:key:%d:%s:requests", r.ApiKeyID, dateKey), ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("[Collector] Redis write error: %v", err)
	}
}

// estimatePromptTokens 基于消息长度估算 prompt tokens
// 中英文混合粗略估算: ~1.5 chars/token
func estimatePromptTokens(msgs []Message) int {
	total := 0
	for _, m := range msgs {
		chars := len([]rune(m.Content))
		total += int(float64(chars) / 1.5)
	}
	if total < 5 {
		total = 5
	}
	return total
}

// mockContentPool 模拟回复模板池
var mockContentPool = []string{
	"这是一个很好的问题。让我为你详细分析一下。\n\n首先，从技术角度来看，这个问题涉及到多个层面的思考。在当前的行业实践中，我们通常采用以下方法来解决这类问题：\n\n1. **需求分析**：明确问题的边界条件和核心目标\n2. **方案设计**：根据具体场景选择最合适的技术路径\n3. **实现验证**：通过实际测试来验证解决方案的有效性\n\n希望这个回答对你有所帮助！如果还有任何疑问，欢迎继续交流。",
	"关于这个问题，我想从几个方面来回答你。\n\n根据我的理解，这个问题的核心在于如何平衡性能和可维护性。在实践中，我们常常看到以下模式：\n\n- **模式一**：采用分层架构来隔离关注点\n- **模式二**：通过缓存机制来提升响应速度\n- **模式三**：利用异步处理来优化用户体验\n\n每个模式都有其适用场景，选择时需要结合实际业务需求来权衡。",
	"谢谢你提出这个问题。\n\n让我从实践的角度来分享一下我的理解。在真实的业务场景中，这类问题往往需要关注以下几点：\n\n**第一，数据一致性**：确保在任何情况下数据的准确性都是最优先的。\n\n**第二，系统可靠性**：通过合理的错误处理和重试机制来保证服务的稳定性。\n\n**第三，扩展性考量**：在设计之初就考虑到未来可能的增长需求。\n\n综上所述，解决这个问题需要综合运用多种策略，而不仅仅是单一的技术手段。",
	"非常好的问题！这是一个在业界经常被讨论的话题。\n\n基于我的知识积累和理解，我可以为你提供以下信息：\n\n这个问题涉及到的核心技术概念包括：\n\n1. **分布式系统**的基本原理\n2. **容错机制**的设计模式\n3. **性能优化**的最佳实践\n\n实际应用中，很多团队会选择从最简单的方案开始，然后根据反馈逐步优化。这种迭代方式虽然看起来不够「完美」，但往往是最务实的选择。",
	"我理解你的疑问。让我尝试用简洁清晰的方式来解答。\n\n这个问题其实可以从两个维度来看：\n\n**理论维度**：从计算机科学的基础理论出发，这个问题的本质是...\n\n**工程维度**：在实际开发中，我们更关注的是如何用最小的成本达到最好的效果。\n\n综合来看，理解和解决这个问题需要理论和实践的有机结合。如果你对某个具体方面感兴趣，我可以为你进一步展开。",
}

// generateMockContent 生成模拟回复
func generateMockContent(msgs []Message) string {
	if len(msgs) == 0 {
		return mockContentPool[rand.Intn(len(mockContentPool))]
	}

	// 根据用户最后一条消息的内容尝试做简单匹配
	lastMsg := msgs[len(msgs)-1]
	content := strings.ToLower(lastMsg.Content)

	for _, keyword := range []string{"代码", "code", "编程", "programming"} {
		if strings.Contains(content, keyword) {
			return "关于编程相关的问题，我可以为你提供一些思路。\n\n```python\n# 这是一个示例代码片段\ndef solve_problem(input_data):\n    result = []\n    for item in input_data:\n        processed = process(item)\n        result.append(processed)\n    return result\n\n# 使用示例\ndata = [1, 2, 3, 4, 5]\noutput = solve_problem(data)\nprint(output)\n```\n\n以上代码展示了解决这类问题的基本框架。在实际应用中，你可能需要根据具体需求进行调整和优化。如果有具体的业务场景，我可以给出更针对性的建议。"
		}
	}

	for _, keyword := range []string{"量子", "quantum"} {
		if strings.Contains(content, keyword) {
			return "量子计算是一个令人着迷的领域！\n\n量子计算利用量子力学的基本原理——**叠加态**和**纠缠**——来进行信息处理。与传统计算机使用比特（0或1）不同，量子计算机使用**量子比特（qubit）**，它可以同时处于多种状态的叠加。\n\n这种特性使得量子计算机在特定问题上具有指数级的加速优势，例如：\n\n1. **大数分解**（Shor算法）\n2. **搜索问题**（Grover算法）\n3. **量子化学模拟**\n\n目前量子计算仍处于早期阶段，但已经在金融、制药、材料科学等领域展现出巨大的应用潜力。"
		}
	}

	for _, keyword := range []string{"你好", "hello", "hi", "介绍", "你是谁"} {
		if strings.Contains(content, keyword) {
			return "你好！我是 AI 助手，很高兴为你服务。\n\n我可以帮助你解答各种问题，包括但不限于：\n\n- 技术问题的分析和解答\n- 编程相关的问题\n- 学术知识的探讨\n- 日常问题的建议\n\n请随时告诉我你需要什么帮助，我会尽力为你提供详细和准确的回答！"
		}
	}

	// 默认返回一条模板回复
	return mockContentPool[rand.Intn(len(mockContentPool))]
}

// randomID 生成随机 ID
func randomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
