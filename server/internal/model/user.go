package model

import "time"

// User 用户表 (users)
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null" json:"username"`
	DisplayName  string    `gorm:"column:display_name;type:varchar(64);default:''" json:"display_name"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(256);not null" json:"-"`
	Role         string    `gorm:"column:role;type:varchar(32);not null;default:student" json:"role"`
	Department   string    `gorm:"column:department;type:varchar(128);default:''" json:"department"`
	Status       int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// 关联
	Projects []Project `gorm:"foreignKey:OwnerID" json:"projects,omitempty"`
}

func (User) TableName() string {
	return "users"
}
