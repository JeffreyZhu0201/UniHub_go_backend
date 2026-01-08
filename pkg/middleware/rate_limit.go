package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unihub/internal/model"
)

// simple in-memory rate limiter for demo
var (
	rateLimits = make(map[string][]time.Time) // AppID -> timestamps
	mu         sync.Mutex
)

// OpenAPIMiddleware 开放平台鉴权与限流
func OpenAPIMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.GetHeader("X-App-ID")
		appSecret := c.GetHeader("X-App-Secret")

		if appID == "" || appSecret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少 AppID 或 AppSecret"})
			return
		}

		// 1. 验证身份
		var app model.App
		if err := db.Where("app_id = ? AND app_secret = ?", appID, appSecret).First(&app).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
			return
		}

		// 2. 限流 (简单滑动窗口实现)
		limit := app.RateLimit
		now := time.Now()

		mu.Lock()
		history, exists := rateLimits[appID]
		if !exists {
			history = []time.Time{}
		}

		// 清理1分钟前的请求
		var valid []time.Time
		for _, t := range history {
			if now.Sub(t) < time.Minute {
				valid = append(valid, t)
			}
		}

		if len(valid) >= limit {
			rateLimits[appID] = valid
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "调用频率超限"})
			return
		}

		valid = append(valid, now)
		rateLimits[appID] = valid
		mu.Unlock()

		c.Set("app_id", appID) // 存储上下文
		c.Next()
	}
}
