package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应 (200)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// Error 业务错误响应
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 参数校验失败 (422)
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, 422, message)
}

// Unauthorized 认证失败 (401)
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 401, message)
}

// NotFound 资源不存在 (404)
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 404, message)
}

// TooMany 频率/额度超限 (429)
func TooMany(c *gin.Context, message string, data interface{}) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, Response{
		Code:    429,
		Message: message,
		Data:    data,
	})
}

// InternalError 服务内部错误 (500)
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}
