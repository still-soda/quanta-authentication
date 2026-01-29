package providers

import (
	"context"
	"net/http"
	"qauth-server/internal/config"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	oautherrors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauthmodels "github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ITokenInfo 令牌信息接口
type ITokenInfo interface {
	GetClientID() string
	GetUserID() string
	GetScope() string
	GetAccessCreateAt() time.Time
	GetAccessExpiresIn() time.Duration
	GetRefreshCreateAt() time.Time
	GetRefreshExpiresIn() time.Duration
	GetRefresh() string
}

// AuthorizeRequest 授权请求
type AuthorizeRequest struct {
	ResponseType string
	ClientID     string
	Scope        string
	RedirectURI  string
	State        string
}

// IClientInfo OAuth 客户端信息接口
type IClientInfo interface {
	GetID() string
	GetSecret() string
	GetDomain() string
	GetUserID() string
	IsPublic() bool
}

// UserAuthorizationHandler 用户授权处理器
type UserAuthorizationHandler func(w http.ResponseWriter, r *http.Request) (userID string, err error)

// PasswordAuthorizationHandler 密码授权处理器
type PasswordAuthorizationHandler func(ctx context.Context, clientID, username, password string) (userID string, err error)

// ExtensionFieldsHandler 扩展字段处理器（用于生成 ID Token 等）
type ExtensionFieldsHandler func(ti ITokenInfo) (fieldsValue map[string]any)

// IOAuth OAuth 提供者接口
type IOAuth interface {
	// SetResponseErrorHandler 设置响应错误处理器
	SetResponseErrorHandler(handler func(re *oautherrors.Response))

	// SetInternalErrorHandler 设置内部错误处理器
	SetInternalErrorHandler(handler func(err error) (re *oautherrors.Response))

	// HandleAuthorizeRequest 处理授权请求
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error

	// ValidateAuthorizeRequest 验证授权请求
	ValidateAuthorizeRequest(r *http.Request) (*AuthorizeRequest, error)

	// HandleTokenRequest 处理令牌请求
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) error

	// ValidateToken 验证访问令牌
	ValidateToken(r *http.Request) (ITokenInfo, error)

	// RevokeToken 撤销令牌
	RevokeToken(ctx context.Context, accessToken string) error

	// RevokeRefreshToken 撤销刷新令牌
	RevokeRefreshToken(ctx context.Context, refreshToken string) error

	// GetTokenInfo 获取令牌信息
	GetTokenInfo(ctx context.Context, accessToken string) (ITokenInfo, error)

	// SetUserAuthorizationHandler 设置用户授权处理器
	SetUserAuthorizationHandler(handler UserAuthorizationHandler)

	// SetPasswordAuthorizationHandler 设置密码授权处理器
	SetPasswordAuthorizationHandler(handler PasswordAuthorizationHandler)

	// SetExtensionFieldsHandler 设置扩展字段处理器
	SetExtensionFieldsHandler(handler ExtensionFieldsHandler)
}

// IClientStore 客户端存储接口
type IClientStore interface {
	// GetByID 根据 ID 获取客户端
	GetByID(ctx context.Context, id string) (IClientInfo, error)
}

// ITokenGenerator 令牌生成器接口
type ITokenGenerator interface {
	// GenerateAccessToken 生成访问令牌
	GenerateAccessToken(ctx context.Context, isGenerateRefresh bool) (access, refresh string, err error)
}

// ========= provider 实现 ========

// GoOAuth 基于 go-oauth2 库的 OAuth 提供者实现
type GoOAuth struct {
	server      *server.Server
	manager     *manage.Manager
	clientStore *gormClientStore
	cfg         *config.Config
}

// gormClientStore 基于 GORM 的 OAuth2 客户端存储
type gormClientStore struct {
	db *gorm.DB
}

// newGormClientStore 创建新的 GORM 客户端存储
func newGormClientStore(db *gorm.DB) *gormClientStore {
	return &gormClientStore{db: db}
}

