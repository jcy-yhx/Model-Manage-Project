package handler

import (
	"github.com/gin-gonic/gin"

	"ai-gateway-server/internal/pkg"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) { pkg.Success(c, data) }

// BadRequest 参数校验失败
func BadRequest(c *gin.Context, msg string) { pkg.BadRequest(c, msg) }

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) { pkg.NotFound(c, msg) }

// InternalError 内部错误
func InternalError(c *gin.Context, msg string) { pkg.InternalError(c, msg) }

// NoRouteHandler 404
func NoRouteHandler(c *gin.Context) { pkg.NotFound(c, "请求的资源不存在") }

// NoMethodHandler 405
func NoMethodHandler(c *gin.Context) { pkg.Error(c, 405, 405, "请求方法不允许") }
