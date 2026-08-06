package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ai-gateway-server/internal/model"
)

// ModelService 模型信息服务
type ModelService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewModelService(db *gorm.DB, rdb *redis.Client) *ModelService {
	return &ModelService{DB: db, RDB: rdb}
}

// ModelWithStats 模型信息 + 今日统计
type ModelWithStats struct {
	model.Model
	TodayTokens   int     `json:"today_tokens"`
	TodayRequests int     `json:"today_requests"`
	AvgLatencyMs  int     `json:"avg_latency_ms"`
	TodayCost     float64 `json:"today_cost"`
}

// GetModels 获取模型列表 + 今日指标
func (s *ModelService) GetModels(ctx context.Context) ([]ModelWithStats, error) {
	var models []model.Model
	if err := s.DB.Find(&models).Error; err != nil {
		return nil, err
	}

	dateKey := time.Now().Format("20060102")
	result := make([]ModelWithStats, 0, len(models))

	for _, m := range models {
		tokens, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:model:%d:tokens", dateKey, m.ID)).Int()
		requests, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:model:%d:requests", dateKey, m.ID)).Int()
		cost, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:model:%d:cost", dateKey, m.ID)).Float64()
		latency, _ := s.RDB.Get(ctx, fmt.Sprintf("stats:%s:model:%d:latency_ms", dateKey, m.ID)).Int()

		var avgLat int
		if requests > 0 {
			avgLat = latency / requests
		}

		result = append(result, ModelWithStats{
			Model:         m,
			TodayTokens:   tokens,
			TodayRequests: requests,
			AvgLatencyMs:  avgLat,
			TodayCost:     cost,
		})
	}

	return result, nil
}
