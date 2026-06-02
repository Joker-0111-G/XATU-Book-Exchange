package middleware

import (
	"strings"

	"xatu-book-exchange/common"
	"xatu-book-exchange/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired 需要登录的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.Error(c, common.CodeUnauthorized)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			common.Error(c, common.CodeUnauthorized)
			c.Abort()
			return
		}

		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil {
			common.Error(c, common.CodeUnauthorized)
			c.Abort()
			return
		}

		// 检查 Token 是否在 Redis 黑名单中（已登出）
		if utils.IsBlacklisted(claims) {
			common.Error(c, common.CodeUnauthorized)
			c.Abort()
			return
		}

		// 将用户信息注入上下文
		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("token_claims", claims)
		c.Next()
	}
}

// AdminRequired 需要管理员权限的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists {
			common.Error(c, common.CodeUnauthorized)
			c.Abort()
			return
		}

		if isAdmin.(int8) != 1 {
			common.Error(c, common.CodeForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（有 Token 就解析用户信息，没有也继续）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if utils.IsBlacklisted(claims) {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("token_claims", claims)
		c.Next()
	}
}
