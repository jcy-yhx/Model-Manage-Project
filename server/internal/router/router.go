package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ai-gateway-server/internal/handler"
	"ai-gateway-server/internal/middleware"
	"ai-gateway-server/internal/service"
)

// SetupRouter 注册所有路由和中间件
func SetupRouter(db *gorm.DB, rdb *redis.Client, chatSvc *service.ChatService, defaultQuota int) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// 自定义 404 / 405
	r.NoRoute(handler.NoRouteHandler)
	r.NoMethod(handler.NoMethodHandler)

	// 公开路由
	healthHandler := &handler.HealthHandler{DB: db, RDB: rdb}
	r.GET("/api/health", healthHandler.Health)

	// API Key 认证路由组
	authorized := r.Group("/api")
	authorized.Use(middleware.ApiKeyAuth(db, rdb, defaultQuota))
	{
		chatHandler := &handler.ChatHandler{Svc: chatSvc}
		authorized.POST("/chat", chatHandler.Chat)
	}

	return r
}
