package simulator

import (
	"context"
	"log"
	"math/rand"
	"time"

	"ai-gateway-server/internal/service"
)

// Simulator 后台模拟数据生成器
type Simulator struct {
	ctx    context.Context
	cancel context.CancelFunc
	svc    *service.ChatService

	// 默认使用种子数据中的项目 1 和 API Key 1
	defaultProjectID uint
	defaultApiKeyID  uint
}

// New 创建模拟器
func New(svc *service.ChatService) *Simulator {
	return &Simulator{
		svc:              svc,
		defaultProjectID: 1,
		defaultApiKeyID:  1,
	}
}

// Start 启动后台 goroutine
func (s *Simulator) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	go func() {
		log.Println("[Simulator] started")
		ticker := time.NewTicker(s.randomInterval())
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				log.Println("[Simulator] stopped")
				return
			case <-ticker.C:
				s.generateOne()
				ticker.Reset(s.randomInterval())
			}
		}
	}()
}

// Stop 停止模拟器
func (s *Simulator) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// generateOne 生成一次模拟调用
func (s *Simulator) generateOne() {
	// 按权重随机选择模型
	modelName := s.pickModel()

	// 构造模拟请求
	req := service.ChatRequest{
		Model: modelName,
		Messages: []service.Message{
			{Role: "user", Content: s.mockUserMsg()},
		},
	}

	_, err := s.svc.Chat(s.ctx, s.defaultProjectID, s.defaultApiKeyID, req)
	if err != nil {
		log.Printf("[Simulator] generate error: %v", err)
	}
}

// pickModel 按权重随机选择模型
func (s *Simulator) pickModel() string {
	// 权重分布（总和 20）：Qwen-72B:5, Qwen-Coder:4, DeepSeek:4, MiniMax:4, GLM:3
	weights := []struct {
		name   string
		weight int
	}{
		{"qwen2.5-72b", 5},
		{"qwen2.5-coder", 4},
		{"deepseek-v3", 4},
		{"minimax-text-01", 4},
		{"glm-4", 3},
	}

	total := 0
	for _, w := range weights {
		total += w.weight
	}

	r := rand.Intn(total)
	cumulative := 0
	for _, w := range weights {
		cumulative += w.weight
		if r < cumulative {
			return w.name
		}
	}
	return "qwen2.5-72b"
}

// randomInterval 返回 3~6 秒随机间隔
func (s *Simulator) randomInterval() time.Duration {
	ms := 3000 + rand.Intn(3001) // 3000~6000 ms
	return time.Duration(ms) * time.Millisecond
}

// mockUserMsg 返回一条模拟用户消息
func (s *Simulator) mockUserMsg() string {
	messages := []string{
		"请帮我解释一下什么是机器学习",
		"如何优化数据库查询性能？",
		"写一段Python代码实现快速排序",
		"介绍一下微服务架构的优缺点",
		"什么是RESTful API？请举例说明",
		"如何实现一个简单的LRU缓存？",
		"请解释HTTP和HTTPS的区别",
		"什么是Docker容器化技术？",
		"如何设计一个高并发系统？",
		"请解释一下MapReduce的原理",
		"什么是面向对象编程的三大特性？",
		"如何保证接口的安全性？",
		"介绍一下Redis的使用场景",
	}
	return messages[rand.Intn(len(messages))]
}
