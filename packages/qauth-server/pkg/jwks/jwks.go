package jwks

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// KeyInfo 密钥信息
type KeyInfo struct {
	ID         string          // Key ID (kid)
	PrivateKey *rsa.PrivateKey // RSA 私钥
	PublicKey  *rsa.PublicKey  // RSA 公钥
	CreatedAt  time.Time       // 创建时间
	ExpiresAt  time.Time       // 过期时间
	Status     KeyStatus       // 密钥状态
}

// KeyStatus 密钥状态
type KeyStatus string

const (
	KeyStatusActive   KeyStatus = "active"   // 当前使用的密钥
	KeyStatusRotating KeyStatus = "rotating" // 正在轮换的密钥（仍可验证）
	KeyStatusExpired  KeyStatus = "expired"  // 已过期的密钥
)

// JWK JSON Web Key 结构
type JWK struct {
	Kty string `json:"kty"` // Key Type (RSA)
	Use string `json:"use"` // Key Use (sig)
	Alg string `json:"alg"` // Algorithm (RS256)
	Kid string `json:"kid"` // Key ID
	N   string `json:"n"`   // RSA modulus
	E   string `json:"e"`   // RSA exponent
}

// JWKS JSON Web Key Set 结构
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWKSManager JWKS 管理器
type JWKSManager struct {
	mu               sync.RWMutex
	keys             map[string]*KeyInfo // kid -> KeyInfo
	activeKeyID      string              // 当前活跃的密钥 ID
	keySize          int                 // RSA 密钥大小
	rotationInterval time.Duration       // 密钥轮换间隔
	gracePeriod      time.Duration       // 旧密钥保留时间（用于验证）
	stopChan         chan struct{}       // 停止信号
}

// JWKSConfig JWKS 配置
type JWKSConfig struct {
	KeySize          int           // RSA 密钥大小，默认 2048
	RotationInterval time.Duration // 密钥轮换间隔，默认 24 小时
	GracePeriod      time.Duration // 旧密钥保留时间，默认 1 小时
}

// DefaultJWKSConfig 默认配置
func DefaultJWKSConfig() *JWKSConfig {
	return &JWKSConfig{
		KeySize:          2048,
		RotationInterval: 24 * time.Hour,
		GracePeriod:      1 * time.Hour,
	}
}

// NewJWKSManager 创建新的 JWKS 管理器
func NewJWKSManager(cfg *JWKSConfig) (*JWKSManager, error) {
	if cfg == nil {
		cfg = DefaultJWKSConfig()
	}

	manager := &JWKSManager{
		keys:             make(map[string]*KeyInfo),
		keySize:          cfg.KeySize,
		rotationInterval: cfg.RotationInterval,
		gracePeriod:      cfg.GracePeriod,
		stopChan:         make(chan struct{}),
	}

	// 生成初始密钥
	if err := manager.generateNewKey(); err != nil {
		return nil, err
	}

	return manager, nil
}

// generateNewKey 生成新的 RSA 密钥对
func (m *JWKSManager) generateNewKey() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, m.keySize)
	if err != nil {
		return err
	}

	kid := uuid.New().String()
	now := time.Now()

	keyInfo := &KeyInfo{
		ID:         kid,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		CreatedAt:  now,
		ExpiresAt:  now.Add(m.rotationInterval + m.gracePeriod),
		Status:     KeyStatusActive,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 将当前活跃密钥标记为轮换中
	if m.activeKeyID != "" {
		if oldKey, exists := m.keys[m.activeKeyID]; exists {
			oldKey.Status = KeyStatusRotating
			oldKey.ExpiresAt = now.Add(m.gracePeriod)
		}
	}

	m.keys[kid] = keyInfo
	m.activeKeyID = kid

	return nil
}

