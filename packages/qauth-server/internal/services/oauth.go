package services

import (
	"context"
	"net/http"
	"qauth-server/internal/config"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/repository"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwks"
	app_jwt "qauth-server/pkg/jwt"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	oautherrors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauthmodels "github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// OAuthService OAuth2 服务
type OAuthService struct {
	db              *gorm.DB
	cfg             *config.Config
	server          *server.Server
	manager         *manage.Manager
	clientStore     *GormClientStore
	jwksManager     *jwks.JWKSManager
	userService     *UserService
	appGroupService *AppGroupService
	oauthRepo       *repository.OAuthRepository
}

// SetAppGroupService 设置应用组服务（用于延迟注入以避免循环依赖）
func (s *OAuthService) SetAppGroupService(appGroupService *AppGroupService) {
	s.appGroupService = appGroupService
}

// GormClientStore 基于 GORM 的 OAuth2 客户端存储
type GormClientStore struct {
	db *gorm.DB
}

// NewGormClientStore 创建新的 GORM 客户端存储
func NewGormClientStore(db *gorm.DB) *GormClientStore {
	return &GormClientStore{db: db}
}

// GetByID 根据 ID 获取 OAuth2 客户端
func (s *GormClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var client models.OAuth2Client
	if err := s.db.WithContext(ctx).First(&client, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, app_error.ErrClientNotFound.WithMessage(id)
		}
		return nil, app_error.ErrClientQueryFailed.Wrap(err)
	}

	return &oauthmodels.Client{
		ID:     client.ID,
		Secret: client.Secret,
		Domain: client.Domain,
		UserID: client.Data.UserID,
		Public: client.Data.Public,
	}, nil
}

// CustomAccessTokenGenerate 自定义令牌生成器
type CustomAccessTokenGenerate struct {
}

func (g *CustomAccessTokenGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	access, err = utilities.GenerateRandomString(32)
	if isGenRefresh {
		refresh, err = utilities.GenerateRandomString(32)
	}

	if err != nil {
		return "", "", app_error.ErrTokenGenerationFailed.Wrap(err)
	}

	access = "atk:" + access
	refresh = "rtk:" + refresh

	return access, refresh, err
}

// NewOAuthService 创建新的 OAuth2 服务
func NewOAuthService(
	db *gorm.DB,
	cfg *config.Config,
	jwksManager *jwks.JWKSManager,
	userService *UserService,
	oauthRepo *repository.OAuthRepository,
) *OAuthService {
	logger := utilities.GetLogger()

	// 创建基于 GORM 的客户端存储
	clientStore := NewGormClientStore(db)

	// 创建 Redis Token 存储
	tokenStore := oredis.NewRedisStore(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, "oauth2_token:")

	// 创建 OAuth2 管理器
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	// 使用自定义的令牌生成器
	manager.MapAccessGenerate(&CustomAccessTokenGenerate{})

	// 配置授权码模式
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    cfg.JWT.AccessTokenExp,
		RefreshTokenExp:   cfg.JWT.RefreshTokenExp,
		IsGenerateRefresh: true,
	})

	// 配置密码模式
	manager.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp:    cfg.JWT.AccessTokenExp,
		RefreshTokenExp:   cfg.JWT.RefreshTokenExp,
		IsGenerateRefresh: true,
	})

	// 配置客户端凭证模式
	manager.SetClientTokenCfg(&manage.Config{
		AccessTokenExp: cfg.JWT.AccessTokenExp,
	})

	// 配置刷新令牌
	manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
		IsGenerateRefresh:  true,
		IsRemoveAccess:     true,
		IsRemoveRefreshing: true,
	})

	// 创建 OAuth2 服务器
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	service := &OAuthService{
		db:          db,
		cfg:         cfg,
		server:      srv,
		manager:     manager,
		clientStore: clientStore,
		jwksManager: jwksManager,
		userService: userService,
		oauthRepo:   oauthRepo,
	}

	// 设置内部错误处理
	srv.SetInternalErrorHandler(func(err error) (re *oautherrors.Response) {
		service.recordError("", nil, err.Error())
		logger.Error("OAuth2 internal error", "error", err)
		return
	})

	// 设置响应错误处理
	srv.SetResponseErrorHandler(func(re *oautherrors.Response) {
		logger.Error("OAuth2 response error", "error", re.Error, "description", re.Description)
	})

	// 仅允许授权码模式
	gt := []oauth2.GrantType{}
	for _, grantType := range cfg.OAuth.GrantTypesSupported {
		gt = append(gt, oauth2.GrantType(grantType))
	}
	srv.SetAllowedGrantType(gt...)

	// 设置用户授权处理器
	srv.SetUserAuthorizationHandler(service.userAuthorizationHandler)

	// 设置密码授权处理器
	srv.SetPasswordAuthorizationHandler(service.passwordAuthorizationHandler)

	// 设置扩展字段处理器
	srv.SetExtensionFieldsHandler(service.extensionFieldsHandler)

	logger.Info("OAuth2 service initialized")
	return service
}

