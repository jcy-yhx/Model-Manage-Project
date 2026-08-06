package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ai-gateway-server/internal/pkg"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// Health 健康检查响应
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "connected"
	if err := h.DB.Raw("SELECT 1").Error; err != nil {
		dbStatus = "disconnected: " + err.Error()
	}

	redisStatus := "connected"
	if err := h.RDB.Ping(c.Request.Context()).Err(); err != nil {
		redisStatus = "disconnected: " + err.Error()
	}

	// 统计模型状态
	var totalModels int64
	var onlineModels int64
	h.DB.Model(&Model{}).Count(&totalModels)
	h.DB.Model(&Model{}).Where("status = ?", "online").Count(&onlineModels)

	// 检测异常模型
	var alerts []map[string]string
	var degradedModels []Model
	h.DB.Where("status IN ?", []string{"degraded", "offline"}).Find(&degradedModels)
	for _, m := range degradedModels {
		alerts = append(alerts, map[string]string{
			"level":   "warning",
			"model":   m.Name,
			"message": "模型状态为 " + m.Status,
		})
	}

	c.JSON(http.StatusOK, pkg.Response{
		Code:    0,
		Message: "ok",
		Data: map[string]interface{}{
			"db":            dbStatus,
			"redis":         redisStatus,
			"online_models": onlineModels,
			"total_models":  totalModels,
			"alerts":        alerts,
		},
	})
}

// Model 本地引用，避免 import cycle（与 internal/model 包同名但仅用于此 handler）
type Model struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"column:name"`
	Status string `gorm:"column:status"`
}

func (Model) TableName() string { return "models" }