// GetByID 根据 ID 获取 OAuth2 客户端
func (s *gormClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var client models.OAuth2Client
	if err := s.db.WithContext(ctx).First(&client, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.ErrClientNotFound.WithMessage(id)
		}
		return nil, e.ErrClientQueryFailed.Wrap(err)
	}

	return &oauthmodels.Client{
		ID:     client.ID,
		Secret: client.Secret,
		Domain: client.Domain,
		UserID: client.Data.UserID,
		Public: client.Data.Public,
	}, nil
}

// customAccessTokenGenerate 自定义令牌生成器
type customAccessTokenGenerate struct{}

func (g *customAccessTokenGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	access, err = utilities.GenerateRandomString(32)
	if isGenRefresh {
		refresh, err = utilities.GenerateRandomString(32)
	}

	if err != nil {
		return "", "", e.ErrTokenGenerationFailed.Wrap(err)
	}

	access = "atk:" + access
	refresh = "rtk:" + refresh

	return access, refresh, err
}

// tokenInfoWrapper 包装 go-oauth2 的 TokenInfo
type tokenInfoWrapper struct {
	oauth2.TokenInfo
}

func (w *tokenInfoWrapper) GetClientID() string {
	return w.TokenInfo.GetClientID()
}

func (w *tokenInfoWrapper) GetUserID() string {
	return w.TokenInfo.GetUserID()
}

func (w *tokenInfoWrapper) GetScope() string {
	return w.TokenInfo.GetScope()
}

func (w *tokenInfoWrapper) GetAccessCreateAt() time.Time {
	return w.TokenInfo.GetAccessCreateAt()
}

func (w *tokenInfoWrapper) GetAccessExpiresIn() time.Duration {
	return w.TokenInfo.GetAccessExpiresIn()
}

func (w *tokenInfoWrapper) GetRefreshCreateAt() time.Time {
	return w.TokenInfo.GetRefreshCreateAt()
}

func (w *tokenInfoWrapper) GetRefreshExpiresIn() time.Duration {
	return w.TokenInfo.GetRefreshExpiresIn()
}

func (w *tokenInfoWrapper) GetRefresh() string {
	return w.TokenInfo.GetRefresh()
}

// NewGoOAuth 创建新的 go-oauth2 提供者
func NewGoOAuth(db *gorm.DB, cfg *config.Config) *GoOAuth {
	// 创建基于 GORM 的客户端存储
	clientStore := newGormClientStore(db)

	// 创建 Redis Token 存储
	tokenStore := oredis.NewRedisStore(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, "oauth2_token:")

	// 创建 OAuth2 管理器
	mgr := manage.NewDefaultManager()
	mgr.MapTokenStorage(tokenStore)
	mgr.MapClientStorage(clientStore)

	// 使用自定义的令牌生成器
	mgr.MapAccessGenerate(&customAccessTokenGenerate{})

	// 配置授权码模式
	mgr.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    cfg.JWT.AccessTokenExp,
		RefreshTokenExp:   cfg.JWT.RefreshTokenExp,
		IsGenerateRefresh: true,
	})

	// 配置密码模式
	mgr.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp:    cfg.JWT.AccessTokenExp,
		RefreshTokenExp:   cfg.JWT.RefreshTokenExp,
		IsGenerateRefresh: true,
	})

	// 配置客户端凭证模式
	mgr.SetClientTokenCfg(&manage.Config{
		AccessTokenExp: cfg.JWT.AccessTokenExp,
	})

	// 配置刷新令牌
	mgr.SetRefreshTokenCfg(&manage.RefreshingConfig{
		IsGenerateRefresh:  true,
		IsRemoveAccess:     true,
		IsRemoveRefreshing: true,
	})

	// 创建 OAuth2 服务器
	srv := server.NewDefaultServer(mgr)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	// 仅允许配置的授权类型
	gt := []oauth2.GrantType{}
	for _, grantType := range cfg.OAuth.GrantTypesSupported {
		gt = append(gt, oauth2.GrantType(grantType))
	}
	srv.SetAllowedGrantType(gt...)

	return &GoOAuth{
		server:      srv,
		manager:     mgr,
		clientStore: clientStore,
		cfg:         cfg,
	}
}

