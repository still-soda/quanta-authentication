package handlers

import (
	"qauth-server/internal/services"
	"qauth-server/pkg/app_error"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// OIDCHandler OIDC 处理器
type OIDCHandler struct {
	oidcService *services.OIDCService
}

// NewOIDCHandler 创建新的 OIDC 处理器
func NewOIDCHandler(oidcService *services.OIDCService) *OIDCHandler {
	return &OIDCHandler{oidcService: oidcService}
}

// GetOpenIDConfiguration 获取 OpenID Connect 配置
// GET /.well-known/openid-configuration
func (h *OIDCHandler) GetOpenIDConfiguration(c *gin.Context) {
	config := h.oidcService.GetOpenIDConfiguration()
	c.JSON(200, config)
}

// GetJWKS 获取 JWKS
// GET /.well-known/jwks.json
func (h *OIDCHandler) GetJWKS(c *gin.Context) {
	jwks := h.oidcService.GetJWKS()

	// 设置缓存头
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("Content-Type", "application/json")

	c.JSON(200, jwks)
}

// ForceKeyRotation 强制密钥轮换（管理员端点）
// POST /admin/jwks/rotate
func (h *OIDCHandler) ForceKeyRotation(c *gin.Context) {
	if err := h.oidcService.ForceKeyRotation(); err != nil {
		response.HandlerError(c, err)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"message":       "key rotation completed",
		"active_key_id": h.oidcService.GetActiveKeyID(),
	})
}

// GetKeyRotationInfo 获取密钥轮换信息（管理员端点）
// GET /admin/jwks/keys
func (h *OIDCHandler) GetKeyRotationInfo(c *gin.Context) {
	keys := h.oidcService.GetKeyRotationInfo()

	// 转换为安全的响应格式（不暴露私钥信息）
	keyInfos := make([]gin.H, 0, len(keys))
	for _, key := range keys {
		keyInfos = append(keyInfos, gin.H{
			"id":         key.ID,
			"status":     key.Status,
			"created_at": key.CreatedAt,
			"expires_at": key.ExpiresAt,
		})
	}

	response.HandlerSuccess(c, gin.H{
		"active_key_id": h.oidcService.GetActiveKeyID(),
		"keys":          keyInfos,
		"total":         len(keys),
	})
}

// GetUserInfo 获取用户信息
// GET /oauth/userinfo
func (h *OIDCHandler) GetUserInfo(c *gin.Context) {
	// TODO: 从 access token 中获取用户信息
	// 这需要验证 Bearer token 并返回用户信息

	// 暂时返回未实现
	response.HandlerError(c, app_error.NewAppError(501, "userinfo endpoint is not implemented yet"))
}

// Logout 登出端点
// GET/POST /oauth/logout
func (h *OIDCHandler) Logout(c *gin.Context) {
	// TODO: 实现登出逻辑
	// 1. 撤销所有相关 token
	// 2. 清除 session
	// 3. 重定向到 post_logout_redirect_uri

	postLogoutRedirectURI := c.Query("post_logout_redirect_uri")
	if postLogoutRedirectURI == "" {
		postLogoutRedirectURI = "/"
	}

	c.Redirect(302, postLogoutRedirectURI)
}