// GetConfig 获取 OAuth 配置
func (s *OAuthService) GetConfig() *config.Config {
	return s.cfg
}

// userAuthorizationHandler 用户授权处理器（授权码模式）
func (s *OAuthService) userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	clientID := r.FormValue("client_id")

	// 从 Cookie 中获取 JWT Token
	cookie, err := r.Cookie("access_token")
	if err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not logged in")
		return "", app_error.ErrUnauthorized.WithMessage("missing access token")
	}

	// 解析 JWT Token
	claims, err := app_jwt.ParseAccessToken(cookie.Value)
	if err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "invalid access token")
		return "", app_error.ErrInvalidToken.Wrap(err)
	}

	// 验证用户是否存在
	var user models.Users
	if err := s.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not found")
		if err == gorm.ErrRecordNotFound {
			return "", app_error.ErrUserNotFound.WithMessage(claims.UserID)
		}
		return "", app_error.ErrUserNotFound.Wrap(err)
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "user account is not active")
		return "", app_error.ErrUserInactive.WithMessage(user.ID)
	}

	s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, true, "")
	return user.ID, nil
}

// passwordAuthorizationHandler 密码授权处理器（密码模式）
func (s *OAuthService) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	// 根据学号查找用户
	var user models.Users
	if err := s.db.First(&user, "student_id = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not found")
			return "", app_error.ErrUserNotFound.WithMessage(username)
		}
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "database error: "+err.Error())
		return "", app_error.ErrUserNotFound.Wrap(err)
	}

	// 验证密码
	if !utilities.VerifyPassword(password, user.Salt, user.PasswordHash) {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "invalid password")
		return "", app_error.ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "user account is not active")
		return "", app_error.ErrUserInactive.WithMessage(user.ID)
	}

	// 记录登录状态
	s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, true, "")
	return user.ID, nil
}

// extensionFieldsHandler 扩展字段处理器（用于生成 ID Token）
func (s *OAuthService) extensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]any) {
	scopes := strings.Split(ti.GetScope(), " ")
	if !slices.Contains(scopes, string(config.ScopeOpenID)) {
		return nil
	}

	user, err := s.userService.GetUserByID(ti.GetUserID(), false)
	if err != nil {
		clientID := ti.GetClientID()
		s.recordError(ti.GetUserID(), &clientID, app_error.ErrUserNotFound.Wrap(err).Error())
		return nil
	}

	if user.DisplayName == nil {
		name := user.Name
		user.DisplayName = &name
	}

	basic := &jwks.BasicClaims{
		Issuer:   s.cfg.OIDC.Issuer,
		Subject:  user.ID,
		Audience: ti.GetClientID(),
		Expiry:   ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Unix(),
		IssuedAt: ti.GetAccessCreateAt().Unix(),
		AuthTime: ti.GetAccessCreateAt().Unix(),
	}
	data := map[string]any{}

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

	// 添加应用组权限（permissions scope）
	if slices.Contains(scopes, string(config.ScopePermissions)) && s.appGroupService != nil {
		clientID := ti.GetClientID()
		userID := ti.GetUserID()

		// 获取用户在该应用组的角色
		appGroupRoles, err := s.appGroupService.GetUserAppGroupRoles(userID, clientID)
		if err == nil && len(appGroupRoles) > 0 {
			roleNames := make([]string, len(appGroupRoles))
			for i, role := range appGroupRoles {
				roleNames[i] = role.Code
			}
			data[string(config.ClaimAppGroupRoles)] = roleNames
		}

		// 获取用户在该应用组的权限
		appGroupPermissions, err := s.appGroupService.GetUserAppGroupPermissions(userID, clientID)
		if err == nil && len(appGroupPermissions) > 0 {
			permCodes := make([]string, len(appGroupPermissions))
			for i, perm := range appGroupPermissions {
				permCodes[i] = perm.Code
			}
			data[string(config.ClaimAppGroupPermissions)] = permCodes
		}
	}

	// 生成 ID Token
	idToken, err := s.jwksManager.SignToken(basic, data)
	if err != nil {
		clientID := ti.GetClientID()
		s.recordError(ti.GetUserID(), &clientID, app_error.ErrIDTokenGenerationFailed.Wrap(err).Error())
		return nil
	}

	return map[string]any{"id_token": idToken}
}

// recordLoginState 记录登录状态
func (s *OAuthService) recordLoginState(userID string, clientID *string, loginType models.LoginType, isSuccess bool, failReason string) {
	loginState := &models.LoginState{
		UserID:     userID,
		ClientID:   clientID,
		Type:       loginType,
		IsSuccess:  isSuccess,
		FailReason: failReason,
	}
	if err := s.oauthRepo.CreateLoginState(loginState); err != nil {
		utilities.GetLogger().Error("failed to record login state", "error", err)
	}
}

