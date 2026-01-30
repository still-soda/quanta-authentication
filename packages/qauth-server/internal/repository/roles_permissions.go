package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type RolesPermissionsRepository struct {
	db *gorm.DB
}

func NewRolesPermissionsRepository(db *gorm.DB) *RolesPermissionsRepository {
	return &RolesPermissionsRepository{
		db: db,
	}
}

// DeleteByRoleID 删除角色的所有权限关联
func (r *RolesPermissionsRepository) DeleteByRoleID(roleID string) error {
	return r.db.Where("roles_id = ?", roleID).Delete(&models.RolesPermissions{}).Error
}

// FindPermissionsByAction 根据操作类型查询权限
func (r *RolesPermissionsRepository) FindPermissionsByAction(action int8) ([]models.Permissions, error) {
	var permissions []models.Permissions
	if err := r.db.Where("action = ?", action).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
