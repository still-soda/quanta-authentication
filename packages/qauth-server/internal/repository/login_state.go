package repository

import (
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type LoginStateRepository struct {
	db *gorm.DB
}

// NewLoginStateRepository 创建登录状态仓库实例
func NewLoginStateRepository(db *gorm.DB) *LoginStateRepository {
	return &LoginStateRepository{
		db: db,
	}
}

// Create 创建登录状态记录
func (r *LoginStateRepository) Create(state *models.LoginState) error {
	if err := r.db.Create(state).Error; err != nil {
		return e.ErrLoginStateRecordFailed.Wrap(err)
	}
	return nil
}
