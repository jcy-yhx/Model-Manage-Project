package model

import "time"

// Model 大模型信息表 (models)
type Model struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(64);uniqueIndex;not null" json:"name"`
	DisplayName string    `gorm:"column:display_name;type:varchar(128);default:''" json:"display_name"`
	Provider    string    `gorm:"column:provider;type:varchar(32);default:''" json:"provider"`
	Status      string    `gorm:"column:status;type:varchar(16);not null;default:online" json:"status"`
	InputPrice  float64   `gorm:"column:input_price;type:decimal(10,6);not null;default:0" json:"input_price"`
	OutputPrice float64   `gorm:"column:output_price;type:decimal(10,6);not null;default:0" json:"output_price"`
	Description string    `gorm:"column:description;type:varchar(256);default:''" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Model) TableName() string {
	return "models"
}
