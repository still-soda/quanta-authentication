package config

import (
	"log"
	"os"
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

// Config 配置结构体
type Config struct {
	// 数据库连接字符串
	DatabaseURL       string
	//  Redis 连接地址
	RedisURL          string
	// Redis 连接密码
	RedisPassword     string
	// 访问令牌过期时间
	AccessTokenExpire time.Duration
	// 服务器端口
	Port              string
}

func New() *Config {
	LoadConfig()

	// 读取 ACCESS_TOKEN_EXPIRE 环境变量，默认值为 3600 秒
	ACCESS_TOKEN_EXPIRE := GetEnv("ACCESS_TOKEN_EXPIRE", "3600")
	accessTokenExp, err := time.ParseDuration(ACCESS_TOKEN_EXPIRE + "s")
	if err != nil {
		accessTokenExp = time.Hour
	}

	return &Config{
		DatabaseURL:       GetEnv("DATABASE_URL", "postgres://user:password@localhost:5432/qauth_db"),
		RedisURL:          GetEnv("REDIS_URL", "127.0.0.1:6379"),
		RedisPassword:     GetEnv("REDIS_PASSWORD", ""),
		AccessTokenExpire: accessTokenExp,
		Port:              GetEnv("PORT", "8080"),
	}
}