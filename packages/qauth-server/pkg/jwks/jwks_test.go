package jwks

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewJWKSManager(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *JWKSConfig
		wantErr bool
	}{
		{
			name:    "default config",
			cfg:     nil,
			wantErr: false,
		},
		{
			name: "custom config",
			cfg: &JWKSConfig{
				KeySize:          2048,
				RotationInterval: 12 * time.Hour,
				GracePeriod:      30 * time.Minute,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewJWKSManager(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJWKSManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if manager == nil {
				t.Error("NewJWKSManager() returned nil manager")
				return
			}

			keyInfo, err := manager.GetActiveKey()
			if err != nil {
				t.Errorf("GetActiveKey() error = %v", err)
				return
			}
			if keyInfo.PrivateKey == nil {
				t.Error("Active key has nil PrivateKey")
			}
			if keyInfo.PublicKey == nil {
				t.Error("Active key has nil PublicKey")
			}
			if keyInfo.Status != KeyStatusActive {
				t.Errorf("Active key status = %v, want %v", keyInfo.Status, KeyStatusActive)
			}
		})
	}
}

func TestJWKSManager_GenerateNewKey(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	oldKeyID := manager.GetActiveKeyID()

	if err := manager.ForceRotate(); err != nil {
		t.Fatalf("ForceRotate() error = %v", err)
	}

	newKeyID := manager.GetActiveKeyID()

	if oldKeyID == newKeyID {
		t.Error("ForceRotate() did not generate a new key")
	}

	oldKey, err := manager.GetKeyByID(oldKeyID)
	if err != nil {
		t.Fatalf("GetKeyByID() for old key error = %v", err)
	}
	if oldKey.Status != KeyStatusRotating {
		t.Errorf("Old key status = %v, want %v", oldKey.Status, KeyStatusRotating)
	}

	newKey, err := manager.GetKeyByID(newKeyID)
	if err != nil {
		t.Fatalf("GetKeyByID() for new key error = %v", err)
	}
	if newKey.Status != KeyStatusActive {
		t.Errorf("New key status = %v, want %v", newKey.Status, KeyStatusActive)
	}
}

func TestJWKSManager_GetJWKS(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	if err := manager.ForceRotate(); err != nil {
		t.Fatalf("ForceRotate() error = %v", err)
	}

	jwks := manager.GetJWKS()

	if len(jwks.Keys) != 2 {
		t.Errorf("JWKS has %d keys, want 2", len(jwks.Keys))
	}

	for _, jwk := range jwks.Keys {
		if jwk.Kty != "RSA" {
			t.Errorf("JWK Kty = %v, want RSA", jwk.Kty)
		}
		if jwk.Use != "sig" {
			t.Errorf("JWK Use = %v, want sig", jwk.Use)
		}
		if jwk.Alg != "RS256" {
			t.Errorf("JWK Alg = %v, want RS256", jwk.Alg)
		}
		if jwk.Kid == "" {
			t.Error("JWK Kid is empty")
		}
	}
}

func TestJWKSManager_GetJWKSJSON(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	jsonBytes, err := manager.GetJWKSJSON()
	if err != nil {
		t.Fatalf("GetJWKSJSON() error = %v", err)
	}

	jsonStr := string(jsonBytes)

	if !strings.Contains(jsonStr, `"keys"`) {
		t.Error("JWKS JSON missing 'keys' field")
	}
	if !strings.Contains(jsonStr, `"kty":"RSA"`) {
		t.Error("JWKS JSON missing 'kty' field")
	}
}

func TestJWKSManager_SignAndVerifyToken(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	claims := jwt.MapClaims{
		"sub":  "user123",
		"name": "Test User",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	tokenString, err := manager.SignToken(claims)
	if err != nil {
		t.Fatalf("SignToken() error = %v", err)
	}

	token, err := manager.VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}

	if !token.Valid {
		t.Error("Token is not valid")
	}

	parsedClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to parse claims")
	}

	if parsedClaims["sub"] != "user123" {
		t.Errorf("sub = %v, want user123", parsedClaims["sub"])
	}
}

