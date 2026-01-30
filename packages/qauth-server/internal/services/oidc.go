package services

import (
	"qauth-server/internal/config"
	"qauth-server/internal/providers"
	"qauth-server/pkg/jwks"
)

// OIDCDiscoveryConfig OIDC 发现配置
type OIDCDiscoveryConfig struct {
	Issuer                            string                           `json:"issuer"`
	AuthorizationEndpoint             string                           `json:"authorization_endpoint"`
	TokenEndpoint                     string                           `json:"token_endpoint"`
	UserinfoEndpoint                  string                           `json:"userinfo_endpoint"`
	JwksURI                           string                           `json:"jwks_uri"`
	RegistrationEndpoint              string                           `json:"registration_endpoint,omitempty"`
	RevocationEndpoint                string                           `json:"revocation_endpoint,omitempty"`
	IntrospectionEndpoint             string                           `json:"introspection_endpoint,omitempty"`
	EndSessionEndpoint                string                           `json:"end_session_endpoint,omitempty"`
	ScopesSupported                   []config.Scope                   `json:"scopes_supported"`
	ResponseTypesSupported            []config.ResponseType            `json:"response_types_supported"`
	ResponseModesSupported            []config.ResponseMode            `json:"response_modes_supported,omitempty"`
	GrantTypesSupported               []config.GrantType               `json:"grant_types_supported"`
	SubjectTypesSupported             []config.SubjectType             `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []config.IDTokenSigningAlg       `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []config.TokenEndpointAuthMethod `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   []config.Claim                   `json:"claims_supported,omitempty"`
	CodeChallengeMethodsSupported     []config.CodeChallengeMethod     `json:"code_challenge_methods_supported,omitempty"`
}

// OIDCService OIDC 服务
type OIDCService struct {
	cfg         *config.Config
	jwksManager *jwks.JWKSManager
	issuer      string
	logger      providers.ILogger
}

// NewOIDCService 创建新的 OIDC 服务
func NewOIDCService(
	cfg *config.Config,
	issuer string,
	jwksManager *jwks.JWKSManager,
	logger providers.ILogger,
) (*OIDCService, error) {
	// 如果配置中有 issuer，优先使用配置
	if cfg.OIDC.Issuer != "" {
		issuer = cfg.OIDC.Issuer
	}

	service := &OIDCService{
		cfg:         cfg,
		jwksManager: jwksManager,
		issuer:      issuer,
		logger:      logger.With("service", "OIDCService"),
	}

	return service, nil
}

func (s *OIDCService) GetConfig() *config.OIDCConfig {
	return &s.cfg.OIDC
}

// GetJWKSManager 获取 JWKS 管理器
func (s *OIDCService) GetJWKSManager() *jwks.JWKSManager {
	return s.jwksManager
}

// GetOpenIDConfiguration 获取 OpenID Connect 配置
func (s *OIDCService) GetOpenIDConfiguration() *OIDCDiscoveryConfig {
	return &OIDCDiscoveryConfig{
		Issuer:                            s.issuer,
		AuthorizationEndpoint:             s.issuer + "/oauth/authorize",
		TokenEndpoint:                     s.issuer + "/oauth/token",
		UserinfoEndpoint:                  s.issuer + "/oauth/userinfo",
		JwksURI:                           s.issuer + "/.well-known/jwks.json",
		RevocationEndpoint:                s.issuer + "/oauth/revoke",
		IntrospectionEndpoint:             s.issuer + "/oauth/validate",
		EndSessionEndpoint:                s.issuer + "/oauth/logout",
		ScopesSupported:                   s.cfg.OAuth.ScopeSupported,
		ResponseTypesSupported:            s.cfg.OAuth.AllowedResponseTypes,
		ResponseModesSupported:            s.cfg.OAuth.ResponseModesSupported,
		GrantTypesSupported:               s.cfg.OAuth.GrantTypesSupported,
		SubjectTypesSupported:             s.cfg.OAuth.SubjectTypesSupported,
		IDTokenSigningAlgValuesSupported:  s.cfg.OAuth.IDTokenSigningAlgValuesSupported,
		TokenEndpointAuthMethodsSupported: s.cfg.OAuth.TokenEndpointAuthMethodsSupported,
		ClaimsSupported:                   s.cfg.OAuth.ClaimsSupported,
		CodeChallengeMethodsSupported:     s.cfg.OAuth.CodeChallengeMethodsSupported,
	}
}

// GetJWKS 获取 JWKS
func (s *OIDCService) GetJWKS() *jwks.JWKS {
	return s.jwksManager.GetJWKS()
}

// GetJWKSJSON 获取 JWKS JSON
func (s *OIDCService) GetJWKSJSON() ([]byte, error) {
	return s.jwksManager.GetJWKSJSON()
}

// Close 关闭服务
func (s *OIDCService) Close() {
	// jwksManager 由外部管理，不在这里停止
}

// ForceKeyRotation 强制密钥轮换
func (s *OIDCService) ForceKeyRotation() error {
	return s.jwksManager.ForceRotate()
}

// GetKeyRotationInfo 获取密钥轮换信息
func (s *OIDCService) GetKeyRotationInfo() []jwks.KeyInfo {
	return s.jwksManager.GetAllKeys()
}

// GetActiveKeyID 获取当前活跃密钥 ID
func (s *OIDCService) GetActiveKeyID() string {
	return s.jwksManager.GetActiveKeyID()
}

// ValidateIDToken 验证 ID 令牌
func (s *OIDCService) ValidateIDToken(idToken string) (*jwks.BasicClaims, error) {
	return s.jwksManager.VerifyToken(idToken)
}
