package router

import "github.com/gin-gonic/gin"

// SetupRouter 注册所有路由和中间件，返回 gin Engine 实例
func SetupRouter() *gin.Engine {
	r := gin.New()
	return r
}
