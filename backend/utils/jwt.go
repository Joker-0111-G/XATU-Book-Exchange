package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"xatu-book-exchange/config"
	"xatu-book-exchange/database"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	IsAdmin   int8   `json:"is_admin"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// generateJTI 生成唯一的 JWT ID（用于黑名单）
func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateToken 生成 access_token 和 refresh_token
func GenerateToken(userID uint, isAdmin int8) (accessToken string, refreshToken string, err error) {
	cfg := config.AppConfig.JWT

	// access_token
	accessClaims := Claims{
		UserID:    userID,
		IsAdmin:   isAdmin,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "xatu-book-exchange",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// refresh_token
	refreshClaims := Claims{
		UserID:    userID,
		IsAdmin:   isAdmin,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "xatu-book-exchange",
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	return
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ParseAccessToken 解析并校验是 access_token
func ParseAccessToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, errors.New("无效的 access_token")
	}
	return claims, nil
}

// ParseRefreshToken 解析并校验是 refresh_token
func ParseRefreshToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, errors.New("无效的 refresh_token")
	}
	return claims, nil
}

// AddToBlacklist 将 Token 的 JTI 加入 Redis 黑名单（到期自动删除）
func AddToBlacklist(claims *Claims) error {
	if database.RDB == nil {
		return nil // Redis 不可用时静默跳过
	}
	key := fmt.Sprintf("token:blacklist:%s", claims.ID)
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return nil
	}
	return database.RDB.Set(context.Background(), key, "1", remaining).Err()
}

// IsBlacklisted 检查 Token 的 JTI 是否在黑名单中
func IsBlacklisted(claims *Claims) bool {
	if database.RDB == nil {
		return false // Redis 不可用时放行
	}
	key := fmt.Sprintf("token:blacklist:%s", claims.ID)
	exists, err := database.RDB.Exists(context.Background(), key).Result()
	if err != nil {
		return false // 查询出错时放行（避免 Redis 故障导致全站不可用）
	}
	return exists == 1
}
