package services

import (
	"qauth-server/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	cfg   *config.Config
	redis *redis.Client
}

func NewCacheService(cfg *config.Config) *CacheService {
	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return &CacheService{
		cfg:   cfg,
		redis: redis,
	}
}

func (s *CacheService) GetRedisClient() *redis.Client {
	return s.redis
}

func (s *CacheService) Close() error {
	return s.redis.Close()
}

// SetKeyValue 设置键值对，exp 以秒为单位
func (s *CacheService) SetKeyValue(key string, value string, exp int) error {
	return s.redis.Set(s.redis.Context(), key, value, time.Duration(exp)*time.Second).Err()
}

// SetKeyExpire 设置键的过期时间，exp 以秒为单位
func (s *CacheService) SetKeyExpire(key string, exp int) error {
	return s.redis.Expire(s.redis.Context(), key, time.Duration(exp)*time.Second).Err()
}

// GetKeyValue 获取键值对
func (s *CacheService) GetKeyValue(key string) (string, error) {
	return s.redis.Get(s.redis.Context(), key).Result()
}

// DeleteKey 删除键值对
func (s *CacheService) DeleteKey(key string) error {
	return s.redis.Del(s.redis.Context(), key).Err()
}
