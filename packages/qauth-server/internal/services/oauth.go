package services

import (
	"context"
	"net/http"
	"qauth-server/internal/config"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/providers"
	"qauth-server/internal/repository"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwks"
	app_jwt "qauth-server/pkg/jwt"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"gorm.io/gorm"
)

// OAuthService OAuth2 服务
type OAuthService struct {
	cfg             *config.Config
	provider        providers.IOAuth
	jwksManager     *jwks.JWKSManager
	userService     *UserService
	appGroupService *AppGroupService
	oauthRepo       *repository.OAuthRepository
	loginStateRepo  *repository.LoginStateRepository
	errRecordRepo   *repository.ErrorRecordRepository
	logger          providers.ILogger
}

// SetAppGroupService 设置应用组服务（用于延迟注入以避免循环依赖）
func (s *OAuthService) SetAppGroupService(appGroupService *AppGroupService) {
	s.appGroupService = appGroupService
}

// NewOAuthService 创建新的 OAuth2 服务
func NewOAuthService(
	cfg *config.Config,
	oauth providers.IOAuth,
	jwksManager *jwks.JWKSManager,
	userService *UserService,
	oauthRepo *repository.OAuthRepository,
	loginStateRepo *repository.LoginStateRepository,
	errRecordRepo *repository.ErrorRecordRepository,
	logger providers.ILogger,
) *OAuthService {
	service := &OAuthService{
		cfg:            cfg,
		provider:       oauth,
		jwksManager:    jwksManager,
		userService:    userService,
		oauthRepo:      oauthRepo,
		loginStateRepo: loginStateRepo,
		errRecordRepo:  errRecordRepo,
		logger:         logger,
	}

	// 设置处理器到 provider
	oauth.SetUserAuthorizationHandler(service.userAuthorizationHandler)
	oauth.SetPasswordAuthorizationHandler(service.passwordAuthorizationHandler)
	oauth.SetExtensionFieldsHandler(service.extensionFieldsHandler)
	oauth.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		service.oauthErrorHandler("", nil, err.Error())
		return
	})
	oauth.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Error("oauth response error", "error", re.Error.Error())
	})

	return service
}

// GetConfig 获取 OAuth 配置
func (s *OAuthService) GetConfig() *config.Config {
	return s.cfg
}

// oauthErrorHandler OAuth 错误处理器
func (s *OAuthService) oauthErrorHandler(userID string, clientID *string, errMsg string) {
	record := &models.ErrorRecord{
		UserID:    userID,
		ClientID:  clientID,
		ErrorType: "OAuthProvider",
		Message:   errMsg,
		Timestamp: time.Now().Unix(),
	}
	if err := s.errRecordRepo.Create(record); err != nil {
		s.logger.Error("failed to record oauth error", "error", err)
	}
}

// userAuthorizationHandler 用户授权处理器（授权码模式）
func (s *OAuthService) userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	clientID := r.FormValue("client_id")

	// 从 Cookie 中获取 JWT Token
	cookie, err := r.Cookie("access_token")
	if err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not logged in")
		return "", e.ErrUnauthorized.WithMessage("missing access token")
	}

	// 解析 JWT Token
	claims, err := app_jwt.ParseAccessToken(cookie.Value)
	if err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "invalid access token")
		return "", e.ErrInvalidToken.Wrap(err)
	}

	// 验证用户是否存在
	user, err := s.userService.GetUserByID(claims.UserID, false)
	if err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not found")
		if err == gorm.ErrRecordNotFound {
			return "", e.ErrUserNotFound.WithMessage(claims.UserID)
		}
		return "", e.ErrUserNotFound.Wrap(err)
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "user account is not active")
		return "", e.ErrUserInactive.WithMessage(user.ID)
	}

	s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, true, "")
	return user.ID, nil
}