// recordError 记录错误信息
func (s *OAuthService) recordError(userID string, clientID *string, errMsg string) {
	errorRecord := &models.ErrorRecord{
		UserID:    userID,
		ClientID:  clientID,
		ErrorType: "OAuthService",
		Message:   errMsg,
		Timestamp: time.Now().Unix(),
	}
	if err := s.oauthRepo.CreateErrorRecord(errorRecord); err != nil {
		utilities.GetLogger().Error("failed to record error", "error", err)
	}
}

// GetServer 获取 OAuth2 服务器实例
func (s *OAuthService) GetServer() *server.Server {
	return s.server
}

// GetManager 获取 OAuth2 管理器实例
func (s *OAuthService) GetManager() *manage.Manager {
	return s.manager
}

// HandleAuthorizeRequest 处理授权请求
func (s *OAuthService) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.server.HandleAuthorizeRequest(w, r)
}

// HandlerAuthorizeDeny 处理授权拒绝
func (s *OAuthService) HandlerAuthorizeDeny(c *gin.Context, redirectUri string) {
	url := redirectUri + "?error=access_denied&error_description=The+user+denied+the+request"
	c.Redirect(http.StatusFound, url)
}

// ValidateAuthorizeRequest 验证授权请求
func (s *OAuthService) ValidateAuthorizeRequest(r *http.Request) (*server.AuthorizeRequest, error) {
	return s.server.ValidationAuthorizeRequest(r)
}

// HandleTokenRequest 处理令牌请求
func (s *OAuthService) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.server.HandleTokenRequest(w, r)
}

// ValidateToken 验证访问令牌
func (s *OAuthService) ValidateToken(r *http.Request) (oauth2.TokenInfo, error) {
	return s.server.ValidationBearerToken(r)
}

// CreateClient 创建 OAuth2 客户端
func (s *OAuthService) CreateClient(name, domain, secret string, isPublic bool, userID string) (*models.OAuth2Client, error) {
	client := &models.OAuth2Client{
		Name:   name,
		Domain: domain,
		Secret: secret,
		Data: models.ClientData{
			Domain: domain,
			Public: isPublic,
			UserID: userID,
		},
	}

	if err := s.oauthRepo.CreateClient(client); err != nil {
		return nil, err
	}

	// 更新 Data 中的 ID
	client.Data.ID = client.ID
	client.Data.Secret = secret
	if err := s.oauthRepo.UpdateClient(client); err != nil {
		return nil, err
	}

	return client, nil
}

// CreateClientParams 创建客户端参数
type CreateClientParams struct {
	Name         string
	Description  string
	Domain       string
	RedirectURIs []string
	Scopes       []string
	GrantTypes   []string
	Status       models.ClientStatus
	Trusted      bool
	Logo         string
	Icon         string
	IconBg       string
	Secret       string
	IsPublic     bool
	UserID       string
}

// CreateClientFull 创建完整的 OAuth2 客户端
func (s *OAuthService) CreateClientFull(params *CreateClientParams) (*models.OAuth2Client, error) {
	// 设置默认值
	if params.Status == "" {
		params.Status = models.ClientStatusDevelopment
	}
	if params.Icon == "" {
		params.Icon = "pi pi-box"
	}
	if params.IconBg == "" {
		params.IconBg = "linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)"
	}
	if len(params.GrantTypes) == 0 {
		params.GrantTypes = []string{"authorization_code", "refresh_token"}
	}
	if len(params.Scopes) == 0 {
		params.Scopes = []string{"openid", "profile"}
	}

	client := &models.OAuth2Client{
		Name:         params.Name,
		Description:  params.Description,
		Domain:       params.Domain,
		Secret:       params.Secret,
		RedirectURIs: params.RedirectURIs,
		Scopes:       params.Scopes,
		GrantTypes:   params.GrantTypes,
		Status:       params.Status,
		Trusted:      params.Trusted,
		Logo:         params.Logo,
		Icon:         params.Icon,
		IconBg:       params.IconBg,
		Data: models.ClientData{
			Domain: params.Domain,
			Public: params.IsPublic,
			UserID: params.UserID,
		},
	}

	if err := s.oauthRepo.CreateClient(client); err != nil {
		return nil, err
	}

	// 更新 Data 中的 ID
	client.Data.ID = client.ID
	client.Data.Secret = params.Secret
	if err := s.oauthRepo.UpdateClient(client); err != nil {
		return nil, err
	}

	return client, nil
}

