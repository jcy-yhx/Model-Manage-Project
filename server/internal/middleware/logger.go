package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件：方法 | 路径 | 耗时 | 状态码 | 客户端 IP
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Printf("[HTTP] %s | %s | %d | %v | %s",
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency.Round(time.Microsecond),
			clientIP,
		)
	}
}
