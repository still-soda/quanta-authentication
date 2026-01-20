package jwt

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getTokenExpireTime(tokenType string) int64 {
	switch tokenType {
	case "access":
		return time.Now().Add(24 * time.Hour).Unix()
	case "refresh":
		return time.Now().Add(7 * 24 * time.Hour).Unix()
	default:
		return time.Now().Add(24 * time.Hour).Unix()
	}
}

func generateToken(tokenType string, claims map[string]any) (string, error) {
	var key string
	switch tokenType {
	case "access":
		key = "your-access-secret-key"
	default:
		key = "your-refresh-secret-key"
	}

	token, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims)).
		SignedString([]byte(key))
	return token, err
}

func parseToken(tokenType, tokenString string) (*jwt.Token, error) {
	var key string
	switch tokenType {
	case "access":
		key = "your-access-secret-key"
	default:
		key = "your-refresh-secret-key"
	}
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})
}

// JWTInfo 包含生成 JWT 所需的信息
type JWTInfo struct {
	UserID    string
	StudentID string
	Role      string
}

// GenerateAccessToken 生成 JWT 令牌
func GenerateAccessToken(info *JWTInfo) (string, error) {
	jwtInfo := map[string]any{
		"user_id":    info.UserID,
		"student_id": info.StudentID,
		"role":       info.Role,
		"exp":        getTokenExpireTime("access"),
	}
	return generateToken("access", jwtInfo)
}

// GenerateRefreshToken 生成刷新用的 JWT 令牌
func GenerateRefreshToken(info *JWTInfo) (string, error) {
	jwtInfo := map[string]any{
		"user_id":    info.UserID,
		"student_id": info.StudentID,
		"role":       info.Role,
		"exp":        getTokenExpireTime("refresh"),
	}
	return generateToken("refresh", jwtInfo)
}

// UserJWTClaims 包含解析后的 JWT 信息
type UserJWTClaims struct {
	UserID    string
	StudentID string
	Role      string
	Exp       int64
}

// ParseAccessToken 解析 JWT 令牌
func ParseAccessToken(tokenString string) (*UserJWTClaims, error) {
	token, err := parseToken("access", tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &UserJWTClaims{
			UserID:    claims["user_id"].(string),
			StudentID: claims["student_id"].(string),
			Role:      claims["role"].(string),
			Exp:       int64(claims["exp"].(float64)),
		}, nil
	} else {
		return nil, err
	}
}

// ParseRefreshToken 解析刷新用的 JWT 令牌
func ParseRefreshToken(tokenString string) (*UserJWTClaims, error) {
	token, err := parseToken("refresh", tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &UserJWTClaims{
			UserID:    claims["user_id"].(string),
			StudentID: claims["student_id"].(string),
			Role:      claims["role"].(string),
			Exp:       int64(claims["exp"].(float64)),
		}, nil
	} else {
		return nil, err
	}
}

// VerifyJWTToken 验证 JWT 令牌的有效性
func VerifyJWTToken(tokenString string) bool {
	_, err := ParseAccessToken(tokenString)
	return err == nil
}

// ExtractTokenFromHeader 从 Gin 上下文中提取 JWT 令牌
func ExtractTokenFromHeader(ctx *gin.Context) (string, error) {
	// 获取 Authorization 头
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "authorization header is required"})
		return "", fmt.Errorf("authorization header is required")
	}
	// 提取 token
	var accessToken string
	_, err := fmt.Sscanf(authorization, "Bearer %s", &accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header format"})
		return "", fmt.Errorf("invalid authorization header format")
	}

	return accessToken, nil
}
