package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ai-gateway-server/internal/service"
)

// LogHandler 日志接口
type LogHandler struct {
	StatsSvc *service.StatsService
}

// RealtimeLogs GET /api/realtime/logs
func (h *LogHandler) RealtimeLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var modelID *int
	if v := c.Query("model_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			modelID = &id
		}
	}

	var status *string
	if v := c.Query("status"); v != "" {
		status = &v
	}

	logs, pagination, err := h.StatsSvc.GetRealtimeLogs(c.Request.Context(), page, pageSize, modelID, status)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{
		"logs":       logs,
		"pagination": pagination,
	})
}
