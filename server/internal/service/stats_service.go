package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// StatsService 统计服务
type StatsService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// OverviewData 概览数据
type OverviewData struct {
	Today      TodayStats      `json:"today"`
	Total      TotalStats      `json:"total"`
	Trend      TrendPct        `json:"trend"`
	ModelCount ModelCountStats `json:"model_count"`
}

type TodayStats struct {
	TotalTokens   int     `json:"total_tokens"`
	TotalRequests int     `json:"total_requests"`
	SuccessRate   float64 `json:"success_rate"`
	AvgLatencyMs  int     `json:"avg_latency_ms"`
	TotalCost     float64 `json:"total_cost"`
}

type TotalStats struct {
	TotalTokens   int     `json:"total_tokens"`
	TotalRequests int     `json:"total_requests"`
	TotalCost     float64 `json:"total_cost"`
}

type TrendPct struct {
	TokenChangePct   float64 `json:"token_change_pct"`
	RequestChangePct float64 `json:"request_change_pct"`
}

type ModelCountStats struct {
	Online int `json:"online"`
	Total  int `json:"total"`
}

// ModelUsageItem 模型用量
type ModelUsageItem struct {
	ModelName  string  `json:"model_name"`
	Tokens     int     `json:"tokens"`
	Percentage float64 `json:"percentage"`
	Cost       float64 `json:"cost"`
}

// TrendPoint 趋势数据点
type TrendPoint struct {
	Date     string `json:"date"`
	Tokens   int64  `json:"tokens"`
	Requests int    `json:"requests"`
}

// LogEntry 日志条目
type LogEntry struct {
	ID          uint      `json:"id"`
	ModelName   string    `json:"model_name"`
	ProjectName string    `json:"project_name"`
	TotalTokens int       `json:"total_tokens"`
	Cost        float64   `json:"cost"`
	Status      string    `json:"status"`
	LatencyMs   int       `json:"latency_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

// Pagination 分页信息
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewStatsService(db *gorm.DB, rdb *redis.Client) *StatsService {
	return &StatsService{DB: db, RDB: rdb}
}

// GetOverview 获取概览数据
func (s *StatsService) GetOverview(ctx context.Context) (*OverviewData, error) {
	dateKey := time.Now().Format("20060102")
	yesterdayKey := time.Now().AddDate(0, 0, -1).Format("20060102")

	// 从 Redis 读取今日指标
	todayTokens, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:total_tokens", dateKey)).Int()
	todayRequests, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:total_requests", dateKey)).Int()
	succ, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:success_count", dateKey)).Int()
	fail, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:fail_count", dateKey)).Int()
	timeout, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:timeout_count", dateKey)).Int()
	totalLat, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:total_latency_ms", dateKey)).Int()
	todayCost, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:total_cost", dateKey)).Float64()

	// 昨日数据（用于趋势对比）
	yesterdayTokens, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:total_tokens", yesterdayKey)).Int()
	yesterdayRequests, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:total_requests", yesterdayKey)).Int()

	// 平均值计算
	var successRate float64
	totalCalls := succ + fail + timeout
	if totalCalls > 0 {
		successRate = float64(succ) / float64(totalCalls)
	}
	var avgLatency int
	if todayRequests > 0 {
		avgLatency = totalLat / todayRequests
	}

	// 趋势变化百分比
	var tokenPct, requestPct float64
	if yesterdayTokens > 0 {
		tokenPct = float64(todayTokens-yesterdayTokens) / float64(yesterdayTokens) * 100
	}
	if yesterdayRequests > 0 {
		requestPct = float64(todayRequests-yesterdayRequests) / float64(yesterdayRequests) * 100
	}

	// MySQL 累计
	var totalTokens, totalRequests int64
	var totalCost float64
	s.DB.Raw("SELECT COALESCE(SUM(total_tokens),0), COALESCE(COUNT(*),0), COALESCE(SUM(cost),0) FROM api_usage_logs").
		Row().Scan(&totalTokens, &totalRequests, &totalCost)

	// MySQL 模型统计
	var onlineModels, allModels int64
	s.DB.Raw("SELECT COALESCE(SUM(CASE WHEN status='online' THEN 1 ELSE 0 END),0) FROM models").Row().Scan(&onlineModels)
	s.DB.Raw("SELECT COUNT(*) FROM models").Row().Scan(&allModels)

	return &OverviewData{
		Today: TodayStats{
			TotalTokens:   todayTokens,
			TotalRequests: todayRequests,
			SuccessRate:   successRate,
			AvgLatencyMs:  avgLatency,
			TotalCost:     todayCost,
		},
		Total: TotalStats{
			TotalTokens:   int(totalTokens),
			TotalRequests: int(totalRequests),
			TotalCost:     totalCost,
		},
		Trend: TrendPct{
			TokenChangePct:   tokenPct,
			RequestChangePct: requestPct,
		},
		ModelCount: ModelCountStats{
			Online: int(onlineModels),
			Total:  int(allModels),
		},
	}, nil
}

// GetModelUsage 获取各模型用量分布
func (s *StatsService) GetModelUsage(ctx context.Context) ([]ModelUsageItem, int, error) {
	dateKey := time.Now().Format("20060102")
	keys, _ := s.RDB.Keys(ctx, fmt.Sprintf("stats:%s:model:*:tokens", dateKey)).Result()

	var items []ModelUsageItem
	totalTokens := 0
	for _, k := range keys {
		var modelID int
		fmt.Sscanf(k, fmt.Sprintf("stats:%s:model:%%d:tokens", dateKey), &modelID)
		tokens, _ := s.RDB.Get(ctx, k).Int()
		cost, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:model:%d:cost", dateKey, modelID)).Float64()

		var modelName string
		s.DB.Raw("SELECT name FROM models WHERE id = ?", modelID).Row().Scan(&modelName)
		if modelName == "" {
			continue
		}

		totalTokens += tokens
		items = append(items, ModelUsageItem{ModelName: modelName, Tokens: tokens, Cost: cost})
	}

	for i := range items {
		if totalTokens > 0 {
			items[i].Percentage = float64(items[i].Tokens) / float64(totalTokens)
		}
	}

	return items, totalTokens, nil
}

// GetTrend 获取近 N 天趋势
func (s *StatsService) GetTrend(ctx context.Context, days int) ([]TrendPoint, error) {
	var points []TrendPoint
	err := s.DB.Raw(`
		SELECT DATE(created_at) AS date,
			COALESCE(SUM(total_tokens),0) AS tokens,
			COALESCE(COUNT(*),0) AS requests
		FROM api_usage_logs
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, days).Scan(&points).Error
	return points, err
}

// GetRealtimeLogs 获取实时日志（分页）
func (s *StatsService) GetRealtimeLogs(ctx context.Context, page, pageSize int, modelID *int, status *string) ([]LogEntry, *Pagination, error) {
	var logs []LogEntry

	query := s.DB.Table("api_usage_logs l").
		Select("l.id, m.name AS model_name, p.name AS project_name, l.total_tokens, l.cost, l.status, l.latency_ms, l.created_at").
		Joins("LEFT JOIN models m ON m.id = l.model_id").
		Joins("LEFT JOIN projects p ON p.id = l.project_id")

	if modelID != nil && *modelID > 0 {
		query = query.Where("l.model_id = ?", *modelID)
	}
	if status != nil && *status != "" {
		query = query.Where("l.status = ?", *status)
	}

	var total int64
	query.Count(&total)

	if err := query.Order("l.id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Scan(&logs).Error; err != nil {
		return nil, nil, err
	}

	if logs == nil {
		logs = []LogEntry{}
	}

	return logs, &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
