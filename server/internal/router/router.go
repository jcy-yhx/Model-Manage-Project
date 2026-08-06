package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ai-gateway-server/internal/handler"
	"ai-gateway-server/internal/middleware"
)

// SetupRouter 注册所有路由和中间件
// db 和 rdb 注入到 handler 中
func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	// 根据配置设置 gin mode
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// 自定义 404 / 405
	r.NoRoute(handler.NoRouteHandler)
	r.NoMethod(handler.NoMethodHandler)

	// 健康检查 handler
	healthHandler := &handler.HealthHandler{DB: db, RDB: rdb}

	// 公开路由（无需认证）
	r.GET("/api/health", healthHandler.Health)

	// TODO 阶段 4：API Key 认证路由组
	// authorized := r.Group("/api")
	// authorized.Use(middleware.ApiKeyAuth(db))
	// { ... }

	return r
}
