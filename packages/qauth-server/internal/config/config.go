package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadConfig 加载 .env 文件到系统环境变量
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[INFO] 未找到 .env 文件，将使用系统环境变量")
	}
}

// GetEnv 获取变量，带默认值
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvAsInt 获取整数环境变量
func GetEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// Config 配置结构体
type Config struct {
	// 服务器配置
	Server ServerConfig

	// 数据库配置
	Database DatabaseConfig

	// Redis 配置
	Redis RedisConfig

	// 存储配置
	Storage StorageConfig

	// JWT 配置
	JWT JWTConfig

	// OIDC 配置
	OIDC OIDCConfig

	// OAuth2 配置
	OAuth OAuthConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// StorageConfig 存储配置
type StorageConfig struct {
	LocalDir string // 本地存储目录
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret          string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// OIDCConfig OIDC 配置
type OIDCConfig struct {
	Issuer              string        // 发行者 URL
	KeyRotationInterval time.Duration // 密钥轮换间隔
	KeySize             int           // RSA 密钥大小
}

type OAuthConfig struct {
	AllowedResponseTypes              []ResponseType
	ScopeSupported                    []Scope
	ResponseModesSupported            []ResponseMode
	GrantTypesSupported               []GrantType
	SubjectTypesSupported             []SubjectType
	IDTokenSigningAlgValuesSupported  []IDTokenSigningAlg
	TokenEndpointAuthMethodsSupported []TokenEndpointAuthMethod
	ClaimsSupported                   []Claim
	CodeChallengeMethodsSupported     []CodeChallengeMethod
}

func New() *Config {
	LoadConfig()

	// JWT 配置
	accessTokenExpStr := GetEnv("ACCESS_TOKEN_EXPIRE", "3600")
	accessTokenExp, err := time.ParseDuration(accessTokenExpStr + "s")
	if err != nil {
		accessTokenExp = time.Hour
	}

	refreshTokenExpStr := GetEnv("REFRESH_TOKEN_EXPIRE", "604800")
	refreshTokenExp, err := time.ParseDuration(refreshTokenExpStr + "s")
	if err != nil {
		refreshTokenExp = 7 * 24 * time.Hour
	}

	// OIDC 密钥轮换间隔
	keyRotationIntervalStr := GetEnv("OIDC_KEY_ROTATION_INTERVAL", "86400")
	keyRotationInterval, err := time.ParseDuration(keyRotationIntervalStr + "s")
	if err != nil {
		keyRotationInterval = 24 * time.Hour
	}

	return &Config{
		Server: ServerConfig{
			Port: GetEnv("PORT", "8080"),
			Mode: GetEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			DSN: GetEnv("DATABASE_URL", ""),
		},
		Redis: RedisConfig{
			Addr:     GetEnv("REDIS_ADDR", "localhost:6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       GetEnvAsInt("REDIS_DB", 0),
		},
		Storage: StorageConfig{
			LocalDir: GetEnv("STORAGE_LOCAL_DIR", "./uploads"),
		},
		JWT: JWTConfig{
			Secret:          GetEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenExp:  accessTokenExp,
			RefreshTokenExp: refreshTokenExp,
		},
		OIDC: OIDCConfig{
			Issuer:              GetEnv("OIDC_ISSUER", "http://localhost:8080"),
			KeyRotationInterval: keyRotationInterval,
			KeySize:             GetEnvAsInt("OIDC_KEY_SIZE", 2048),
		},
		OAuth: OAuthConfig{
			AllowedResponseTypes:              []ResponseType{ResponseTypeCode, ResponseTypeIDToken},
			ScopeSupported:                    []Scope{ScopeOpenID, ScopeProfile, ScopeEmail},
			ResponseModesSupported:            []ResponseMode{ResponseModeQuery, ResponseModeFragment, ResponseModeFormPost},
			GrantTypesSupported:               []GrantType{GrantTypeAuthorizationCode, GrantTypeRefreshToken},
			SubjectTypesSupported:             []SubjectType{SubjectTypePublic},
			IDTokenSigningAlgValuesSupported:  []IDTokenSigningAlg{IDTokenSigningAlgRS256},
			TokenEndpointAuthMethodsSupported: []TokenEndpointAuthMethod{TokenEndpointAuthMethodClientSecretBasic, TokenEndpointAuthMethodClientSecretPost},
			ClaimsSupported: []Claim{
				ClaimSub, ClaimIss, ClaimAud, ClaimExp, ClaimIat,
				ClaimName, ClaimStudentID, ClaimDisplayName, ClaimPicture,
				ClaimEmail, ClaimEmailVerified,
				ClaimRoles,
			},
			CodeChallengeMethodsSupported: []CodeChallengeMethod{CodeChallengeMethodS256, CodeChallengeMethodPlain},
		},
	}
}
