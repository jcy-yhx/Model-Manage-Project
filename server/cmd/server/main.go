package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"ai-gateway-server/internal/config"
	"ai-gateway-server/internal/router"
	"ai-gateway-server/internal/service"
	"ai-gateway-server/internal/simulator"
)

func main() {
	// ── 1. 加载配置 ──
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Init] failed to load config: %v", err)
	}

	// ── 2. 连接 MySQL ──
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		log.Fatalf("[Init] failed to connect MySQL: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Println("[Init] MySQL connected")

	// ── 3. 连接 Redis ──
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("[Init] failed to connect Redis: %v", err)
	}
	log.Println("[Init] Redis connected")

	// ── 4. 初始化服务 ──
	chatSvc := service.NewChatService(db, rdb, cfg.Quota.DefaultDailyTokens)

	// ── 5. 启动模拟器（后台 goroutine） ──
	sim := simulator.New(chatSvc)
	if cfg.Simulator.Enabled {
		sim.Start()
		log.Println("[Init] Simulator started")
	}

	// ── 6. 路由 ──
	r := router.SetupRouter(db, rdb, chatSvc, cfg.Quota.DefaultDailyTokens)

	// ── 7. 启动 HTTP Server ──
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("[Server] listening on %s (mode: %s)", addr, cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] failed to start: %v", err)
		}
	}()

	// ── 8. 优雅关闭 ──
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Server] shutting down...")
	sim.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[Server] forced shutdown: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("[Server] MySQL close error: %v", err)
	}
	if err := rdb.Close(); err != nil {
		log.Printf("[Server] Redis close error: %v", err)
	}

	log.Println("[Server] stopped")
}
