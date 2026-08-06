package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ai-gateway-server/internal/pkg"
)

// context keys for downstream pkgs
const (
	CtxApiKeyID  = "api_key_id"
	CtxProjectID = "project_id"
	CtxKeyID     = "key_id"
)

// ApiKey 数据库查询用的结构
type ApiKey struct {
	ID        uint   `gorm:"primaryKey"`
	ProjectID uint   `gorm:"column:project_id"`
	KeyHash   string `gorm:"column:key_hash"`
	KeyPrefix string `gorm:"column:key_prefix"`
	Status    int8   `gorm:"column:status"`
}

func (ApiKey) TableName() string { return "api_keys" }

// ApiKeyAuth 返回 API Key 校验 + 额度检查中间件
//   - 从 Authorization: Bearer sk-xxx 提取 Key
//   - key_prefix 查库 → bcrypt 比对
//   - 查询 Redis 额度 → 超限返回 429
//   - 注入 api_key_id / project_id / key_id 到 context
func ApiKeyAuth(db *gorm.DB, rdb *redis.Client, defaultQuota int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 提取 API Key
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			pkg.Unauthorized(c, "缺少 API Key，请在 Authorization Header 中提供 Bearer sk-xxx")
			c.Abort()
			return
		}
		rawKey := strings.TrimPrefix(auth, "Bearer ")
		if rawKey == "" {
			pkg.Unauthorized(c, "API Key 不能为空")
			c.Abort()
			return
		}

		// 2. key_prefix 查库
		prefix := rawKey
		if len(prefix) > 12 {
			prefix = prefix[:12]
		}

		var keys []ApiKey
		if err := db.Where("key_prefix = ? AND status = 1", prefix).Find(&keys).Error; err != nil {
			pkg.InternalError(c, "校验 API Key 时发生错误")
			c.Abort()
			return
		}

		// 3. bcrypt 比对
		var matched *ApiKey
		for i := range keys {
			if bcrypt.CompareHashAndPassword([]byte(keys[i].KeyHash), []byte(rawKey)) == nil {
				matched = &keys[i]
				break
			}
		}
		if matched == nil {
			pkg.Unauthorized(c, "无效的 API Key")
			c.Abort()
			return
		}

		// 4. 额度检查
		dateKey := time.Now().Format("20060102")
		quotaKey := fmt.Sprintf("quota:key:%d:%s:tokens", matched.ID, dateKey)

		used, _ := rdb.Get(c.Request.Context(), quotaKey).Int()
		if used >= defaultQuota {
			pkg.TooMany(c,
				"当日 Token 额度已用尽，请明天再试或联系项目管理员提升额度",
				gin.H{
					"quota_limit": defaultQuota,
					"quota_used":  used,
					"reset_at":    time.Now().Format("2006-01-02") + "T00:00:00Z",
				})
			c.Abort()
			return
		}

		// 5. 注入 context
		c.Set(CtxApiKeyID, matched.ID)
		c.Set(CtxProjectID, matched.ProjectID)
		c.Set(CtxKeyID, matched.ID)

		c.Next()
	}
}

// GetApiKeyID 从 context 中提取 api_key_id
func GetApiKeyID(c *gin.Context) uint {
	v, _ := c.Get(CtxApiKeyID)
	id, _ := v.(uint)
	return id
}

// GetProjectID 从 context 中提取 project_id
func GetProjectID(c *gin.Context) uint {
	v, _ := c.Get(CtxProjectID)
	id, _ := v.(uint)
	return id
}
