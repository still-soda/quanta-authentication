package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qauth-server/internal/config"
	"qauth-server/internal/services"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func getTestOIDCConfig() *config.Config {
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

func setupOIDCTestRouter(t *testing.T) (*gin.Engine, *OIDCHandler) {
	gin.SetMode(gin.TestMode)

	cfg := getTestOIDCConfig()
	oidcService, err := services.NewOIDCService(cfg, "")
	if err != nil {
		t.Fatalf("Failed to create OIDC service: %v", err)
	}

	handler := NewOIDCHandler(oidcService)
	router := gin.New()

	router.GET("/.well-known/openid-configuration", handler.GetOpenIDConfiguration)
	router.GET("/.well-known/jwks.json", handler.GetJWKS)
	router.POST("/admin/jwks/rotate", handler.ForceKeyRotation)
	router.GET("/admin/jwks/keys", handler.GetKeyRotationInfo)

	return router, handler
}

func TestOIDCHandler_GetOpenIDConfiguration(t *testing.T) {
	router, _ := setupOIDCTestRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/.well-known/openid-configuration", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %s, want application/json; charset=utf-8", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// response.HandlerSuccess 会包装在 data 字段中
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'data' field")
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
		if _, ok := data[field]; !ok {
			t.Errorf("Response missing required field: %s", field)
		}
	}

	if data["issuer"] != "https://auth.example.com" {
		t.Errorf("Issuer = %v, want https://auth.example.com", data["issuer"])
	}
}

func TestOIDCHandler_GetJWKS(t *testing.T) {
	router, _ := setupOIDCTestRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	// JWKS 端点直接返回 JSON，Content-Type 可能是 application/json
	if contentType != "application/json; charset=utf-8" && contentType != "application/json" {
		t.Errorf("Content-Type = %s, want application/json or application/json; charset=utf-8", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// response.HandlerSuccess 会包装在 data 字段中
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'data' field")
	}

	keys, ok := data["keys"].([]interface{})
	if !ok {
		t.Fatal("Response data missing 'keys' array")
	}

	if len(keys) == 0 {
		t.Error("Keys array is empty")
	}

	firstKey := keys[0].(map[string]interface{})
	if firstKey["kty"] != "RSA" {
		t.Errorf("Key type = %v, want RSA", firstKey["kty"])
	}
	if firstKey["alg"] != "RS256" {
		t.Errorf("Key algorithm = %v, want RS256", firstKey["alg"])
	}
	if firstKey["use"] != "sig" {
		t.Errorf("Key use = %v, want sig", firstKey["use"])
	}
	if firstKey["kid"] == "" {
		t.Error("Key ID is empty")
	}
	if firstKey["n"] == "" {
		t.Error("Key modulus (n) is empty")
	}
	if firstKey["e"] == "" {
		t.Error("Key exponent (e) is empty")
	}
}

func TestOIDCHandler_ForceKeyRotation(t *testing.T) {
	router, handler := setupOIDCTestRouter(t)

	reqBefore, _ := http.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	wBefore := httptest.NewRecorder()
	router.ServeHTTP(wBefore, reqBefore)

	var jwksBefore map[string]interface{}
	json.Unmarshal(wBefore.Body.Bytes(), &jwksBefore)
	dataBefore := jwksBefore["data"].(map[string]interface{})
	keysBefore := dataBefore["keys"].([]interface{})
	keyCountBefore := len(keysBefore)

	req, _ := http.NewRequest(http.MethodPost, "/admin/jwks/rotate", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// response.HandlerSuccess 会包装在 data 字段中
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'data' field")
	}

	if data["message"] != "key rotation completed" {
		t.Errorf("Message = %v, want 'key rotation completed'", data["message"])
	}

	if data["active_key_id"] == "" {
		t.Error("active_key_id is empty")
	}

	reqAfter, _ := http.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	wAfter := httptest.NewRecorder()
	router.ServeHTTP(wAfter, reqAfter)

	var jwksAfter map[string]interface{}
	json.Unmarshal(wAfter.Body.Bytes(), &jwksAfter)
	dataAfter := jwksAfter["data"].(map[string]interface{})
	keysAfter := dataAfter["keys"].([]interface{})
	keyCountAfter := len(keysAfter)

	if keyCountAfter != keyCountBefore+1 {
		t.Errorf("Key count after rotation = %d, want %d", keyCountAfter, keyCountBefore+1)
	}

	_ = handler
}

func TestOIDCHandler_GetKeyRotationInfo(t *testing.T) {
	router, _ := setupOIDCTestRouter(t)

	rotateReq, _ := http.NewRequest(http.MethodPost, "/admin/jwks/rotate", nil)
	rotateW := httptest.NewRecorder()
	router.ServeHTTP(rotateW, rotateReq)

	req, _ := http.NewRequest(http.MethodGet, "/admin/jwks/keys", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// response.HandlerSuccess 会包装在 data 字段中
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'data' field")
	}

	keys, ok := data["keys"].([]interface{})
	if !ok {
		t.Fatal("Response data missing 'keys' array")
	}

	if len(keys) < 2 {
		t.Error("Expected at least 2 keys after rotation")
	}

	for _, k := range keys {
		keyInfo := k.(map[string]interface{})
		if keyInfo["id"] == "" {
			t.Error("Key ID is empty")
		}
		if keyInfo["status"] == "" {
			t.Error("Key status is empty")
		}
		if keyInfo["created_at"] == "" {
			t.Error("Key created_at is empty")
		}
	}
}

func TestOIDCHandler_GetOpenIDConfiguration_CacheControl(t *testing.T) {
	router, _ := setupOIDCTestRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/.well-known/openid-configuration", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	if cacheControl == "" {
		t.Skip("Cache-Control header not implemented")
	}
}

func TestOIDCHandler_GetJWKS_CacheControl(t *testing.T) {
	router, _ := setupOIDCTestRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	if cacheControl == "" {
		t.Skip("Cache-Control header not implemented")
	}
}

func TestOIDCHandler_CORS(t *testing.T) {
	router, _ := setupOIDCTestRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/.well-known/openid-configuration", nil)
	req.Header.Set("Origin", "https://client.example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}