// SetInternalErrorHandler 设置内部错误处理器
func (p *GoOAuth) SetInternalErrorHandler(handler func(err error) (re *oautherrors.Response)) {
	p.server.SetInternalErrorHandler(handler)
}

// SetResponseErrorHandler 设置响应错误处理器
func (p *GoOAuth) SetResponseErrorHandler(handler func(re *oautherrors.Response)) {
	p.server.SetResponseErrorHandler(handler)
}

// HandleAuthorizeRequest 处理授权请求
func (p *GoOAuth) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return p.server.HandleAuthorizeRequest(w, r)
}

// ValidateAuthorizeRequest 验证授权请求
func (p *GoOAuth) ValidateAuthorizeRequest(r *http.Request) (*AuthorizeRequest, error) {
	req, err := p.server.ValidationAuthorizeRequest(r)
	if err != nil {
		return nil, err
	}

	return &AuthorizeRequest{
		ResponseType: string(req.ResponseType),
		ClientID:     req.ClientID,
		Scope:        req.Scope,
		RedirectURI:  req.RedirectURI,
		State:        req.State,
	}, nil
}

// HandleTokenRequest 处理令牌请求
func (p *GoOAuth) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return p.server.HandleTokenRequest(w, r)
}

// ValidateToken 验证访问令牌
func (p *GoOAuth) ValidateToken(r *http.Request) (ITokenInfo, error) {
	ti, err := p.server.ValidationBearerToken(r)
	if err != nil {
		return nil, err
	}
	return &tokenInfoWrapper{TokenInfo: ti}, nil
}

// RevokeToken 撤销令牌
func (p *GoOAuth) RevokeToken(ctx context.Context, accessToken string) error {
	if err := p.manager.RemoveAccessToken(ctx, accessToken); err != nil {
		return e.ErrTokenRevocationFailed.Wrap(err)
	}
	return nil
}

// RevokeRefreshToken 撤销刷新令牌
func (p *GoOAuth) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if err := p.manager.RemoveRefreshToken(ctx, refreshToken); err != nil {
		return e.ErrTokenRevocationFailed.Wrap(err)
	}
	return nil
}

// GetTokenInfo 获取令牌信息
func (p *GoOAuth) GetTokenInfo(ctx context.Context, accessToken string) (ITokenInfo, error) {
	info, err := p.manager.LoadAccessToken(ctx, accessToken)
	if err != nil {
		return nil, e.ErrTokenLoadFailed.Wrap(err)
	}
	return &tokenInfoWrapper{TokenInfo: info}, nil
}

// SetUserAuthorizationHandler 设置用户授权处理器
func (p *GoOAuth) SetUserAuthorizationHandler(handler UserAuthorizationHandler) {
	// 创建适配器，将我们的接口转换为 go-oauth2 的接口
	p.server.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return handler(w, r)
	})
}

// SetPasswordAuthorizationHandler 设置密码授权处理器
func (p *GoOAuth) SetPasswordAuthorizationHandler(handler PasswordAuthorizationHandler) {
	// 创建适配器，将我们的接口转换为 go-oauth2 的接口
	p.server.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		return handler(ctx, clientID, username, password)
	})
}

// SetExtensionFieldsHandler 设置扩展字段处理器
func (p *GoOAuth) SetExtensionFieldsHandler(handler ExtensionFieldsHandler) {
	// 创建适配器，将我们的接口转换为 go-oauth2 的接口
	p.server.SetExtensionFieldsHandler(func(ti oauth2.TokenInfo) (fieldsValue map[string]any) {
		wrapper := &tokenInfoWrapper{TokenInfo: ti}
		return handler(wrapper)
	})
}
