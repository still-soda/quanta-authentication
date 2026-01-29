package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppGroupRoleRepository struct {
	db *gorm.DB
}

func NewAppGroupRoleRepository(db *gorm.DB) *AppGroupRoleRepository {
	return &AppGroupRoleRepository{
		db: db,
	}
}

// WithTx 使用事务创建新的仓储实例
func (r *AppGroupRoleRepository) WithTx(tx *gorm.DB) *AppGroupRoleRepository {
	return &AppGroupRoleRepository{
		db: tx,
	}
}

// Create 创建应用组角色
func (r *AppGroupRoleRepository) Create(role *models.AppGroupRole) error {
	return r.db.Create(role).Error
}

// Update 更新应用组角色
func (r *AppGroupRoleRepository) Update(role *models.AppGroupRole) error {
	return r.db.Save(role).Error
}

// Delete 删除应用组角色
func (r *AppGroupRoleRepository) Delete(roleID string) error {
	return r.db.Delete(&models.AppGroupRole{}, "id = ?", roleID).Error
}

// FindByID 根据ID查找角色
func (r *AppGroupRoleRepository) FindByID(roleID string) (*models.AppGroupRole, error) {
	var role models.AppGroupRole
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByClientID 根据客户端ID获取所有角色
func (r *AppGroupRoleRepository) FindByClientID(clientID string) ([]models.AppGroupRole, error) {
	var roles []models.AppGroupRole
	if err := r.db.Where("client_id = ?", clientID).
		Order("is_default DESC, created_at ASC").
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindDefaultRole 获取默认角色
func (r *AppGroupRoleRepository) FindDefaultRole(clientID string) (*models.AppGroupRole, error) {
	var role models.AppGroupRole
	if err := r.db.Where("client_id = ? AND is_default = ?", clientID, true).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ClearDefaultRole 清除默认角色标志
func (r *AppGroupRoleRepository) ClearDefaultRole(clientID string) error {
	return r.db.Model(&models.AppGroupRole{}).
		Where("client_id = ? AND is_default = ?", clientID, true).
		Update("is_default", false).Error
}

// SetDefaultRole 设置默认角色
func (r *AppGroupRoleRepository) SetDefaultRole(roleID, clientID string) error {
	return r.db.Model(&models.AppGroupRole{}).
		Where("id = ? AND client_id = ?", roleID, clientID).
		Update("is_default", true).Error
}

// FindWithPermissions 查找角色及其权限
func (r *AppGroupRoleRepository) FindWithPermissions(roleID string) (*models.AppGroupRole, error) {
	var role models.AppGroupRole
	if err := r.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CountUsers 统计角色的用户数
func (r *AppGroupRoleRepository) CountUsers(roleID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.AppGroupUsersRoles{}).
		Where("app_group_role_id = ?", roleID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountPermissions 统计角色的权限数
func (r *AppGroupRoleRepository) CountPermissions(roleID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.AppGroupRolesPermissions{}).
		Where("app_group_role_id = ?", roleID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GrantPermissions 为角色分配权限
func (r *AppGroupRoleRepository) GrantPermissions(roleID string, permissionIDs []string) error {
	for _, permID := range permissionIDs {
		assoc := &models.AppGroupRolesPermissions{
			AppGroupRoleID:       roleID,
			AppGroupPermissionID: permID,
		}
		if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(assoc).Error; err != nil {
			return err
		}
	}
	return nil
}

// RevokePermissions 撤销角色权限
func (r *AppGroupRoleRepository) RevokePermissions(roleID string, permissionIDs []string) error {
	return r.db.Where("app_group_role_id = ? AND app_group_permission_id IN ?", roleID, permissionIDs).
		Delete(&models.AppGroupRolesPermissions{}).Error
}

// ClearPermissions 清空角色的所有权限
func (r *AppGroupRoleRepository) ClearPermissions(roleID string) error {
	return r.db.Where("app_group_role_id = ?", roleID).
		Delete(&models.AppGroupRolesPermissions{}).Error
}

// DeleteRolePermissionAssociations 删除角色权限关联
func (r *AppGroupRoleRepository) DeleteRolePermissionAssociations(roleID string) error {
	return r.db.Where("app_group_role_id = ?", roleID).
		Delete(&models.AppGroupRolesPermissions{}).Error
}

// DeleteUserRoleAssociations 删除用户角色关联
func (r *AppGroupRoleRepository) DeleteUserRoleAssociations(roleID string) error {
	return r.db.Where("app_group_role_id = ?", roleID).
		Delete(&models.AppGroupUsersRoles{}).Error
}

// DeletePermissionRoleAssociations 删除指定权限的角色关联
func (r *AppGroupRoleRepository) DeletePermissionRoleAssociations(permissionID string) error {
	return r.db.Where("app_group_permission_id = ?", permissionID).
		Delete(&models.AppGroupRolesPermissions{}).Error
}
