package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsersRolesRepository struct {
	db *gorm.DB
}

func NewUsersRolesRepository(db *gorm.DB) *UsersRolesRepository {
	return &UsersRolesRepository{
		db: db,
	}
}

// Create 创建用户角色关联
func (r *UsersRolesRepository) Create(userRole *models.UsersRoles) error {
	return r.db.Create(userRole).Error
}

// BatchCreate 批量创建用户角色关联（忽略冲突）
func (r *UsersRolesRepository) BatchCreate(userRoles []models.UsersRoles) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error
}

// DeleteByRoleID 删除角色的所有用户关联
func (r *UsersRolesRepository) DeleteByRoleID(roleID string) error {
	return r.db.Where("roles_id = ?", roleID).Delete(&models.UsersRoles{}).Error
}

// DeleteByUserIDAndRoleID 删除特定用户的角色关联
func (r *UsersRolesRepository) DeleteByUserIDAndRoleID(userID, roleID string) error {
	return r.db.Delete(&models.UsersRoles{}, "users_id = ? AND roles_id = ?", userID, roleID).Error
}

// FindByUserID 查询用户的系统角色
func (r *UsersRolesRepository) FindByUserID(userID string) (*models.UsersRoles, error) {
	var userRole models.UsersRoles
	if err := r.db.Preload("Role", "is_system = ?", true).
		First(&userRole, "users_id = ?", userID).
		Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}
