package handlers

import (
	"encoding/json"
	"net/http"
	"qauth-server/internal/config"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwt"
	"qauth-server/pkg/response"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OAuthHandler OAuth2 处理器
type OAuthHandler struct {
	oauthService *services.OAuthService
	roleService  *services.RoleService
	userService  *services.UserService
	oidcService  *services.OIDCService
	cacheService *services.CacheService
	auditService *services.AuditService
}

// NewOAuthHandler 创建新的 OAuth2 处理器
func NewOAuthHandler(
	oauthService *services.OAuthService,
	roleService *services.RoleService,
	userService *services.UserService,
	oidcService *services.OIDCService,
	cacheService *services.CacheService,
	auditService *services.AuditService,
) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		roleService:  roleService,
		userService:  userService,
		oidcService:  oidcService,
		cacheService: cacheService,
		auditService: auditService,
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

// requestQueryToJSON 将请求查询参数转换为 JSON 字符串
func (h *OAuthHandler) requestQueryToJSON(r *http.Request) (string, error) {
	q := r.URL.Query()
	params := make(map[string]any)
	for key, values := range q {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	jsonStr, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}

// AuthorizePage 显示授权页面
// GET /oauth/authorize
func (h *OAuthHandler) AuthorizePage(c *gin.Context) {
	// 验证授权请求
	_, err := h.oauthService.ValidateAuthorizeRequest(c.Request)
	if err != nil {
		utilities.GetLogger().Error("Validate authorize request error", "error", err)
		c.Error(app_error.ErrBadRequest)
		return
	}

	// 将请求查询参数转换为 JSON 字符串
	jsonStr, err := h.requestQueryToJSON(c.Request)
	if err != nil {
		utilities.GetLogger().Error("Convert request query to JSON error", "error", err)
		c.Error(app_error.ErrInternalServerError)
		return
	}

	// 生成 authorizeID
	authorizeID, err := utilities.GenerateRandomString(32)
	if err != nil {
		utilities.GetLogger().Error("Generate authorize ID error", "error", err)
		c.Error(app_error.ErrInternalServerError)
		return
	}

	// 将 authorizeID 和请求参数存储在缓存中，有效期为 60 秒
	h.cacheService.SetKeyValue("authorize-"+authorizeID, jsonStr, 120)

	// 重定向到授权页面
	authorizePageUrl := h.oauthService.GetConfig().AuthorizePageURL + "?aid=" + authorizeID
	c.Redirect(302, authorizePageUrl)
}

// AuthorizeInfo 获取授权信息
// GET /oauth/authorize/info
func (h *OAuthHandler) AuthorizeInfo(c *gin.Context) {
	r := c.Request

	// TODO: 验证 referer 是否来自授权页面

	// 从缓存中获取请求参数
	authorizeID := r.URL.Query().Get("aid")
	data, err := h.cacheService.GetKeyValue("authorize-" + authorizeID)
	if err != nil || data == "" {
		utilities.GetLogger().Error("Get authorize info from cache error", "error", err)
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 将 JSON 字符串解析为 map
	var params map[string]any
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		utilities.GetLogger().Error("Parse authorize info JSON error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, params)
}

// verifyAuthorizeID 验证 authorizeID 是否有效且请求 URI 一致
func (h *OAuthHandler) restoreRequestQuery(c *gin.Context, aid string) error {
	// 从缓存中获取请求参数
	data, err := h.cacheService.GetKeyValue("authorize-" + aid)
	if err != nil || data == "" {
		utilities.GetLogger().Error("Get authorize info from cache error", "error", err)
		return app_error.ErrBadRequest
	}

	// 将 JSON 字符串解析为 map
	var params map[string]any
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		utilities.GetLogger().Error("Parse authorize info JSON error", "error", err)
		return app_error.ErrInternalServerError
	}

	// 修改请求的 URL 查询参数
	q := c.Request.URL.Query()
	for key, value := range params {
		q.Set(key, value.(string))
	}

	c.Request.URL.RawQuery = q.Encode()

	return nil
}

// Authorize 处理授权请求（授权码模式）
// POST /oauth/authorize
func (h *OAuthHandler) Authorize(c *gin.Context) {
	startTime := time.Now()

	// 验证 authorizeID
	authorizeID := c.Query("aid")
	if err := h.restoreRequestQuery(c, authorizeID); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 获取授权参数
	clientID := c.Request.URL.Query().Get("client_id")
	scope := c.Request.URL.Query().Get("scope")
	redirectUri := c.Query("redirect_uri")
	responseType := c.Request.URL.Query().Get("response_type")
	scopes := strings.Split(scope, " ")

	// 验证授权结果
	action := c.PostForm("action")
	if action == "deny" {
		// 记录授权拒绝
		h.auditService.LogOAuthAuthorize(c, "", "", &clientID, scopes, "", responseType, redirectUri, false, "user denied")
		h.oauthService.HandlerAuthorizeDeny(c, redirectUri)
		return
	}

	// 处理非法请求
	if action != "approve" {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	c.Request.RequestURI = ""

	// 删除缓存中的授权请求
	h.cacheService.DeleteKey("authorize-" + authorizeID)

	// 处理授权请求
	rt := c.Request.URL.Query().Get("response_type")
	arts := h.oauthService.GetConfig().AllowedResponseTypes
	isValid := slices.Contains(arts, config.ResponseType(rt))

	// 根据响应类型修改请求以支持隐式授权模式
	h.modifyResponseType(c, isValid, rt)

	// 处理授权请求
	err := h.oauthService.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		utilities.GetLogger().Error("OAuth authorize error", "error", err)
		// 记录授权失败
		h.auditService.LogOAuthAuthorize(c, "", "", &clientID, scopes, "", responseType, redirectUri, false, err.Error())
		c.Error(app_error.ErrBadRequest)
		return
	}

	// 记录授权成功
	userID := ""
	userName := ""
	if uid, exists := c.Get("userID"); exists {
		userID = uid.(string)
		if user, err := h.userService.GetUserByID(userID, false); err == nil {
			userName = user.Name
		}
	}
	var operatorID *string
	if userID != "" {
		operatorID = &userID
	}
	h.auditService.Log(&services.AuditContext{
		OperatorID:   operatorID,
		OperatorName: userName,
		IP:           c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		ClientID:     &clientID,
	}, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionOAuthAuthorize,
		TargetID:   clientID,
		TargetType: "oauth_client",
		Detail: &models.AuditLogDetail{
			Scopes:       scopes,
			ResponseType: responseType,
			RedirectURI:  redirectUri,
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	// 增加今日授权计数器
	h.cacheService.IncrKey("todays-authcount")

	// 根据响应类型修改重定向位置以支持隐式授权模式
	h.modifyRedirectLocation(c, isValid, rt)
}

// Token 处理令牌请求
// POST /oauth/token
func (h *OAuthHandler) Token(c *gin.Context) {
	startTime := time.Now()
	clientID := c.PostForm("client_id")
	grantType := c.PostForm("grant_type")

	err := h.oauthService.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		utilities.GetLogger().Error("OAuth token error", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:     models.AuditModuleOAuth,
			Action:     models.AuditActionOAuthToken,
			TargetID:   clientID,
			TargetType: "oauth_client",
			Detail: &models.AuditLogDetail{
				GrantType:  grantType,
				FailReason: err.Error(),
			},
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		c.Error(app_error.ErrBadRequest)
		return
	}

	// 记录令牌请求成功
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionOAuthToken,
		TargetID:   clientID,
		TargetType: "oauth_client",
		Detail: &models.AuditLogDetail{
			GrantType: grantType,
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})
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
	startTime := time.Now()

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
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionOAuthRevoke,
			TargetType:   "token",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		c.Error(app_error.ErrBadRequest)
		return
	}

	// 记录令牌撤销成功
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionOAuthRevoke,
		TargetType: "token",
		Detail: &models.AuditLogDetail{
			Metadata: map[string]any{
				"token_type_hint": req.TokenTypeHint,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	c.Status(200)
}

// CreateClient 创建 OAuth2 客户端
// POST /oauth/clients
func (h *OAuthHandler) CreateClient(c *gin.Context) {
	startTime := time.Now()

	if err := services.VerifyPermissions(c, h.roleService, []string{
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
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleClient,
			Action:       models.AuditActionClientCreate,
			TargetType:   "oauth_client",
			TargetName:   req.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录客户端创建成功
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientCreate,
		TargetID:   client.ID,
		TargetType: "oauth_client",
		TargetName: client.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"name":      client.Name,
				"domain":    client.Domain,
				"is_public": req.IsPublic,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{
		"client_id":     client.ID,
		"client_name":   client.Name,
		"client_domain": client.Domain,
	})
}

// GetClient 获取客户端信息
// GET /oauth/clients/:id
func (h *OAuthHandler) GetClient(c *gin.Context) {
	if err := services.VerifyPermissions(c, h.roleService, []string{
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
	startTime := time.Now()

	if err := services.VerifyPermissions(c, h.roleService, []string{
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

	// 获取旧的客户端信息
	oldClient, _ := h.oauthService.GetClientByID(clientID)

	if err := h.oauthService.UpdateClient(clientID, req.Name, req.Domain); err != nil {
		utilities.GetLogger().Error("Update OAuth client error", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleClient,
			Action:       models.AuditActionClientUpdate,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录客户端更新
	var oldValue map[string]any
	if oldClient != nil {
		oldValue = map[string]any{
			"name":   oldClient.Name,
			"domain": oldClient.Domain,
		}
	}
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientUpdate,
		TargetID:   clientID,
		TargetType: "oauth_client",
		TargetName: req.Name,
		Detail: &models.AuditLogDetail{
			OldValue: oldValue,
			NewValue: map[string]any{
				"name":   req.Name,
				"domain": req.Domain,
			},
			ChangedFields: []string{"name", "domain"},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"updated": true})
}

// DeleteClient 删除客户端
// DELETE /oauth/clients/:id
func (h *OAuthHandler) DeleteClient(c *gin.Context) {
	startTime := time.Now()

	if err := services.VerifyPermissions(c, h.roleService, []string{
		permissions.OAuthClientDelete,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	// 获取客户端信息用于审计
	client, _ := h.oauthService.GetClientByID(clientID)
	clientName := ""
	if client != nil {
		clientName = client.Name
	}

	if err := h.oauthService.DeleteClient(clientID); err != nil {
		utilities.GetLogger().Error("Delete OAuth client error", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleClient,
			Action:       models.AuditActionClientDelete,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			TargetName:   clientName,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录客户端删除
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientDelete,
		TargetID:   clientID,
		TargetType: "oauth_client",
		TargetName: clientName,
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"deleted": true})
}

// ListClients 获取客户端列表
// GET /oauth/clients
func (h *OAuthHandler) ListClients(c *gin.Context) {
	if err := services.VerifyPermissions(c, h.roleService, []string{
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
		response.HandlerError(c, app_error.ErrBadRequest)
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

// isURL 检查字符串是否为有效的 URL
func (h *OAuthHandler) isURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

// Logout 登出端点
func (h *OAuthHandler) Logout(c *gin.Context) {
	// 提取访问令牌
	token, err := jwt.ExtractTokenFromHeader(c)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	// 获取令牌信息
	info, err := h.oauthService.GetTokenInfo(token)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	// 撤销访问令牌和刷新令牌
	refreshToken := info.GetRefresh()
	h.oauthService.GetManager().RemoveAccessToken(c, token)
	h.oauthService.GetManager().RemoveRefreshToken(c, refreshToken)

	// 处理登出重定向，只允许重定向到相对路径或空路径
	postLogoutRedirectURI := c.Query("post_logout_redirect_uri")
	if postLogoutRedirectURI == "" || h.isURL(postLogoutRedirectURI) {
		postLogoutRedirectURI = "/"
	}

	c.Redirect(302, postLogoutRedirectURI)
}
