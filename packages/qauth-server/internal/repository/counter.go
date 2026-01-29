package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type CounterRepository struct {
	db *gorm.DB
}

// NewCounterRepository 创建计数器仓库
func NewCounterRepository(db *gorm.DB) *CounterRepository {
	return &CounterRepository{
		db: db,
	}
}

// FindByKeyWithinTimeRange 返回时间范围内的指定计数器
func (r *CounterRepository) FindByKeyWithinTimeRange(key string, startTime int64, endTime int64) ([]models.Counters, error) {
	var counters []models.Counters
	if err := r.db.
		Where("key = ? AND timestamp >= ? AND timestamp <= ?", key, startTime, endTime).
		Order("timestamp DESC").
		Find(&counters).Error; err != nil {
		return nil, err
	}
	return counters, nil
}

// Create 创建一个新的计数器记录
func (r *CounterRepository) Create(counter *models.Counters) error {
	return r.db.Create(counter).Error
}

// FindRecentByKey 获取最近指定数量的计数器记录
func (r *CounterRepository) FindRecentByKey(key string, limit int) ([]models.Counters, error) {
	var counters []models.Counters
	if err := r.db.
		Where("key = ?", key).
		Order("timestamp DESC").
		Limit(limit).
		Find(&counters).Error; err != nil {
		return nil, err
	}
	return counters, nil
}
