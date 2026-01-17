package services

import (
	"encoding/json"
	"qauth-server/internal/config"
	"testing"
	"time"
)

func getTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
			Mode: "test",
		},
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			AccessTokenExp:  time.Hour,
			RefreshTokenExp: 24 * time.Hour,
		},
		OIDC: config.OIDCConfig{
			Issuer:              "https://auth.example.com",
			KeyRotationInterval: 24 * time.Hour,
			KeySize:             2048,
		},
	}
}

func TestNewOIDCService(t *testing.T) {
	cfg := getTestConfig()

	service, err := NewOIDCService(cfg, "https://test.example.com")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	if service == nil {
		t.Fatal("NewOIDCService() returned nil")
	}

	if service.GetJWKSManager() == nil {
		t.Error("JWKS manager is nil")
	}

	oidcCfg := service.GetOpenIDConfiguration()
	if oidcCfg.Issuer != "https://auth.example.com" {
		t.Errorf("Issuer = %v, want https://auth.example.com", oidcCfg.Issuer)
	}
}

func TestNewOIDCService_WithEmptyIssuer(t *testing.T) {
	cfg := getTestConfig()
	cfg.OIDC.Issuer = ""

	service, err := NewOIDCService(cfg, "https://fallback.example.com")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	oidcCfg := service.GetOpenIDConfiguration()
	if oidcCfg.Issuer != "https://fallback.example.com" {
		t.Errorf("Issuer = %v, want https://fallback.example.com", oidcCfg.Issuer)
	}
}

func TestOIDCService_GetOpenIDConfiguration(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	oidcCfg := service.GetOpenIDConfiguration()

	if oidcCfg.AuthorizationEndpoint == "" {
		t.Error("AuthorizationEndpoint is empty")
	}
	if oidcCfg.TokenEndpoint == "" {
		t.Error("TokenEndpoint is empty")
	}
	if oidcCfg.JwksURI == "" {
		t.Error("JwksURI is empty")
	}
	if oidcCfg.UserinfoEndpoint == "" {
		t.Error("UserinfoEndpoint is empty")
	}

	if len(oidcCfg.ScopesSupported) == 0 {
		t.Error("ScopesSupported is empty")
	}
	if len(oidcCfg.ResponseTypesSupported) == 0 {
		t.Error("ResponseTypesSupported is empty")
	}
	if len(oidcCfg.GrantTypesSupported) == 0 {
		t.Error("GrantTypesSupported is empty")
	}
	if len(oidcCfg.IDTokenSigningAlgValuesSupported) == 0 {
		t.Error("IDTokenSigningAlgValuesSupported is empty")
	}

	hasOpenID := false
	for _, scope := range oidcCfg.ScopesSupported {
		if scope == "openid" {
			hasOpenID = true
			break
		}
	}
	if !hasOpenID {
		t.Error("ScopesSupported does not contain 'openid'")
	}

	hasRS256 := false
	for _, alg := range oidcCfg.IDTokenSigningAlgValuesSupported {
		if alg == "RS256" {
			hasRS256 = true
			break
		}
	}
	if !hasRS256 {
		t.Error("IDTokenSigningAlgValuesSupported does not contain 'RS256'")
	}
}

func TestOIDCService_GetOpenIDConfiguration_JSON(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	oidcCfg := service.GetOpenIDConfiguration()

	jsonData, err := json.Marshal(oidcCfg)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	requiredFields := []string{
		"issuer",
		"authorization_endpoint",
		"token_endpoint",
		"jwks_uri",
		"scopes_supported",
		"response_types_supported",
		"grant_types_supported",
		"subject_types_supported",
		"id_token_signing_alg_values_supported",
	}

	for _, field := range requiredFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("JSON missing required field: %s", field)
		}
	}
}

func TestOIDCService_GetJWKS(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	jwks := service.GetJWKS()

	if len(jwks.Keys) == 0 {
		t.Error("JWKS has no keys")
	}

	for _, key := range jwks.Keys {
		if key.Kty != "RSA" {
			t.Errorf("Key type = %v, want RSA", key.Kty)
		}
		if key.Alg != "RS256" {
			t.Errorf("Key algorithm = %v, want RS256", key.Alg)
		}
		if key.Kid == "" {
			t.Error("Key ID is empty")
		}
	}
}