// passwordAuthorizationHandler 密码授权处理器（密码模式）
func (s *OAuthService) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	// 根据学号查找用户
	user, err := s.userService.GetUserByStudentID(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not found")
			return "", e.ErrUserNotFound.WithMessage(username)
		}
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "database error: "+err.Error())
		return "", e.ErrUserNotFound.Wrap(err)
	}

	// 验证密码
	if !utilities.VerifyPassword(password, user.Salt, user.PasswordHash) {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "invalid password")
		return "", e.ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "user account is not active")
		return "", e.ErrUserInactive.WithMessage(user.ID)
	}

	// 记录登录状态
	s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, true, "")
	return user.ID, nil
}

// extensionFieldsHandler 扩展字段处理器（用于生成 ID Token）
func (s *OAuthService) extensionFieldsHandler(ti providers.ITokenInfo) (fieldsValue map[string]any) {
	scopes := strings.Split(ti.GetScope(), " ")
	if !slices.Contains(scopes, string(config.ScopeOpenID)) {
		return nil
	}

	user, err := s.userService.GetUserByID(ti.GetUserID(), false)
	if err != nil {
		clientID := ti.GetClientID()
		s.recordError(ti.GetUserID(), &clientID, e.ErrUserNotFound.Wrap(err).Error())
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
		s.recordError(ti.GetUserID(), &clientID, e.ErrIDTokenGenerationFailed.Wrap(err).Error())
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
	if err := s.loginStateRepo.Create(loginState); err != nil {
		utilities.GetLogger().Error("failed to record login state", "error", err)
	}
}

// recordError 记录错误信息
func (s *OAuthService) recordError(userID string, clientID *string, errMsg string) {
	record := &models.ErrorRecord{
		UserID:    userID,
		ClientID:  clientID,
		ErrorType: "OAuthService",
		Message:   errMsg,
		Timestamp: time.Now().Unix(),
	}
	if err := s.errRecordRepo.Create(record); err != nil {
		utilities.GetLogger().Error("failed to record error", "error", err)
	}
}

// GetServer 获取 OAuth2 服务器实例（已弃用，使用 provider 接口）
func (s *OAuthService) GetServer() providers.IOAuth {
	return s.provider
}

// GetManager 获取 OAuth2 管理器实例（已弃用，使用 provider 接口）
func (s *OAuthService) GetManager() providers.IOAuth {
	return s.provider
}

// HandleAuthorizeRequest 处理授权请求
func (s *OAuthService) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.provider.HandleAuthorizeRequest(w, r)
}

// HandlerAuthorizeDeny 处理授权拒绝
func (s *OAuthService) HandlerAuthorizeDeny(c *gin.Context, redirectUri string) {
	url := redirectUri + "?error=access_denied&error_description=The+user+denied+the+request"
	c.Redirect(http.StatusFound, url)
}

// ValidateAuthorizeRequest 验证授权请求
func (s *OAuthService) ValidateAuthorizeRequest(r *http.Request) (*providers.AuthorizeRequest, error) {
	return s.provider.ValidateAuthorizeRequest(r)
}

// HandleTokenRequest 处理令牌请求
func (s *OAuthService) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.provider.HandleTokenRequest(w, r)
}

// ValidateToken 验证访问令牌
func (s *OAuthService) ValidateToken(r *http.Request) (providers.ITokenInfo, error) {
	return s.provider.ValidateToken(r)
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
		return e.ErrNoFieldsToUpdate
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
	return s.provider.RevokeToken(context.Background(), accessToken)
}

// RevokeRefreshToken 撤销刷新令牌
func (s *OAuthService) RevokeRefreshToken(refreshToken string) error {
	return s.provider.RevokeRefreshToken(context.Background(), refreshToken)
}

// GetTokenInfo 获取令牌信息
func (s *OAuthService) GetTokenInfo(accessToken string) (providers.ITokenInfo, error) {
	return s.provider.GetTokenInfo(context.Background(), accessToken)
}
