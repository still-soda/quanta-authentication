package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type AppGroupPermissionRepository struct {
	db *gorm.DB
}

func NewAppGroupPermissionRepository(db *gorm.DB) *AppGroupPermissionRepository {
	return &AppGroupPermissionRepository{
		db: db,
	}
}

// Create 创建应用组权限
func (r *AppGroupPermissionRepository) Create(permission *models.AppGroupPermission) error {
	return r.db.Create(permission).Error
}

// Update 更新应用组权限
func (r *AppGroupPermissionRepository) Update(permission *models.AppGroupPermission) error {
	return r.db.Save(permission).Error
}

// Delete 删除应用组权限
func (r *AppGroupPermissionRepository) Delete(permissionID string) error {
	return r.db.Delete(&models.AppGroupPermission{}, "id = ?", permissionID).Error
}

// FindByID 根据ID查找权限
func (r *AppGroupPermissionRepository) FindByID(permissionID string) (*models.AppGroupPermission, error) {
	var permission models.AppGroupPermission
	if err := r.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByClientID 根据客户端ID获取所有权限
func (r *AppGroupPermissionRepository) FindByClientID(clientID string) ([]models.AppGroupPermission, error) {
	var permissions []models.AppGroupPermission
	if err := r.db.Where("client_id = ?", clientID).
		Order("resource ASC, action ASC").
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// FindByIDs 根据权限ID列表获取权限
func (r *AppGroupPermissionRepository) FindByIDs(permissionIDs []string) ([]models.AppGroupPermission, error) {
	var permissions []models.AppGroupPermission
	if err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// FindByCodes 根据权限代码获取权限列表
func (r *AppGroupPermissionRepository) FindByCodes(clientID string, codes []string) ([]models.AppGroupPermission, error) {
	var permissions []models.AppGroupPermission
	if err := r.db.Where("client_id = ? AND code IN ?", clientID, codes).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// FindByRoleIDs 根据角色ID列表获取权限
func (r *AppGroupPermissionRepository) FindByRoleIDs(roleIDs []string) ([]models.AppGroupPermission, error) {
	var permissions []models.AppGroupPermission
	if err := r.db.Distinct().
		Joins("JOIN app_group_roles_permissions ON app_group_permissions.id = app_group_roles_permissions.app_group_permission_id").
		Where("app_group_roles_permissions.app_group_role_id IN ?", roleIDs).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