func TestJWKSManager_VerifyTokenAfterRotation(t *testing.T) {
	manager, err := NewJWKSManager(&JWKSConfig{
		KeySize:          2048,
		RotationInterval: 24 * time.Hour,
		GracePeriod:      1 * time.Hour,
	})
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	oldToken, err := manager.SignToken(claims)
	if err != nil {
		t.Fatalf("SignToken() error = %v", err)
	}

	if err := manager.ForceRotate(); err != nil {
		t.Fatalf("ForceRotate() error = %v", err)
	}

	newToken, err := manager.SignToken(claims)
	if err != nil {
		t.Fatalf("SignToken() error = %v", err)
	}

	_, err = manager.VerifyToken(oldToken)
	if err != nil {
		t.Errorf("VerifyToken() for old token error = %v", err)
	}

	_, err = manager.VerifyToken(newToken)
	if err != nil {
		t.Errorf("VerifyToken() for new token error = %v", err)
	}
}

func TestJWKSManager_VerifyTokenWithWrongSignature(t *testing.T) {
	manager1, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	manager2, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token, err := manager1.SignToken(claims)
	if err != nil {
		t.Fatalf("SignToken() error = %v", err)
	}

	_, err = manager2.VerifyToken(token)
	if err == nil {
		t.Error("VerifyToken() should fail with different key")
	}
}

func TestJWKSManager_GetAllKeys(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	if err := manager.ForceRotate(); err != nil {
		t.Fatalf("ForceRotate() error = %v", err)
	}
	if err := manager.ForceRotate(); err != nil {
		t.Fatalf("ForceRotate() error = %v", err)
	}

	keys := manager.GetAllKeys()

	if len(keys) != 3 {
		t.Errorf("GetAllKeys() returned %d keys, want 3", len(keys))
	}

	for _, key := range keys {
		if key.PrivateKey != nil {
			t.Error("GetAllKeys() returned key with PrivateKey")
		}
	}
}

func TestPublicKeyToJWK(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	keyInfo, err := manager.GetActiveKey()
	if err != nil {
		t.Fatalf("GetActiveKey() error = %v", err)
	}

	jwk := publicKeyToJWK(keyInfo.ID, keyInfo.PublicKey)

	if jwk.Kty != "RSA" {
		t.Errorf("JWK Kty = %v, want RSA", jwk.Kty)
	}
	if jwk.Use != "sig" {
		t.Errorf("JWK Use = %v, want sig", jwk.Use)
	}
	if jwk.Alg != "RS256" {
		t.Errorf("JWK Alg = %v, want RS256", jwk.Alg)
	}
	if jwk.Kid != keyInfo.ID {
		t.Errorf("JWK Kid = %v, want %v", jwk.Kid, keyInfo.ID)
	}
}

func TestJWKSManager_CleanupExpiredKeys(t *testing.T) {
	cfg := &JWKSConfig{
		KeySize:          2048,
		RotationInterval: 100 * time.Millisecond,
		GracePeriod:      50 * time.Millisecond,
	}

	manager, err := NewJWKSManager(cfg)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	initialKeyID := manager.GetActiveKeyID()

	if err := manager.ForceRotate(); err != nil {
		t.Fatalf("ForceRotate() error = %v", err)
	}

	_, err = manager.GetKeyByID(initialKeyID)
	if err != nil {
		t.Errorf("Old key should still exist: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	manager.cleanupExpiredKeys()

	_, err = manager.GetKeyByID(initialKeyID)
	if err == nil {
		t.Error("Old key should have been cleaned up")
	}
}

func TestVerifyTokenWithInvalidMethod(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("Failed to create HS256 token: %v", err)
	}

	_, err = manager.VerifyToken(tokenString)
	if err == nil {
		t.Error("VerifyToken() should fail with wrong signing method")
	}
}

func TestVerifyTokenWithMissingKid(t *testing.T) {
	manager, err := NewJWKSManager(nil)
	if err != nil {
		t.Fatalf("NewJWKSManager() error = %v", err)
	}

	keyInfo, err := manager.GetActiveKey()
	if err != nil {
		t.Fatalf("GetActiveKey() error = %v", err)
	}

	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(keyInfo.PrivateKey)
	if err != nil {
		t.Fatalf("Failed to create RS256 token: %v", err)
	}

	_, err = manager.VerifyToken(tokenString)
	if err == nil {
		t.Error("VerifyToken() should fail with missing kid")
	}
}