// GetClientByID 根据 ID 获取客户端
func (s *OAuthService) GetClientByID(clientID string) (*models.OAuth2Client, error) {
	return s.oauthRepo.FindClientByID(context.Background(), clientID)
}

// UpdateClientParams 更新客户端参数
type UpdateClientParams struct {
	Name         *string
	Description  *string
	Domain       *string
	RedirectURIs []string
	Scopes       []string
	GrantTypes   []string
	Status       *models.ClientStatus
	Trusted      *bool
	Logo         *string
	Icon         *string
	IconBg       *string
}

// UpdateClient 更新客户端信息
func (s *OAuthService) UpdateClient(clientID string, name, domain string) error {
	return s.oauthRepo.UpdateClientFields(clientID, map[string]interface{}{
		"name":   name,
		"domain": domain,
	})
}

// UpdateClientFull 更新完整的客户端信息
func (s *OAuthService) UpdateClientFull(clientID string, params *UpdateClientParams) error {
	updates := make(map[string]interface{})

	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.Description != nil {
		updates["description"] = *params.Description
	}
	if params.Domain != nil {
		updates["domain"] = *params.Domain
	}
	if params.RedirectURIs != nil {
		updates["redirect_uris"] = models.StringArray(params.RedirectURIs)
	}
	if params.Scopes != nil {
		updates["scopes"] = models.StringArray(params.Scopes)
	}
	if params.GrantTypes != nil {
		updates["grant_types"] = models.StringArray(params.GrantTypes)
	}
	if params.Status != nil {
		updates["status"] = *params.Status
	}
	if params.Trusted != nil {
		updates["trusted"] = *params.Trusted
	}
	if params.Logo != nil {
		updates["logo"] = *params.Logo
	}
	if params.Icon != nil {
		updates["icon"] = *params.Icon
	}
	if params.IconBg != nil {
		updates["icon_bg"] = *params.IconBg
	}

	if len(updates) == 0 {
		return app_error.ErrNoFieldsToUpdate
	}

	return s.oauthRepo.UpdateClientFields(clientID, updates)
}

// DeleteClient 删除客户端
func (s *OAuthService) DeleteClient(clientID string) error {
	return s.oauthRepo.DeleteClient(clientID)
}

// ListClientsParams 列出客户端参数
type ListClientsParams struct {
	Page     int
	PageSize int
	Search   string
	Status   string
	SortBy   string
	SortDesc bool
}

// CountClients 计算客户端数量
func (s *OAuthService) CountClients() (int64, error) {
	return s.oauthRepo.CountClients()
}

// ListClients 获取客户端列表
func (s *OAuthService) ListClients(page, pageSize int) ([]models.OAuth2Client, int64, error) {
	return s.oauthRepo.ListClients(page, pageSize)
}

// ListClientsFull 获取客户端列表（支持更多查询参数）
func (s *OAuthService) ListClientsFull(params *ListClientsParams) ([]models.OAuth2Client, int64, error) {
	return s.oauthRepo.ListClientsWithFilter(
		params.Search,
		params.Status,
		params.SortBy,
		params.SortDesc,
		params.Page,
		params.PageSize,
	)
}

// RegenerateClientSecret 重新生成客户端密钥
func (s *OAuthService) RegenerateClientSecret(clientID, newSecret string) error {
	return s.oauthRepo.UpdateClientSecret(clientID, newSecret)
}

// IncrementClientRequestCount 增加客户端请求计数
func (s *OAuthService) IncrementClientRequestCount(clientID string) error {
	return s.oauthRepo.IncrementClientRequestCount(clientID)
}

// GetClientStats 获取客户端统计数据
func (s *OAuthService) GetClientStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 获取总客户端数
	total, err := s.oauthRepo.CountClients()
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// 获取各状态客户端数
	statusCounts, err := s.oauthRepo.CountClientsByStatus()
	if err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		stats[string(sc.Status)] = sc.Count
	}

	return stats, nil
}

// RevokeToken 撤销令牌
func (s *OAuthService) RevokeToken(accessToken string) error {
	if err := s.manager.RemoveAccessToken(context.Background(), accessToken); err != nil {
		return app_error.ErrTokenRevocationFailed.Wrap(err)
	}
	return nil
}

// RevokeRefreshToken 撤销刷新令牌
func (s *OAuthService) RevokeRefreshToken(refreshToken string) error {
	if err := s.manager.RemoveRefreshToken(context.Background(), refreshToken); err != nil {
		return app_error.ErrTokenRevocationFailed.Wrap(err)
	}
	return nil
}

// GetTokenInfo 获取令牌信息
func (s *OAuthService) GetTokenInfo(accessToken string) (oauth2.TokenInfo, error) {
	info, err := s.manager.LoadAccessToken(context.Background(), accessToken)
	if err != nil {
		return nil, app_error.ErrTokenLoadFailed.Wrap(err)
	}
	return info, nil
}
