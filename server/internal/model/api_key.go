package model

import "time"

// ApiKey API 密钥表 (api_keys)
type ApiKey struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID  uint       `gorm:"column:project_id;type:bigint;not null;index" json:"project_id"`
	Name       string     `gorm:"column:name;type:varchar(64);default:''" json:"name"`
	KeyHash    string     `gorm:"column:key_hash;type:varchar(256);not null" json:"-"`
	KeyPrefix  string     `gorm:"column:key_prefix;type:varchar(16);default:'';index" json:"key_prefix"`
	Status     int8       `gorm:"column:status;type:tinyint;not null;default:1;index" json:"status"`
	LastUsedAt *time.Time `gorm:"column:last_used_at" json:"last_used_at"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// 关联
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (ApiKey) TableName() string {
	return "api_keys"
}
