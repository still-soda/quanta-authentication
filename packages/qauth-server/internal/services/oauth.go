package services

import (
	"context"
	"errors"
	"net/http"
	"qauth-server/internal/config"
	"qauth-server/internal/models"
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
	db          *gorm.DB
	cfg         *config.Config
	server      *server.Server
	manager     *manage.Manager
	clientStore *GormClientStore
	jwksManager *jwks.JWKSManager
	userService *UserService
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
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
		return "", "", err
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
		return "", errors.New("user not logged in")
	}

	// 解析 JWT Token
	claims, err := app_jwt.ParseAccessToken(cookie.Value)
	if err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "invalid access token")
		return "", errors.New("invalid access token")
	}

	// 验证用户是否存在
	var user models.Users
	if err := s.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not found")
		return "", errors.New("user not found")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "user account is not active")
		return "", errors.New("user account is not active")
	}

	s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, true, "")
	return user.ID, nil
}

// passwordAuthorizationHandler 密码授权处理器（密码模式）
func (s *OAuthService) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	// 根据学号查找用户
	var user models.Users
	if err := s.db.First(&user, "student_id = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "user not found")
			return "", errors.New("user not found")
		}
		s.recordLoginState("", &clientID, models.LoginTypeOAuth2, false, "database error: "+err.Error())
		return "", err
	}

	// 验证密码
	if !utilities.VerifyPassword(password, user.Salt, user.PasswordHash) {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "invalid password")
		return "", errors.New("invalid password")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, false, "user account is not active")
		return "", errors.New("user account is not active")
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

	// 生成 ID Token
	idToken, err := s.jwksManager.SignToken(basic, data)
	if err != nil {
		utilities.GetLogger().Error("failed to sign ID token", "error", err)
		return nil
	}

	return map[string]any{"id_token": idToken}
}

// recordLoginState 记录登录状态
func (s *OAuthService) recordLoginState(userID string, clientID *string, loginType models.LoginType, isSuccess bool, failReason string) {
	loginState := models.LoginState{
		UserID:     userID,
		ClientID:   clientID,
		Type:       loginType,
		IsSuccess:  isSuccess,
		FailReason: failReason,
	}
	if err := s.db.Create(&loginState).Error; err != nil {
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
	if err := s.db.Create(errorRecord).Error; err != nil {
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

	if err := s.db.Create(client).Error; err != nil {
		return nil, err
	}

	// 更新 Data 中的 ID
	client.Data.ID = client.ID
	client.Data.Secret = secret
	if err := s.db.Save(client).Error; err != nil {
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

	if err := s.db.Create(client).Error; err != nil {
		return nil, err
	}

	// 更新 Data 中的 ID
	client.Data.ID = client.ID
	client.Data.Secret = params.Secret
	if err := s.db.Save(client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

// GetClientByID 根据 ID 获取客户端
func (s *OAuthService) GetClientByID(clientID string) (*models.OAuth2Client, error) {
	var client models.OAuth2Client
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
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
	return s.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Updates(map[string]interface{}{
			"name":   name,
			"domain": domain,
		}).Error
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
		return nil
	}

	return s.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Updates(updates).Error
}

// DeleteClient 删除客户端
func (s *OAuthService) DeleteClient(clientID string) error {
	return s.db.Delete(&models.OAuth2Client{}, "id = ?", clientID).Error
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
	var count int64
	if err := s.db.Model(&models.OAuth2Client{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListClients 获取客户端列表
func (s *OAuthService) ListClients(page, pageSize int) ([]models.OAuth2Client, int64, error) {
	var clients []models.OAuth2Client
	var total int64

	if err := s.db.Model(&models.OAuth2Client{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

// ListClientsFull 获取客户端列表（支持更多查询参数）
func (s *OAuthService) ListClientsFull(params *ListClientsParams) ([]models.OAuth2Client, int64, error) {
	var clients []models.OAuth2Client
	var total int64

	query := s.db.Model(&models.OAuth2Client{})

	// 搜索条件
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR id::text ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// 状态过滤
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderClause := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.SortDesc {
			order = "DESC"
		}
		orderClause = params.SortBy + " " + order
	}
	query = query.Order(orderClause)

	// 分页
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

// RegenerateClientSecret 重新生成客户端密钥
func (s *OAuthService) RegenerateClientSecret(clientID, newSecret string) error {
	return s.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Update("secret", newSecret).Error
}

// IncrementClientRequestCount 增加客户端请求计数
func (s *OAuthService) IncrementClientRequestCount(clientID string) error {
	return s.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Updates(map[string]interface{}{
			"request_count": gorm.Expr("request_count + 1"),
			"last_used_at":  time.Now(),
		}).Error
}

// GetClientStats 获取客户端统计数据
func (s *OAuthService) GetClientStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总客户端数
	var total int64
	if err := s.db.Model(&models.OAuth2Client{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 各状态客户端数
	var statusCounts []struct {
		Status models.ClientStatus
		Count  int64
	}
	if err := s.db.Model(&models.OAuth2Client{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		stats[string(sc.Status)] = sc.Count
	}

	return stats, nil
}

// RevokeToken 撤销令牌
func (s *OAuthService) RevokeToken(accessToken string) error {
	return s.manager.RemoveAccessToken(context.Background(), accessToken)
}

// RevokeRefreshToken 撤销刷新令牌
func (s *OAuthService) RevokeRefreshToken(refreshToken string) error {
	return s.manager.RemoveRefreshToken(context.Background(), refreshToken)
}

// GetTokenInfo 获取令牌信息
func (s *OAuthService) GetTokenInfo(accessToken string) (oauth2.TokenInfo, error) {
	return s.manager.LoadAccessToken(context.Background(), accessToken)
}
