package repository

import (
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type ErrorRecordRepository struct {
	db *gorm.DB
}

// NewErrorRecordRepository 创建错误记录仓库实例
func NewErrorRecordRepository(db *gorm.DB) *ErrorRecordRepository {
	return &ErrorRecordRepository{
		db: db,
	}
}

// Create 创建错误记录
func (r *ErrorRecordRepository) Create(record *models.ErrorRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		return e.ErrErrorRecordFailed.Wrap(err)
	}
	return nil
}
