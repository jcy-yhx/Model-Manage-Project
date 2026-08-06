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
func SetupRouter(db *gorm.DB, rdb *redis.Client, chatSvc *service.ChatService, statsSvc *service.StatsService, modelSvc *service.ModelService, defaultQuota int) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.NoRoute(handler.NoRouteHandler)
	r.NoMethod(handler.NoMethodHandler)

	// 公开路由（无需认证）
	healthHandler := &handler.HealthHandler{DB: db, RDB: rdb}
	r.GET("/api/health", healthHandler.Health)

	// 统计路由（无需认证，方便 Dashboard 轮询）
	statsHandler := &handler.StatsHandler{StatsSvc: statsSvc}
	logHandler := &handler.LogHandler{StatsSvc: statsSvc}
	modelHandler := &handler.ModelHandler{ModelSvc: modelSvc}

	r.GET("/api/stats/overview", statsHandler.Overview)
	r.GET("/api/stats/model-usage", statsHandler.ModelUsage)
	r.GET("/api/stats/trend", statsHandler.Trend)
	r.GET("/api/realtime/logs", logHandler.RealtimeLogs)
	r.GET("/api/models", modelHandler.List)

	// API Key 认证路由组
	authorized := r.Group("/api")
	authorized.Use(middleware.ApiKeyAuth(db, rdb, defaultQuota))
	{
		chatHandler := &handler.ChatHandler{Svc: chatSvc}
		authorized.POST("/chat", chatHandler.Chat)
	}

	return r
}
