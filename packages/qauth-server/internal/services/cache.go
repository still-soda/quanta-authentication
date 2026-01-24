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
	val, err := s.redis.Get(s.redis.Context(), key).Result()

	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

// GetKeyValueAsInt64 获取键值对并转换为int64
func (s *CacheService) GetKeyValueAsInt64(key string) (int64, error) {
	val, err := s.redis.Get(s.redis.Context(), key).Int64()

	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return val, nil
}

// DeleteKey 删除键值对
func (s *CacheService) DeleteKey(key string) error {
	return s.redis.Del(s.redis.Context(), key).Err()
}

// IncrKey 自增键的值
func (s *CacheService) IncrKey(key string) (int64, error) {
	return s.redis.Incr(s.redis.Context(), key).Result()
}

// CountPrefixKeys 计算具有指定前缀的键的数量
func (s *CacheService) CountPrefixKeys(prefix string) (int64, error) {
	var cursor uint64
	var total int64 = 0

	const batchSize = 1000

	for {
		// 使用带有超时的 Context，防止扫描时间过长导致请求堆积
		keys, nextCursor, err := s.redis.Scan(s.redis.Context(), cursor, prefix+"*", batchSize).Result()
		if err != nil {
			return total, err
		}

		total += int64(len(keys))
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return total, nil
}
