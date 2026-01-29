package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppGroupAdminRepository struct {
	db *gorm.DB
}

func NewAppGroupAdminRepository(db *gorm.DB) *AppGroupAdminRepository {
	return &AppGroupAdminRepository{
		db: db,
	}
}

// Create 创建应用组管理员记录，若已存在则忽略
func (r *AppGroupAdminRepository) Create(admin *models.AppGroupAdmin) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(admin).Error
}

// Delete 删除应用组管理员
func (r *AppGroupAdminRepository) Delete(clientID, userID string, adminType models.AppGroupAdminType) error {
	return r.db.Where("client_id = ? AND user_id = ? AND admin_type = ?", clientID, userID, adminType).
		Delete(&models.AppGroupAdmin{}).Error
}

// FindByClientID 获取应用组的所有管理员
func (r *AppGroupAdminRepository) FindByClientID(clientID string) ([]models.AppGroupAdmin, error) {
	var admins []models.AppGroupAdmin
	if err := r.db.Preload("User").Preload("Granter").
		Where("client_id = ?", clientID).
		Order("admin_type ASC, created_at ASC").
		Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// FindByUserID 获取用户在所有应用组的管理员身份
func (r *AppGroupAdminRepository) FindByUserID(userID string) ([]models.AppGroupAdmin, error) {
	var admins []models.AppGroupAdmin
	if err := r.db.Preload("Client").
		Where("user_id = ?", userID).
		Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// FindByClientIDAndUserID 检查用户是否为应用组管理员
func (r *AppGroupAdminRepository) FindByClientIDAndUserID(clientID, userID string) (*models.AppGroupAdmin, error) {
	var admin models.AppGroupAdmin
	err := r.db.Where("client_id = ? AND user_id = ?", clientID, userID).
		Order("CASE admin_type WHEN 'owner' THEN 1 WHEN 'admin' THEN 2 ELSE 3 END").
		First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
