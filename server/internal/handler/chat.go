package handler

import (
	"github.com/gin-gonic/gin"

	"ai-gateway-server/internal/service"
)

// ChatHandler 大模型调用接口
type ChatHandler struct {
	Svc *service.ChatService
}

// Chat POST /api/chat
func (h *ChatHandler) Chat(c *gin.Context) {
	var req service.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数校验失败："+err.Error())
		return
	}

	// 从 context 获取中间件注入的信息
	apiKeyID, _ := c.Get("api_key_id")
	projectID, _ := c.Get("project_id")

	result, err := h.Svc.Chat(c.Request.Context(),
		projectID.(uint),
		apiKeyID.(uint),
		req,
	)
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	Success(c, result.Response)
}
