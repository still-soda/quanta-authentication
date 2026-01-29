package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppGroupUserRoleRepository struct {
	db *gorm.DB
}

func NewAppGroupUserRoleRepository(db *gorm.DB) *AppGroupUserRoleRepository {
	return &AppGroupUserRoleRepository{
		db: db,
	}
}

// Assign 为用户分配角色
func (r *AppGroupUserRoleRepository) Assign(assignment *models.AppGroupUsersRoles) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(assignment).Error
}

// Revoke 撤销用户角色
func (r *AppGroupUserRoleRepository) Revoke(userID, roleID string) error {
	return r.db.Where("user_id = ? AND app_group_role_id = ?", userID, roleID).
		Delete(&models.AppGroupUsersRoles{}).Error
}

// FindUserRoles 获取用户在应用组的所有角色
func (r *AppGroupUserRoleRepository) FindUserRoles(userID, clientID string) ([]models.AppGroupUsersRoles, error) {
	var assignments []models.AppGroupUsersRoles
	if err := r.db.Preload("AppGroupRole").
		Joins("JOIN app_group_roles ON app_group_users_roles.app_group_role_id = app_group_roles.id").
		Where("app_group_users_roles.user_id = ? AND app_group_roles.client_id = ?", userID, clientID).
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

// FindRoleUsers 获取角色的所有用户
func (r *AppGroupUserRoleRepository) FindRoleUsers(roleID string) ([]models.AppGroupUsersRoles, error) {
	var assignments []models.AppGroupUsersRoles
	if err := r.db.Preload("User").
		Where("app_group_role_id = ?", roleID).
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}
