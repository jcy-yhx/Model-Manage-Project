package model

import "time"

// ApiUsageLog API 调用日志表 (api_usage_logs)
type ApiUsageLog struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID        uint      `gorm:"column:project_id;type:bigint;not null;index:idx_project_id_created_at,priority:1" json:"project_id"`
	ApiKeyID         uint      `gorm:"column:api_key_id;type:bigint;not null;index:idx_api_key_id_created_at,priority:1" json:"api_key_id"`
	ModelID          uint      `gorm:"column:model_id;type:bigint;not null;index:idx_model_id_created_at,priority:1" json:"model_id"`
	PromptTokens     int       `gorm:"column:prompt_tokens;type:int;not null;default:0" json:"prompt_tokens"`
	CompletionTokens int       `gorm:"column:completion_tokens;type:int;not null;default:0" json:"completion_tokens"`
	TotalTokens      int       `gorm:"column:total_tokens;type:int;not null;default:0" json:"total_tokens"`
	Cost             float64   `gorm:"column:cost;type:decimal(10,6);not null;default:0" json:"cost"`
	LatencyMs        int       `gorm:"column:latency_ms;type:int;not null;default:0" json:"latency_ms"`
	Status           string    `gorm:"column:status;type:varchar(16);not null;default:success;index" json:"status"`
	ErrorMessage     string    `gorm:"column:error_message;type:varchar(256);default:''" json:"error_message"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;not null;index:idx_created_at;index:idx_model_id_created_at,priority:2;index:idx_project_id_created_at,priority:2;index:idx_api_key_id_created_at,priority:2;autoCreateTime" json:"created_at"`

	// 关联
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	ApiKey  ApiKey  `gorm:"foreignKey:ApiKeyID" json:"api_key,omitempty"`
	Model   Model   `gorm:"foreignKey:ModelID" json:"model,omitempty"`
}

func (ApiUsageLog) TableName() string {
	return "api_usage_logs"
}
