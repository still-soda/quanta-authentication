package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓库
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

// Create 创建文件记录
func (r *FileRepository) Create(file *models.Files) error {
	return r.db.Create(file).Error
}