func TestOIDCService_GetJWKSJSON(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	jsonData, err := service.GetJWKSJSON()
	if err != nil {
		t.Fatalf("GetJWKSJSON() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	if _, ok := parsed["keys"]; !ok {
		t.Error("JWKS JSON missing 'keys' field")
	}
}

func TestOIDCService_ForceKeyRotation(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	oldKeyID := service.GetActiveKeyID()

	if err := service.ForceKeyRotation(); err != nil {
		t.Fatalf("ForceKeyRotation() error = %v", err)
	}

	newKeyID := service.GetActiveKeyID()

	if oldKeyID == newKeyID {
		t.Error("Key was not rotated")
	}

	jwks := service.GetJWKS()
	if len(jwks.Keys) != 2 {
		t.Errorf("JWKS has %d keys after rotation, want 2", len(jwks.Keys))
	}
}

func TestOIDCService_GetKeyRotationInfo(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	service.ForceKeyRotation()
	service.ForceKeyRotation()

	keys := service.GetKeyRotationInfo()

	if len(keys) != 3 {
		t.Errorf("GetKeyRotationInfo() returned %d keys, want 3", len(keys))
	}

	for _, key := range keys {
		if key.PrivateKey != nil {
			t.Error("Key info should not contain private key")
		}
		if key.ID == "" {
			t.Error("Key ID is empty")
		}
	}
}

func TestOIDCService_GetActiveKeyID(t *testing.T) {
	cfg := getTestConfig()
	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	keyID := service.GetActiveKeyID()

	if keyID == "" {
		t.Error("Active key ID is empty")
	}
}

func TestOIDCService_EndpointURLs(t *testing.T) {
	cfg := getTestConfig()
	cfg.OIDC.Issuer = "https://auth.example.com"

	service, err := NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("NewOIDCService() error = %v", err)
	}
	defer service.Close()

	oidcCfg := service.GetOpenIDConfiguration()

	expectedEndpoints := map[string]string{
		"authorization_endpoint": "https://auth.example.com/oauth/authorize",
		"token_endpoint":         "https://auth.example.com/oauth/token",
		"userinfo_endpoint":      "https://auth.example.com/oauth/userinfo",
		"jwks_uri":               "https://auth.example.com/.well-known/jwks.json",
		"revocation_endpoint":    "https://auth.example.com/oauth/revoke",
		"end_session_endpoint":   "https://auth.example.com/oauth/logout",
	}

	if oidcCfg.AuthorizationEndpoint != expectedEndpoints["authorization_endpoint"] {
		t.Errorf("AuthorizationEndpoint = %v, want %v", oidcCfg.AuthorizationEndpoint, expectedEndpoints["authorization_endpoint"])
	}
	if oidcCfg.TokenEndpoint != expectedEndpoints["token_endpoint"] {
		t.Errorf("TokenEndpoint = %v, want %v", oidcCfg.TokenEndpoint, expectedEndpoints["token_endpoint"])
	}
	if oidcCfg.UserinfoEndpoint != expectedEndpoints["userinfo_endpoint"] {
		t.Errorf("UserinfoEndpoint = %v, want %v", oidcCfg.UserinfoEndpoint, expectedEndpoints["userinfo_endpoint"])
	}
	if oidcCfg.JwksURI != expectedEndpoints["jwks_uri"] {
		t.Errorf("JwksURI = %v, want %v", oidcCfg.JwksURI, expectedEndpoints["jwks_uri"])
	}
	if oidcCfg.RevocationEndpoint != expectedEndpoints["revocation_endpoint"] {
		t.Errorf("RevocationEndpoint = %v, want %v", oidcCfg.RevocationEndpoint, expectedEndpoints["revocation_endpoint"])
	}
	if oidcCfg.EndSessionEndpoint != expectedEndpoints["end_session_endpoint"] {
		t.Errorf("EndSessionEndpoint = %v, want %v", oidcCfg.EndSessionEndpoint, expectedEndpoints["end_session_endpoint"])
	}
}
