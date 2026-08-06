package model

import "time"

// Project 项目表 (projects)
type Project struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	OwnerID     uint      `gorm:"column:owner_id;type:bigint;not null;index" json:"owner_id"`
	Description string    `gorm:"column:description;type:varchar(512);default:''" json:"description"`
	Status      int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// 关联
	Owner   User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	ApiKeys []ApiKey `gorm:"foreignKey:ProjectID" json:"api_keys,omitempty"`
}

func (Project) TableName() string {
	return "projects"
}
