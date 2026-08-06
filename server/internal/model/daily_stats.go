package model

import "time"

// DailyUsageStat 日汇总统计表 (daily_usage_stats)
type DailyUsageStat struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StatDate       time.Time `gorm:"column:stat_date;type:date;not null;uniqueIndex:idx_date_model,priority:1;index:idx_stat_date" json:"stat_date"`
	ModelID        *uint     `gorm:"column:model_id;type:bigint;uniqueIndex:idx_date_model,priority:2" json:"model_id"`
	TotalTokens    int64     `gorm:"column:total_tokens;type:bigint;not null;default:0" json:"total_tokens"`
	TotalRequests  int       `gorm:"column:total_requests;type:int;not null;default:0" json:"total_requests"`
	SuccessCount   int       `gorm:"column:success_count;type:int;not null;default:0" json:"success_count"`
	FailCount      int       `gorm:"column:fail_count;type:int;not null;default:0" json:"fail_count"`
	TotalCost      float64   `gorm:"column:total_cost;type:decimal(12,6);not null;default:0" json:"total_cost"`
	AvgLatencyMs   int       `gorm:"column:avg_latency_ms;type:int;not null;default:0" json:"avg_latency_ms"`
}

func (DailyUsageStat) TableName() string {
	return "daily_usage_stats"
}
