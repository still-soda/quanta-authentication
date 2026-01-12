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
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string // 完整连接字符串（优先级高于单独配置）
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// StorageConfig 存储配置
type StorageConfig struct {
	Provider string // local, s3
	LocalDir string // 本地存储目录
	S3Bucket string // S3 bucket 名称
	S3Region string // S3 区域
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret          string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
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

	return &Config{
		Server: ServerConfig{
			Port: GetEnv("PORT", "8080"),
			Mode: GetEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			DSN:      GetEnv("DATABASE_URL", ""),
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnvAsInt("DB_PORT", 5432),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", ""),
			DBName:   GetEnv("DB_NAME", "qauth"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     GetEnv("REDIS_ADDR", "localhost:6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       GetEnvAsInt("REDIS_DB", 0),
		},
		Storage: StorageConfig{
			Provider: GetEnv("STORAGE_PROVIDER", "local"),
			LocalDir: GetEnv("STORAGE_LOCAL_DIR", "./uploads"),
			S3Bucket: GetEnv("STORAGE_S3_BUCKET", ""),
			S3Region: GetEnv("STORAGE_S3_REGION", "us-east-1"),
		},
		JWT: JWTConfig{
			Secret:          GetEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenExp:  accessTokenExp,
			RefreshTokenExp: refreshTokenExp,
		},
	}
}
