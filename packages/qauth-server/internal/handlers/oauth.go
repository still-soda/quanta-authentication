package handlers

import (
	"qauth-server/internal/config"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwt"
	"qauth-server/pkg/response"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// OAuthHandler OAuth2 处理器
type OAuthHandler struct {
	oauthService *services.OAuthService
	roleService  *services.RoleService
	userService  *services.UserService
	oidcService  *services.OIDCService
}

// NewOAuthHandler 创建新的 OAuth2 处理器
func NewOAuthHandler(
	oauthService *services.OAuthService,
	roleService *services.RoleService,
	userService *services.UserService,
	oidcService *services.OIDCService,
) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		roleService:  roleService,
		userService:  userService,
		oidcService:  oidcService,
	}
}

// modifyResponseType 修改响应类型以支持隐式授权模式
func (h *OAuthHandler) modifyResponseType(c *gin.Context, isValid bool, rt string) {
	r := c.Request

	setResponseType := func(t string) {
		q := r.URL.Query()
		q.Set("response_type", t)
		r.URL.RawQuery = q.Encode()
	}

	// 这会让 go-auth2 授权服务器认为请求是不允许的请求
	if !isValid {
		setResponseType("__in_valid__")
		return
	}

	if rt != "code" && rt != "token" {
		setResponseType("token")
	}
}

// modifyRedirectLocation 修改重定向位置以符合隐式授权模式
func (h *OAuthHandler) modifyRedirectLocation(c *gin.Context, isValid bool, rt string) {
	w := c.Writer

	if !isValid || w.Status() != 302 {
		return
	}

	location := w.Header().Get("Location")

	switch rt {
	case "id_token":
		pl, err := utilities.PickHashParams(location, []string{"id_token", "state"})
		if err != nil {
			return
		}
		w.Header().Set("Location", pl)
	}
}

// Authorize 处理授权请求（授权码模式）
// GET/POST /oauth/authorize
func (h *OAuthHandler) Authorize(c *gin.Context) {
	rt := c.Request.URL.Query().Get("response_type")
	arts := h.oauthService.GetConfig().AllowedResponseTypes
	isValid := slices.Contains(arts, config.ResponseType(rt))

	h.modifyResponseType(c, isValid, rt)

	err := h.oauthService.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		utilities.GetLogger().Error("OAuth authorize error", "error", err)
		c.Error(app_error.ErrBadRequest)
	}

	h.modifyRedirectLocation(c, isValid, rt)
}

// Token 处理令牌请求
// POST /oauth/token
func (h *OAuthHandler) Token(c *gin.Context) {
	err := h.oauthService.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		utilities.GetLogger().Error("OAuth token error", "error", err)
		c.Error(app_error.ErrBadRequest)
	}
}

// ValidateToken 验证访问令牌
// POST /oauth/validate
func (h *OAuthHandler) ValidateToken(c *gin.Context) {
	tokenInfo, err := h.oauthService.ValidateToken(c.Request)
	if err != nil {
		c.Error(app_error.ErrUnauthorized)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"active":     true,
		"client_id":  tokenInfo.GetClientID(),
		"user_id":    tokenInfo.GetUserID(),
		"expires_in": int64(tokenInfo.GetAccessExpiresIn().Seconds()),
		"scope":      tokenInfo.GetScope(),
	})
}

