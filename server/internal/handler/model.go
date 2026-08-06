package handler

import (
	"github.com/gin-gonic/gin"

	"ai-gateway-server/internal/service"
)

// ModelHandler 模型接口
type ModelHandler struct {
	ModelSvc *service.ModelService
}

// List GET /api/models
func (h *ModelHandler) List(c *gin.Context) {
	models, err := h.ModelSvc.GetModels(c.Request.Context())
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if models == nil {
		models = []service.ModelWithStats{}
	}
	Success(c, gin.H{"models": models})
}
