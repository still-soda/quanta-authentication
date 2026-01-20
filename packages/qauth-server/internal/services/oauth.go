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
	"strings"

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
	Cfg         *config.Config
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
	})

	// 创建 OAuth2 管理器
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

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

	// 设置内部错误处理
	srv.SetInternalErrorHandler(func(err error) (re *oautherrors.Response) {
		logger.Error("OAuth2 internal error", "error", err)
		return
	})

	// 设置响应错误处理
	srv.SetResponseErrorHandler(func(re *oautherrors.Response) {
		logger.Error("OAuth2 response error", "error", re.Error, "description", re.Description)
	})

	service := &OAuthService{
		db:          db,
		Cfg:         cfg,
		server:      srv,
		manager:     manager,
		clientStore: clientStore,
		jwksManager: jwksManager,
		userService: userService,
	}

	// 设置用户授权处理器
	srv.SetUserAuthorizationHandler(service.userAuthorizationHandler)

	// 设置密码授权处理器
	srv.SetPasswordAuthorizationHandler(service.passwordAuthorizationHandler)

	// 设置扩展字段处理器
	srv.SetExtensionFieldsHandler(service.extensionFieldsHandler)

	logger.Info("OAuth2 service initialized")
	return service
}

// userAuthorizationHandler 用户授权处理器（授权码模式）
func (s *OAuthService) userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	// 从 Cookie 中获取 JWT Token
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", errors.New("user not logged in")
	}

	// 解析 JWT Token
	claims, err := app_jwt.ParseAccessToken(cookie.Value)
	if err != nil {
		return "", errors.New("invalid access token")
	}

	// 验证用户是否存在
	var user models.Users
	if err := s.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		return "", errors.New("user not found")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return "", errors.New("user account is not active")
	}

	return user.ID, nil
}

// passwordAuthorizationHandler 密码授权处理器（密码模式）
func (s *OAuthService) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	// 根据学号查找用户
	var user models.Users
	if err := s.db.First(&user, "student_id = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// 验证密码
	if !utilities.VerifyPassword(password, user.Salt, user.PasswordHash) {
		return "", errors.New("invalid password")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return "", errors.New("user account is not active")
	}

	// 记录登录状态
	s.recordLoginState(user.ID, &clientID, models.LoginTypeOAuth2, true, "")

	return user.ID, nil
}

// extensionFieldsHandler 扩展字段处理器（用于生成 ID Token）
func (s *OAuthService) extensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]any) {
	if !strings.Contains(ti.GetScope(), "openid") {
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

	idTokenClaims := &jwks.IDTokenClaims{
		UserID:      user.ID,
		StudentID:   user.StudentID,
		DisplayName: *user.DisplayName,
		Avatar:      "", // TODO: 添加头像支持

		Issuer:   s.Cfg.OIDC.Issuer,
		Subject:  user.ID,
		Audience: ti.GetClientID(),
		Expiry:   ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Unix(),
		IssuedAt: ti.GetAccessCreateAt().Unix(),
		AuthTime: ti.GetAccessCreateAt().Unix(),
	}
	idToken, err := s.jwksManager.SignToken(idTokenClaims)
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

// GetClientByID 根据 ID 获取客户端
func (s *OAuthService) GetClientByID(clientID string) (*models.OAuth2Client, error) {
	var client models.OAuth2Client
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
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

// DeleteClient 删除客户端
func (s *OAuthService) DeleteClient(clientID string) error {
	return s.db.Delete(&models.OAuth2Client{}, "id = ?", clientID).Error
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