// RevokeToken 撤销令牌
// POST /oauth/revoke
func (h *OAuthHandler) RevokeToken(c *gin.Context) {
	var req struct {
		Token         string `json:"token" form:"token" binding:"required"`
		TokenTypeHint string `json:"token_type_hint" form:"token_type_hint"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(app_error.ErrBadRequest)
		return
	}

	var err error
	if req.TokenTypeHint == "refresh_token" {
		err = h.oauthService.RevokeRefreshToken(req.Token)
	} else {
		err = h.oauthService.RevokeToken(req.Token)
	}

	if err != nil {
		utilities.GetLogger().Error("OAuth revoke error", "error", err)
		c.Error(app_error.ErrBadRequest)
		return
	}

	c.Status(200)
}

// CreateClient 创建 OAuth2 客户端
// POST /oauth/clients
func (h *OAuthHandler) CreateClient(c *gin.Context) {
	if err := services.ValidatePermissionCodes(c, h.roleService, []string{
		permissions.OAuthClientCreate,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		Domain   string `json:"domain" binding:"required"`
		Secret   string `json:"secret" binding:"required"`
		IsPublic bool   `json:"is_public"`
		UserID   string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	client, err := h.oauthService.CreateClient(req.Name, req.Domain, req.Secret, req.IsPublic, req.UserID)
	if err != nil {
		utilities.GetLogger().Error("Create OAuth client error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"client_id":     client.ID,
		"client_name":   client.Name,
		"client_domain": client.Domain,
	})
}

// GetClient 获取客户端信息
// GET /oauth/clients/:id
func (h *OAuthHandler) GetClient(c *gin.Context) {
	if err := services.ValidatePermissionCodes(c, h.roleService, []string{
		permissions.OAuthClientView,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	client, err := h.oauthService.GetClientByID(clientID)
	if err != nil {
		response.HandlerError(c, app_error.ErrNotFound)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"client_id":     client.ID,
		"client_name":   client.Name,
		"client_domain": client.Domain,
		"created_at":    client.CreatedAt,
	})
}

// UpdateClient 更新客户端信息
// PUT /oauth/clients/:id
func (h *OAuthHandler) UpdateClient(c *gin.Context) {
	if err := services.ValidatePermissionCodes(c, h.roleService, []string{
		permissions.OAuthClientUpdate,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	var req struct {
		Name   string `json:"name" binding:"required"`
		Domain string `json:"domain" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	if err := h.oauthService.UpdateClient(clientID, req.Name, req.Domain); err != nil {
		utilities.GetLogger().Error("Update OAuth client error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, gin.H{"updated": true})
}

// DeleteClient 删除客户端
// DELETE /oauth/clients/:id
func (h *OAuthHandler) DeleteClient(c *gin.Context) {
	if err := services.ValidatePermissionCodes(c, h.roleService, []string{
		permissions.OAuthClientDelete,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	if err := h.oauthService.DeleteClient(clientID); err != nil {
		utilities.GetLogger().Error("Delete OAuth client error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, gin.H{"deleted": true})
}

// ListClients 获取客户端列表
// GET /oauth/clients
func (h *OAuthHandler) ListClients(c *gin.Context) {
	if err := services.ValidatePermissionCodes(c, h.roleService, []string{
		permissions.OAuthClientList,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if val, err := utilities.ParseInt(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := utilities.ParseInt(ps); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

	clients, total, err := h.oauthService.ListClients(page, pageSize)
	if err != nil {
		utilities.GetLogger().Error("List OAuth clients error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"items": clients,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// UserInfo 获取用户信息
func (h *OAuthHandler) UserInfo(c *gin.Context) {
	token, err := jwt.ExtractTokenFromHeader(c)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	info, err := h.oauthService.GetTokenInfo(token)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	userID := info.GetUserID()
	issueAt := info.GetAccessCreateAt().Unix()
	expiry := info.GetAccessExpiresIn().Seconds() + float64(issueAt)
	audience := info.GetClientID()
	issuer := h.oidcService.GetConfig().Issuer

	user, err := h.userService.GetUserByID(userID, true)
	if err != nil {
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	data := map[string]any{
		string(config.ClaimSub): userID,
		string(config.ClaimIss): issuer,
		string(config.ClaimIat): issueAt,
		string(config.ClaimExp): expiry,
		string(config.ClaimAud): audience,
	}

	scopes := strings.Split(info.GetScope(), " ")

	// 根据 scope 添加不同的声明
	if slices.Contains(scopes, string(config.ScopeProfile)) {
		data[string(config.ClaimName)] = user.Name
		data[string(config.ClaimStudentID)] = user.StudentID
		data[string(config.ClaimDisplayName)] = *user.DisplayName
		data[string(config.ClaimPicture)] = "" // TODO: 添加头像支持
	}

	if slices.Contains(scopes, string(config.ScopeEmail)) {
		data[string(config.ClaimEmail)] = user.Email
		data[string(config.ClaimEmailVerified)] = user.EmailVerified
	}

	if slices.Contains(scopes, string(config.ScopeRoles)) {
		data[string(config.ClaimRoles)] = user.Roles
	}

	response.HandlerSuccess(c, data)
}

// Logout 登出端点
func (h *OAuthHandler) Logout(c *gin.Context) {
	token, err := jwt.ExtractTokenFromHeader(c)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	info, err := h.oauthService.GetTokenInfo(token)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	refreshToken := info.GetRefresh()
	h.oauthService.GetManager().RemoveAccessToken(c, token)
	h.oauthService.GetManager().RemoveRefreshToken(c, refreshToken)

	postLogoutRedirectURI := c.Query("post_logout_redirect_uri")
	if postLogoutRedirectURI == "" {
		postLogoutRedirectURI = "/"
	}

	c.Redirect(302, postLogoutRedirectURI)
}
