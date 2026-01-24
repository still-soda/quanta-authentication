package services

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type CounterService struct {
	db *gorm.DB
}

func NewCounterService(db *gorm.DB) *CounterService {
	return &CounterService{db: db}
}

// GetCountersWithinTimeRange 返回时间范围内的指定计数器
func (s *CounterService) GetCountersWithinTimeRange(key string, startTime int64, endTime int64) ([]models.Counters, error) {
	var counters []models.Counters
	if err := s.db.
		Where("key = ? AND timestamp >= ? AND timestamp <= ?", key, startTime, endTime).
		Order("timestamp DESC").
		Find(&counters).Error; err != nil {
		return nil, err
	}
	return counters, nil
}

// CreateCounter 创建一个新的计数器记录
func (s *CounterService) CreateCounter(key string, count int64) error {
	counter := &models.Counters{
		Key:   key,
		Count: count,
	}
	return s.db.Create(counter).Error
}

var ZERO int64 = 0

// GetRecentCounter 获取最近指定数量的计数器记录
// padNum: 如果记录数不足，则用该值进行填充；如果为 nil，则使用第一个记录的值进行填充
func (s *CounterService) GetRecentCounter(key string, limit int, padNum *int64) ([]models.Counters, error) {
	var counters []models.Counters
	if err := s.db.
		Where("key = ?", key).
		Order("timestamp DESC").
		Limit(limit).
		Find(&counters).Error; err != nil {
		return nil, err
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
