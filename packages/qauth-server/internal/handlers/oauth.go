package handlers

import (
	"encoding/json"
	"net/http"
	"qauth-server/internal/config"
	"qauth-server/internal/config/permissions"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/providers"
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
	cacheSrv    providers.ICache
	oauthSrv    *services.OAuthService
	roleSrv     *services.RoleService
	userSrv     *services.UserService
	oidcSrv     *services.OIDCService
	auditSrv    *services.AuditService
	appGroupSrv *services.AppGroupService
	logger      providers.ILogger
}

// NewOAuthHandler 创建新的 OAuth2 处理器
func NewOAuthHandler(
	cache providers.ICache,
	oauthSrv *services.OAuthService,
	roleSrv *services.RoleService,
	userSrv *services.UserService,
	oidcSrv *services.OIDCService,
	auditSrv *services.AuditService,
	appGroupSrv *services.AppGroupService,
	logger providers.ILogger,
) *OAuthHandler {
	return &OAuthHandler{
		oauthSrv:    oauthSrv,
		roleSrv:     roleSrv,
		userSrv:     userSrv,
		oidcSrv:     oidcSrv,
		cacheSrv:    cache,
		auditSrv:    auditSrv,
		appGroupSrv: appGroupSrv,
		logger:      logger.With("handler", "OAuthHandler"),
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
	_, err := h.oauthSrv.ValidateAuthorizeRequest(c.Request)
	if err != nil {
		h.logger.Error("Validate authorize request error", "error", err)
		c.Error(app_error.ErrBadRequest)
		return
	}

	// 将请求查询参数转换为 JSON 字符串
	jsonStr, err := h.requestQueryToJSON(c.Request)
	if err != nil {
		h.logger.Error("Convert request query to JSON error", "error", err)
		c.Error(app_error.ErrInternalServerError)
		return
	}

	// 生成 authorizeID
	authorizeID, err := utilities.GenerateRandomString(32)
	if err != nil {
		h.logger.Error("Generate authorize ID error", "error", err)
		c.Error(app_error.ErrInternalServerError)
		return
	}

	// 将 authorizeID 和请求参数存储在缓存中，有效期为 60 秒
	h.cacheSrv.SetKeyValue("authorize-"+authorizeID, jsonStr, 120)

	// 重定向到授权页面
	authorizePageUrl := h.oauthSrv.GetConfig().OAuth.AuthorizePageURL + "?aid=" + authorizeID
	c.Redirect(302, authorizePageUrl)
}

// AuthorizeInfo 获取授权信息
// GET /oauth/authorize/info
func (h *OAuthHandler) AuthorizeInfo(c *gin.Context) {
	r := c.Request

	// TODO: 验证 referer 是否来自授权页面

	// 从缓存中获取请求参数
	authorizeID := r.URL.Query().Get("aid")
	data, err := h.cacheSrv.GetKeyValue("authorize-" + authorizeID)
	if err != nil || data == "" {
		h.logger.Error("Get authorize info from cache error", "error", err)
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 将 JSON 字符串解析为 map
	var params map[string]any
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		h.logger.Error("Parse authorize info JSON error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, params)
}

// verifyAuthorizeID 验证 authorizeID 是否有效且请求 URI 一致
func (h *OAuthHandler) restoreRequestQuery(c *gin.Context, aid string) error {
	// 从缓存中获取请求参数
	data, err := h.cacheSrv.GetKeyValue("authorize-" + aid)
	if err != nil || data == "" {
		h.logger.Error("Get authorize info from cache error", "error", err)
		return app_error.ErrBadRequest
	}

	// 将 JSON 字符串解析为 map
	var params map[string]any
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		h.logger.Error("Parse authorize info JSON error", "error", err)
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
		h.auditSrv.LogOAuthAuthorize(c, "", "", &clientID, scopes, "", responseType, redirectUri, false, "user denied")
		h.oauthSrv.HandlerAuthorizeDeny(c, redirectUri)
		return
	}

	// 处理非法请求
	if action != "approve" {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	c.Request.RequestURI = ""

	// 删除缓存中的授权请求
	h.cacheSrv.DeleteKey("authorize-" + authorizeID)

	// 处理授权请求
	rt := c.Request.URL.Query().Get("response_type")
	arts := h.oauthSrv.GetConfig().OAuth.AllowedResponseTypes
	isValid := slices.Contains(arts, config.ResponseType(rt))

	// 根据响应类型修改请求以支持隐式授权模式
	h.modifyResponseType(c, isValid, rt)

	// 处理授权请求
	err := h.oauthSrv.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		h.logger.Error("OAuth authorize error", "error", err)
		// 记录授权失败
		h.auditSrv.LogOAuthAuthorize(c, "", "", &clientID, scopes, "", responseType, redirectUri, false, err.Error())
		c.Error(app_error.ErrBadRequest)
		return
	}

	// 记录授权成功
	userID := ""
	userName := ""
	if uid, exists := c.Get("userID"); exists {
		userID = uid.(string)
		if user, err := h.userSrv.GetUserByID(userID, false); err == nil {
			userName = user.Name
		}
	}
	var operatorID *string
	if userID != "" {
		operatorID = &userID
	}
	h.auditSrv.Log(&services.AuditContext{
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
	h.cacheSrv.IncrKey("todays-authcount")

	// 根据响应类型修改重定向位置以支持隐式授权模式
	h.modifyRedirectLocation(c, isValid, rt)
}

// Token 处理令牌请求
// POST /oauth/token
func (h *OAuthHandler) Token(c *gin.Context) {
	startTime := time.Now()
	clientID := c.PostForm("client_id")
	grantType := c.PostForm("grant_type")

	err := h.oauthSrv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		h.logger.Error("OAuth token error", "error", err)
		h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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
	h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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
	tokenInfo, err := h.oauthSrv.ValidateToken(c.Request)
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
		err = h.oauthSrv.RevokeRefreshToken(req.Token)
	} else {
		err = h.oauthSrv.RevokeToken(req.Token)
	}

	if err != nil {
		h.logger.Error("OAuth revoke error", "error", err)
		h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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
	h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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

	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientCreate,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		Name         string   `json:"name" binding:"required"`
		Description  string   `json:"description"`
		Domain       string   `json:"domain" binding:"required"`
		RedirectURIs []string `json:"redirect_uris"`
		Scopes       []string `json:"scopes"`
		GrantTypes   []string `json:"grant_types"`
		Status       string   `json:"status"`
		Trusted      bool     `json:"trusted"`
		Logo         string   `json:"logo"`
		Icon         string   `json:"icon"`
		IconBg       string   `json:"icon_bg"`
		IsPublic     bool     `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 生成安全的客户端密钥
	secret, err := utilities.GenerateRandomString(32)
	if err != nil {
		h.logger.Error("Generate client secret error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 获取当前用户 ID
	userID := ""
	if userInfo, exists := c.Get("userInfo"); exists {
		userID = userInfo.(*jwt.UserJWTClaims).UserID
	}

	client, err := h.oauthSrv.CreateClientFull(&services.CreateClientParams{
		Name:         req.Name,
		Description:  req.Description,
		Domain:       req.Domain,
		RedirectURIs: req.RedirectURIs,
		Scopes:       req.Scopes,
		GrantTypes:   req.GrantTypes,
		Status:       models.ClientStatus(req.Status),
		Trusted:      req.Trusted,
		Logo:         req.Logo,
		Icon:         req.Icon,
		IconBg:       req.IconBg,
		Secret:       secret,
		IsPublic:     req.IsPublic,
		UserID:       userID,
	})
	if err != nil {
		h.logger.Error("Create OAuth client error", "error", err)
		h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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

	// 初始化应用组管理员（将创建者设为 owner）
	if h.appGroupSrv != nil && userID != "" {
		if err := h.appGroupSrv.InitializeAppGroupAdmins(client.ID, userID); err != nil {
			h.logger.Error("Initialize app group admins error", "error", err, "clientID", client.ID)
			// 不因此失败整个创建流程，只记录警告
		}
	}

	// 记录客户端创建成功
	h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientCreate,
		TargetID:   client.ID,
		TargetType: "oauth_client",
		TargetName: client.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"name":          client.Name,
				"description":   client.Description,
				"domain":        client.Domain,
				"redirect_uris": client.RedirectURIs,
				"scopes":        client.Scopes,
				"grant_types":   client.GrantTypes,
				"status":        client.Status,
				"trusted":       client.Trusted,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{
		"client":        client.ToResponse(),
		"client_secret": secret,
	})
}

// GetClient 获取客户端信息
// GET /oauth/clients/:id
func (h *OAuthHandler) GetClient(c *gin.Context) {
	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientView,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	client, err := h.oauthSrv.GetClientByID(clientID)
	if err != nil {
		response.HandlerError(c, app_error.ErrNotFound)
		return
	}

	response.HandlerSuccess(c, client.ToResponse())
}

// UpdateClient 更新客户端信息
// PUT /oauth/clients/:id
func (h *OAuthHandler) UpdateClient(c *gin.Context) {
	startTime := time.Now()

	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientUpdate,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	var req struct {
		Name         *string  `json:"name"`
		Description  *string  `json:"description"`
		Domain       *string  `json:"domain"`
		RedirectURIs []string `json:"redirect_uris"`
		Scopes       []string `json:"scopes"`
		GrantTypes   []string `json:"grant_types"`
		Status       *string  `json:"status"`
		Trusted      *bool    `json:"trusted"`
		Logo         *string  `json:"logo"`
		Icon         *string  `json:"icon"`
		IconBg       *string  `json:"icon_bg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 获取旧的客户端信息
	oldClient, _ := h.oauthSrv.GetClientByID(clientID)

	// 构建更新参数
	params := &services.UpdateClientParams{
		Name:         req.Name,
		Description:  req.Description,
		Domain:       req.Domain,
		RedirectURIs: req.RedirectURIs,
		Scopes:       req.Scopes,
		GrantTypes:   req.GrantTypes,
		Trusted:      req.Trusted,
		Logo:         req.Logo,
		Icon:         req.Icon,
		IconBg:       req.IconBg,
	}

	if req.Status != nil {
		status := models.ClientStatus(*req.Status)
		params.Status = &status
	}

	if err := h.oauthSrv.UpdateClientFull(clientID, params); err != nil {
		h.logger.Error("Update OAuth client error", "error", err)
		h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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

	// 获取更新后的客户端
	updatedClient, _ := h.oauthSrv.GetClientByID(clientID)

	// 记录变更字段
	changedFields := []string{}
	if req.Name != nil {
		changedFields = append(changedFields, "name")
	}
	if req.Description != nil {
		changedFields = append(changedFields, "description")
	}
	if req.Domain != nil {
		changedFields = append(changedFields, "domain")
	}
	if req.RedirectURIs != nil {
		changedFields = append(changedFields, "redirect_uris")
	}
	if req.Scopes != nil {
		changedFields = append(changedFields, "scopes")
	}
	if req.GrantTypes != nil {
		changedFields = append(changedFields, "grant_types")
	}
	if req.Status != nil {
		changedFields = append(changedFields, "status")
	}
	if req.Trusted != nil {
		changedFields = append(changedFields, "trusted")
	}

	// 记录客户端更新
	var oldValue map[string]any
	if oldClient != nil {
		oldValue = map[string]any{
			"name":          oldClient.Name,
			"description":   oldClient.Description,
			"domain":        oldClient.Domain,
			"redirect_uris": oldClient.RedirectURIs,
			"scopes":        oldClient.Scopes,
			"grant_types":   oldClient.GrantTypes,
			"status":        oldClient.Status,
			"trusted":       oldClient.Trusted,
		}
	}

	targetName := ""
	if updatedClient != nil {
		targetName = updatedClient.Name
	}

	h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientUpdate,
		TargetID:   clientID,
		TargetType: "oauth_client",
		TargetName: targetName,
		Detail: &models.AuditLogDetail{
			OldValue: oldValue,
			NewValue: map[string]any{
				"name":          req.Name,
				"description":   req.Description,
				"domain":        req.Domain,
				"redirect_uris": req.RedirectURIs,
				"scopes":        req.Scopes,
				"grant_types":   req.GrantTypes,
				"status":        req.Status,
				"trusted":       req.Trusted,
			},
			ChangedFields: changedFields,
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	if updatedClient != nil {
		response.HandlerSuccess(c, updatedClient.ToResponse())
	} else {
		response.HandlerSuccess(c, gin.H{"updated": true})
	}
}

// DeleteClient 删除客户端
// DELETE /oauth/clients/:id
func (h *OAuthHandler) DeleteClient(c *gin.Context) {
	startTime := time.Now()

	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientDelete,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	// 获取客户端信息用于审计
	client, _ := h.oauthSrv.GetClientByID(clientID)
	clientName := ""
	var deletedClientData map[string]any
	if client != nil {
		clientName = client.Name
		deletedClientData = map[string]any{
			"name":          client.Name,
			"description":   client.Description,
			"domain":        client.Domain,
			"redirect_uris": client.RedirectURIs,
			"scopes":        client.Scopes,
			"grant_types":   client.GrantTypes,
			"status":        client.Status,
			"trusted":       client.Trusted,
		}
	}

	if err := h.oauthSrv.DeleteClient(clientID); err != nil {
		h.logger.Error("Delete OAuth client error", "error", err)
		h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
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
	h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientDelete,
		TargetID:   clientID,
		TargetType: "oauth_client",
		TargetName: clientName,
		Detail: &models.AuditLogDetail{
			OldValue: deletedClientData,
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"deleted": true})
}

// ListClients 获取客户端列表
// GET /oauth/clients
func (h *OAuthHandler) ListClients(c *gin.Context) {
	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientList,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 解析查询参数
	params := &services.ListClientsParams{
		Page:     utilities.ParseIntParam(c.Query("page"), 1),
		PageSize: utilities.ParseIntParam(c.Query("page_size"), 10),
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		SortBy:   c.Query("sort_by"),
		SortDesc: c.Query("sort_desc") == "true",
	}

	// 限制 page_size 最大值
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	clients, total, err := h.oauthSrv.ListClientsFull(params)
	if err != nil {
		h.logger.Error("List OAuth clients error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	items := make([]*models.OAuth2ClientResponse, len(clients))
	for i, client := range clients {
		items[i] = client.ToResponse()
	}

	response.HandlerSuccess(c, gin.H{
		"items": items,
		"total": total,
		"page":  params.Page,
		"size":  params.PageSize,
	})
}

// RegenerateSecret 重新生成客户端密钥
// POST /oauth/clients/:id/regenerate-secret
func (h *OAuthHandler) RegenerateSecret(c *gin.Context) {
	startTime := time.Now()

	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientUpdate,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	clientID := c.Param("id")

	// 获取客户端信息
	client, err := h.oauthSrv.GetClientByID(clientID)
	if err != nil {
		response.HandlerError(c, app_error.ErrNotFound)
		return
	}

	// 生成新的密钥
	newSecret, err := utilities.GenerateRandomString(32)
	if err != nil {
		h.logger.Error("Generate client secret error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 更新密钥
	if err := h.oauthSrv.RegenerateClientSecret(clientID, newSecret); err != nil {
		h.logger.Error("Regenerate client secret error", "error", err)
		h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleClient,
			Action:       models.AuditActionClientUpdate,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			TargetName:   client.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录密钥重新生成
	h.auditSrv.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleClient,
		Action:     models.AuditActionClientUpdate,
		TargetID:   clientID,
		TargetType: "oauth_client",
		TargetName: client.Name,
		Detail: &models.AuditLogDetail{
			ChangedFields: []string{"secret"},
			Metadata: map[string]any{
				"action": "regenerate_secret",
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{
		"secret": newSecret,
	})
}

// GetClientStats 获取客户端统计
// GET /oauth/clients/stats
func (h *OAuthHandler) GetClientStats(c *gin.Context) {
	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientList,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	stats, err := h.oauthSrv.GetClientStats()
	if err != nil {
		h.logger.Error("Get client stats error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, stats)
}

// GetClientOptions 获取创建/编辑客户端时的可选配置项
// GET /oauth/clients/options
func (h *OAuthHandler) GetClientOptions(c *gin.Context) {
	if err := services.VerifyPermissions(c, h.roleSrv, []string{
		permissions.OAuthClientList,
	}); err != nil {
		response.HandlerError(c, err)
		return
	}

	oauthConfig := h.oauthSrv.GetConfig()

	// 构建授权范围选项
	scopeOptions := make([]map[string]string, 0, len(oauthConfig.OAuth.ScopeSupported))
	scopeLabels := map[config.Scope]string{
		config.ScopeOpenID:      "OpenID",
		config.ScopeProfile:     "Profile",
		config.ScopeEmail:       "Email",
		config.ScopeRoles:       "Roles",
		config.ScopePermissions: "Permissions",
	}
	for _, scope := range oauthConfig.OAuth.ScopeSupported {
		label := scopeLabels[scope]
		if label == "" {
			label = string(scope)
		}
		scopeOptions = append(scopeOptions, map[string]string{
			"label": label,
			"value": string(scope),
		})
	}

	// 构建授权类型选项
	grantTypeOptions := make([]map[string]string, 0, len(oauthConfig.OAuth.GrantTypesSupported))
	grantTypeLabels := map[config.GrantType]string{
		config.GrantTypeAuthorizationCode: "授权码",
		config.GrantTypeRefreshToken:      "刷新令牌",
		config.GrantTypeClientCredentials: "客户端凭证",
		config.GrantTypePassword:          "密码模式",
		config.GrantTypeImplicit:          "隐式授权",
	}
	for _, grantType := range oauthConfig.OAuth.GrantTypesSupported {
		label := grantTypeLabels[grantType]
		if label == "" {
			label = string(grantType)
		}
		grantTypeOptions = append(grantTypeOptions, map[string]string{
			"label": label,
			"value": string(grantType),
		})
	}

	response.HandlerSuccess(c, gin.H{
		"scopes":      scopeOptions,
		"grant_types": grantTypeOptions,
	})
}

// UserInfo 获取用户信息
func (h *OAuthHandler) UserInfo(c *gin.Context) {
	token, err := jwt.ExtractTokenFromHeader(c)
	if err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	info, err := h.oauthSrv.GetTokenInfo(token)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	userID := info.GetUserID()
	issueAt := info.GetAccessCreateAt().Unix()
	expiry := info.GetAccessExpiresIn().Seconds() + float64(issueAt)
	audience := info.GetClientID()
	issuer := h.oidcSrv.GetConfig().Issuer

	user, err := h.userSrv.GetUserByID(userID, true)
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
	info, err := h.oauthSrv.GetTokenInfo(token)
	if err != nil {
		response.HandlerError(c, app_error.ErrInvalidAccessToken)
		return
	}

	// 撤销访问令牌和刷新令牌
	refreshToken := info.GetRefresh()
	h.oauthSrv.RevokeToken(token)
	h.oauthSrv.RevokeRefreshToken(refreshToken)

	// 处理登出重定向，只允许重定向到相对路径或空路径
	postLogoutRedirectURI := c.Query("post_logout_redirect_uri")
	if postLogoutRedirectURI == "" || h.isURL(postLogoutRedirectURI) {
		postLogoutRedirectURI = "/"
	}

	c.Redirect(302, postLogoutRedirectURI)
}
