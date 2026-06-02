package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
}

type visitorInfo struct {
	count    int
	lastSeen time.Time
}

var (
	visitMu     sync.Mutex
	visitors    = make(map[string]*visitorInfo)
	cleanupOnce sync.Once
)

func startCleanup() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			visitMu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(visitors, ip)
				}
			}
			visitMu.Unlock()
		}
	}()
}

// RateLimiter 简单的内存限流中间件
func RateLimiter(cfg RateLimitConfig) gin.HandlerFunc {
	cleanupOnce.Do(startCleanup)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		visitMu.Lock()
		v, exists := visitors[ip]
		now := time.Now()

		if !exists || now.Sub(v.lastSeen) > cfg.Window {
			visitors[ip] = &visitorInfo{count: 1, lastSeen: now}
			visitMu.Unlock()
			c.Next()
			return
		}

		v.count++
		v.lastSeen = now

		if v.count > cfg.MaxRequests {
			visitMu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    1001,
				"message": "请求过于频繁，请稍后再试",
				"data":    nil,
			})
			return
		}
		visitMu.Unlock()
		c.Next()
	}
}