// StartRotation 启动密钥轮换
func (m *JWKSManager) StartRotation() {
	go func() {
		ticker := time.NewTicker(m.rotationInterval)
		defer ticker.Stop()

		cleanupTicker := time.NewTicker(m.gracePeriod / 2)
		defer cleanupTicker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.generateNewKey(); err != nil {
					// 记录错误但继续运行
					continue
				}
			case <-cleanupTicker.C:
				m.cleanupExpiredKeys()
			case <-m.stopChan:
				return
			}
		}
	}()
}

// StopRotation 停止密钥轮换
func (m *JWKSManager) StopRotation() {
	close(m.stopChan)
}

// cleanupExpiredKeys 清理过期的密钥
func (m *JWKSManager) cleanupExpiredKeys() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for kid, keyInfo := range m.keys {
		if keyInfo.Status == KeyStatusRotating && now.After(keyInfo.ExpiresAt) {
			keyInfo.Status = KeyStatusExpired
			delete(m.keys, kid)
		}
	}
}

// GetActiveKey 获取当前活跃的密钥
func (m *JWKSManager) GetActiveKey() (*KeyInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeKeyID == "" {
		return nil, errors.New("no active key available")
	}

	keyInfo, exists := m.keys[m.activeKeyID]
	if !exists {
		return nil, errors.New("active key not found")
	}

	return keyInfo, nil
}

// GetKeyByID 根据 ID 获取密钥
func (m *JWKSManager) GetKeyByID(kid string) (*KeyInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keyInfo, exists := m.keys[kid]
	if !exists {
		return nil, errors.New("key not found")
	}

	if keyInfo.Status == KeyStatusExpired {
		return nil, errors.New("key has expired")
	}

	return keyInfo, nil
}

// GetJWKS 获取 JWKS（公钥集合）
func (m *JWKSManager) GetJWKS() *JWKS {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jwks := &JWKS{
		Keys: make([]JWK, 0),
	}

	for _, keyInfo := range m.keys {
		if keyInfo.Status == KeyStatusExpired {
			continue
		}

		jwk := publicKeyToJWK(keyInfo.ID, keyInfo.PublicKey)
		jwks.Keys = append(jwks.Keys, jwk)
	}

	return jwks
}

// GetJWKSJSON 获取 JWKS JSON 字符串
func (m *JWKSManager) GetJWKSJSON() ([]byte, error) {
	jwks := m.GetJWKS()
	return json.Marshal(jwks)
}

// publicKeyToJWK 将 RSA 公钥转换为 JWK
func publicKeyToJWK(kid string, publicKey *rsa.PublicKey) JWK {
	return JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: kid,
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}
}

// SignToken 使用当前活跃密钥签名 JWT
func (m *JWKSManager) SignToken(claims jwt.Claims) (string, error) {
	keyInfo, err := m.GetActiveKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyInfo.ID

	return token.SignedString(keyInfo.PrivateKey)
}

// VerifyToken 验证 JWT 签名
func (m *JWKSManager) VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// 获取 kid
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}

		// 获取对应的公钥
		keyInfo, err := m.GetKeyByID(kid)
		if err != nil {
			return nil, err
		}

		return keyInfo.PublicKey, nil
	})
}

// ForceRotate 强制进行密钥轮换
func (m *JWKSManager) ForceRotate() error {
	return m.generateNewKey()
}

// GetAllKeys 获取所有密钥信息（不含私钥）
func (m *JWKSManager) GetAllKeys() []KeyInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]KeyInfo, 0, len(m.keys))
	for _, keyInfo := range m.keys {
		// 返回不含私钥的副本
		keys = append(keys, KeyInfo{
			ID:        keyInfo.ID,
			PublicKey: keyInfo.PublicKey,
			CreatedAt: keyInfo.CreatedAt,
			ExpiresAt: keyInfo.ExpiresAt,
			Status:    keyInfo.Status,
		})
	}

	return keys
}

// GetActiveKeyID 获取当前活跃密钥的 ID
func (m *JWKSManager) GetActiveKeyID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeKeyID
}
