package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ai-gateway-server/internal/service"
)

// StatsHandler 统计接口
type StatsHandler struct {
	StatsSvc *service.StatsService
}

// Overview GET /api/stats/overview
func (h *StatsHandler) Overview(c *gin.Context) {
	data, err := h.StatsSvc.GetOverview(c.Request.Context())
	if err != nil {
		InternalError(c, "获取概览数据失败: "+err.Error())
		return
	}
	Success(c, data)
}

// ModelUsage GET /api/stats/model-usage
func (h *StatsHandler) ModelUsage(c *gin.Context) {
	items, totalTokens, err := h.StatsSvc.GetModelUsage(c.Request.Context())
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, gin.H{
		"models":       items,
		"total_tokens": totalTokens,
	})
}

// Trend GET /api/stats/trend?days=7
func (h *StatsHandler) Trend(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil && v > 0 && v <= 90 {
			days = v
		}
	}
	points, err := h.StatsSvc.GetTrend(c.Request.Context(), days)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if points == nil {
		points = []service.TrendPoint{}
	}
	Success(c, gin.H{"points": points})
}
