package services

import (
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/repository"
)

type CounterService struct {
	repo *repository.CounterRepository
}

// NewCounterService 创建计数器服务
func NewCounterService(repo *repository.CounterRepository) *CounterService {
	return &CounterService{repo: repo}
}

// GetCountersWithinTimeRange 返回时间范围内的指定计数器
func (s *CounterService) GetCountersWithinTimeRange(key string, startTime int64, endTime int64) ([]models.Counters, error) {
	return s.repo.FindByKeyWithinTimeRange(key, startTime, endTime)
}

// CreateCounter 创建一个新的计数器记录
func (s *CounterService) CreateCounter(key string, count int64) error {
	counter := &models.Counters{
		Key:   key,
		Count: count,
	}
	return s.repo.Create(counter)
}

var ZERO int64 = 0

// GetRecentCounter 获取最近指定数量的计数器记录
// padNum: 如果记录数不足，则用该值进行填充；如果为 nil，则使用第一个记录的值进行填充
func (s *CounterService) GetRecentCounter(key string, limit int, padNum *int64) ([]models.Counters, error) {
	counters, err := s.repo.FindRecentByKey(key, limit)
	if err != nil {
		return nil, e.ErrFailedToFindCounters.Wrap(err)
	}

	// 首部填充到指定长度
	diff := limit - len(counters)
	if diff > 0 {
		if padNum == nil {
			if len(counters) == 0 {
				padNum = &ZERO
			} else {
				head := counters[0].Count
				padNum = &head
			}
		}
		padding := make([]models.Counters, limit)
		for i := range diff {
			padding[i] = models.Counters{
				Key:   key,
				Count: *padNum,
			}
		}
		counters = append(padding, counters...)
	}

	return counters, nil
}
