package services

import (
	"errors"
	"fmt"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AppGroupService 应用组权限服务
type AppGroupService struct {
	db          *gorm.DB
	userService *UserService
}

// NewAppGroupService 创建应用组权限服务
func NewAppGroupService(db *gorm.DB, userService *UserService) *AppGroupService {
	return &AppGroupService{
		db:          db,
		userService: userService,
	}
}

// ======================== 应用组管理员相关方法 ========================

// InitializeAppGroupAdmins 初始化应用组管理员（创建应用时调用）
// 自动添加应用创建者为 owner，超级管理员也拥有管理权限
func (s *AppGroupService) InitializeAppGroupAdmins(clientID, creatorID string) error {
	// 添加创建者为 owner
	ownerAdmin := &models.AppGroupAdmin{
		ClientID:  clientID,
		UserID:    creatorID,
		AdminType: models.AppGroupAdminTypeOwner,
		GrantedBy: creatorID,
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(ownerAdmin).Error; err != nil {
		return fmt.Errorf("failed to create owner admin: %w", err)
	}

	utilities.GetLogger().Info("initialized app group admins", "clientID", clientID, "creatorID", creatorID)
	return nil
}

// AddAppGroupAdmin 添加应用组管理员
func (s *AppGroupService) AddAppGroupAdmin(clientID, userID string, adminType models.AppGroupAdminType, grantedBy string) error {
	admin := &models.AppGroupAdmin{
		ClientID:  clientID,
		UserID:    userID,
		AdminType: adminType,
		GrantedBy: grantedBy,
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(admin).Error; err != nil {
		return fmt.Errorf("failed to add app group admin: %w", err)
	}

	return nil
}

// RemoveAppGroupAdmin 移除应用组管理员
func (s *AppGroupService) RemoveAppGroupAdmin(clientID, userID string, adminType models.AppGroupAdminType) error {
	result := s.db.Where("client_id = ? AND user_id = ? AND admin_type = ?", clientID, userID, adminType).
		Delete(&models.AppGroupAdmin{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove app group admin: %w", result.Error)
	}
	return nil
}

// GetAppGroupAdmins 获取应用组的所有管理员
func (s *AppGroupService) GetAppGroupAdmins(clientID string) ([]models.AppGroupAdmin, error) {
	var admins []models.AppGroupAdmin
	if err := s.db.Preload("User").Preload("Granter").
		Where("client_id = ?", clientID).
		Order("admin_type ASC, created_at ASC").
		Find(&admins).Error; err != nil {
		return nil, fmt.Errorf("failed to get app group admins: %w", err)
	}
	return admins, nil
}

// GetAppGroupAdminsByUser 获取用户在所有应用组的管理员身份
func (s *AppGroupService) GetAppGroupAdminsByUser(userID string) ([]models.AppGroupAdmin, error) {
	var admins []models.AppGroupAdmin
	if err := s.db.Preload("Client").
		Where("user_id = ?", userID).
		Find(&admins).Error; err != nil {
		return nil, fmt.Errorf("failed to get user's app group admins: %w", err)
	}
	return admins, nil
}

// IsAppGroupAdmin 检查用户是否为应用组管理员
func (s *AppGroupService) IsAppGroupAdmin(clientID, userID string) (bool, models.AppGroupAdminType, error) {
	var admin models.AppGroupAdmin
	err := s.db.Where("client_id = ? AND user_id = ?", clientID, userID).
		Order("CASE admin_type WHEN 'owner' THEN 1 WHEN 'admin' THEN 2 ELSE 3 END").
		First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check app group admin: %w", err)
	}
	return true, admin.AdminType, nil
}

// HasAppGroupAdminPermission 检查用户是否有应用组管理权限
func (s *AppGroupService) HasAppGroupAdminPermission(clientID, userID string, requiredTypes ...models.AppGroupAdminType) (bool, error) {
	isAdmin, adminType, err := s.IsAppGroupAdmin(clientID, userID)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, nil
	}

	// owner 和 admin 拥有所有权限
	if adminType == models.AppGroupAdminTypeOwner || adminType == models.AppGroupAdminTypeAdmin {
		return true, nil
	}

	// 检查是否有需要的特定权限
	for _, rt := range requiredTypes {
		if adminType == rt {
			return true, nil
		}
	}

	return false, nil
}

// ======================== 应用组权限相关方法 ========================

// CreateAppGroupPermission 创建应用组权限
func (s *AppGroupService) CreateAppGroupPermission(clientID, resource string, action int8, code, name, description string) (*models.AppGroupPermission, error) {
	permission := &models.AppGroupPermission{
		ClientID:    clientID,
		Resource:    resource,
		Action:      action,
		Code:        code,
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(permission).Error; err != nil {
		return nil, fmt.Errorf("failed to create app group permission: %w", err)
	}

	return permission, nil
}

// UpdateAppGroupPermission 更新应用组权限
func (s *AppGroupService) UpdateAppGroupPermission(permissionID, name, description string) (*models.AppGroupPermission, error) {
	var permission models.AppGroupPermission
	if err := s.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return nil, fmt.Errorf("failed to find permission: %w", err)
	}

	permission.Name = name
	permission.Description = description

	if err := s.db.Save(&permission).Error; err != nil {
		return nil, fmt.Errorf("failed to update app group permission: %w", err)
	}

	return &permission, nil
}

// DeleteAppGroupPermission 删除应用组权限
func (s *AppGroupService) DeleteAppGroupPermission(permissionID string) error {
	// 先删除角色权限关联
	if err := s.db.Where("app_group_permission_id = ?", permissionID).Delete(&models.AppGroupRolesPermissions{}).Error; err != nil {
		return fmt.Errorf("failed to delete role permission associations: %w", err)
	}

	// 删除权限
	if err := s.db.Delete(&models.AppGroupPermission{}, "id = ?", permissionID).Error; err != nil {
		return fmt.Errorf("failed to delete app group permission: %w", err)
	}

	return nil
}

// GetAppGroupPermission 获取应用组权限详情
func (s *AppGroupService) GetAppGroupPermission(permissionID string) (*models.AppGroupPermission, error) {
	var permission models.AppGroupPermission
	if err := s.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return nil, fmt.Errorf("failed to get app group permission: %w", err)
	}
	return &permission, nil
}

// GetAppGroupPermissions 获取应用组的所有权限
func (s *AppGroupService) GetAppGroupPermissions(clientID string) ([]models.AppGroupPermission, error) {
	var permissions []models.AppGroupPermission
	if err := s.db.Where("client_id = ?", clientID).
		Order("resource ASC, action ASC").
		Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get app group permissions: %w", err)
	}
	return permissions, nil
}

// GetAppGroupPermissionsByResource 获取应用组按资源分组的权限
func (s *AppGroupService) GetAppGroupPermissionsByResource(clientID string) (map[string][]models.AppGroupPermission, error) {
	permissions, err := s.GetAppGroupPermissions(clientID)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.AppGroupPermission)
	for _, perm := range permissions {
		grouped[perm.Resource] = append(grouped[perm.Resource], perm)
	}

	return grouped, nil
}

// GetAppGroupPermissionByCodes 根据权限代码获取权限列表
func (s *AppGroupService) GetAppGroupPermissionByCodes(clientID string, codes []string) ([]models.AppGroupPermission, error) {
	var permissions []models.AppGroupPermission
	if err := s.db.Where("client_id = ? AND code IN ?", clientID, codes).Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get app group permissions by codes: %w", err)
	}
	return permissions, nil
}

// ======================== 应用组角色相关方法 ========================

// CreateAppGroupRole 创建应用组角色
func (s *AppGroupService) CreateAppGroupRole(clientID, code, name, description string, isDefault bool) (*models.AppGroupRole, error) {
	role := &models.AppGroupRole{
		ClientID:    clientID,
		Code:        code,
		Name:        name,
		Description: description,
		IsDefault:   isDefault,
	}

	if err := s.db.Create(role).Error; err != nil {
		return nil, fmt.Errorf("failed to create app group role: %w", err)
	}

	return role, nil
}

// UpdateAppGroupRole 更新应用组角色
func (s *AppGroupService) UpdateAppGroupRole(roleID, name, description string, isDefault bool) (*models.AppGroupRole, error) {
	var role models.AppGroupRole
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}

	role.Name = name
	role.Description = description
	role.IsDefault = isDefault

	if err := s.db.Save(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to update app group role: %w", err)
	}

	return &role, nil
}

// DeleteAppGroupRole 删除应用组角色
func (s *AppGroupService) DeleteAppGroupRole(roleID string) error {
	// 先删除角色权限关联
	if err := s.db.Where("app_group_role_id = ?", roleID).Delete(&models.AppGroupRolesPermissions{}).Error; err != nil {
		return fmt.Errorf("failed to delete role permission associations: %w", err)
	}

	// 删除用户角色关联
	if err := s.db.Where("app_group_role_id = ?", roleID).Delete(&models.AppGroupUsersRoles{}).Error; err != nil {
		return fmt.Errorf("failed to delete user role associations: %w", err)
	}

	// 删除角色
	if err := s.db.Delete(&models.AppGroupRole{}, "id = ?", roleID).Error; err != nil {
		return fmt.Errorf("failed to delete app group role: %w", err)
	}

	return nil
}

// GetAppGroupRole 获取应用组角色详情
func (s *AppGroupService) GetAppGroupRole(roleID string) (*models.AppGroupRole, error) {
	var role models.AppGroupRole
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, fmt.Errorf("failed to get app group role: %w", err)
	}
	return &role, nil
}

// GetAppGroupRoles 获取应用组的所有角色
func (s *AppGroupService) GetAppGroupRoles(clientID string) ([]models.AppGroupRole, error) {
	var roles []models.AppGroupRole
	if err := s.db.Where("client_id = ?", clientID).
		Order("is_default DESC, created_at ASC").
		Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get app group roles: %w", err)
	}
	return roles, nil
}

// AppGroupRoleWithStats 应用组角色带统计信息
type AppGroupRoleWithStats struct {
	models.AppGroupRole
	UserCount       int64 `json:"user_count"`
	PermissionCount int64 `json:"permission_count"`
}

// GetAppGroupRolesWithStats 获取应用组的所有角色（带统计信息）
func (s *AppGroupService) GetAppGroupRolesWithStats(clientID string) ([]AppGroupRoleWithStats, error) {
	roles, err := s.GetAppGroupRoles(clientID)
	if err != nil {
		return nil, err
	}

	result := make([]AppGroupRoleWithStats, len(roles))
	for i, role := range roles {
		var userCount int64
		s.db.Model(&models.AppGroupUsersRoles{}).Where("app_group_role_id = ?", role.ID).Count(&userCount)

		var permCount int64
		s.db.Model(&models.AppGroupRolesPermissions{}).Where("app_group_role_id = ?", role.ID).Count(&permCount)

		result[i] = AppGroupRoleWithStats{
			AppGroupRole:    role,
			UserCount:       userCount,
			PermissionCount: permCount,
		}
	}

	return result, nil
}

// GetAppGroupRoleWithStats 获取单个应用组角色（带统计信息）
func (s *AppGroupService) GetAppGroupRoleWithStats(roleID string) (*AppGroupRoleWithStats, error) {
	role, err := s.GetAppGroupRole(roleID)
	if err != nil {
		return nil, err
	}

	var userCount int64
	s.db.Model(&models.AppGroupUsersRoles{}).Where("app_group_role_id = ?", role.ID).Count(&userCount)

	var permCount int64
	s.db.Model(&models.AppGroupRolesPermissions{}).Where("app_group_role_id = ?", role.ID).Count(&permCount)

	return &AppGroupRoleWithStats{
		AppGroupRole:    *role,
		UserCount:       userCount,
		PermissionCount: permCount,
	}, nil
}

// GrantPermissionsToAppGroupRole 为应用组角色分配权限
func (s *AppGroupService) GrantPermissionsToAppGroupRole(roleID string, permissionIDs []string) error {
	var role models.AppGroupRole
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	var permissions []models.AppGroupPermission
	if err := s.db.Where("id IN ? AND client_id = ?", permissionIDs, role.ClientID).Find(&permissions).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	// 创建关联
	for _, perm := range permissions {
		assoc := &models.AppGroupRolesPermissions{
			AppGroupRoleID:       roleID,
			AppGroupPermissionID: perm.ID,
		}
		if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(assoc).Error; err != nil {
			return fmt.Errorf("failed to grant permission to role: %w", err)
		}
	}

	return nil
}

// RevokePermissionsFromAppGroupRole 从应用组角色撤销权限
func (s *AppGroupService) RevokePermissionsFromAppGroupRole(roleID string, permissionIDs []string) error {
	if err := s.db.Where("app_group_role_id = ? AND app_group_permission_id IN ?", roleID, permissionIDs).
		Delete(&models.AppGroupRolesPermissions{}).Error; err != nil {
		return fmt.Errorf("failed to revoke permissions from role: %w", err)
	}
	return nil
}

// SetAppGroupRolePermissions 设置应用组角色的权限（替换现有权限）
func (s *AppGroupService) SetAppGroupRolePermissions(roleID string, permissionIDs []string) error {
	var role models.AppGroupRole
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有关联
		if err := tx.Where("app_group_role_id = ?", roleID).Delete(&models.AppGroupRolesPermissions{}).Error; err != nil {
			return fmt.Errorf("failed to clear existing permissions: %w", err)
		}

		// 添加新关联
		for _, permID := range permissionIDs {
			assoc := &models.AppGroupRolesPermissions{
				AppGroupRoleID:       roleID,
				AppGroupPermissionID: permID,
			}
			if err := tx.Create(assoc).Error; err != nil {
				return fmt.Errorf("failed to grant permission to role: %w", err)
			}
		}

		return nil
	})
}

// GetAppGroupRolePermissions 获取应用组角色的权限
func (s *AppGroupService) GetAppGroupRolePermissions(roleID string) ([]models.AppGroupPermission, error) {
	var role models.AppGroupRole
	if err := s.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	return role.Permissions, nil
}

// ======================== 应用组用户角色相关方法 ========================

// AssignAppGroupRoleToUser 为用户分配应用组角色
func (s *AppGroupService) AssignAppGroupRoleToUser(userID, roleID, assignedBy string) error {
	assignment := &models.AppGroupUsersRoles{
		UserID:         userID,
		AppGroupRoleID: roleID,
		AssignedBy:     assignedBy,
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(assignment).Error; err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// RevokeAppGroupRoleFromUser 从用户撤销应用组角色
func (s *AppGroupService) RevokeAppGroupRoleFromUser(userID, roleID string) error {
	if err := s.db.Where("user_id = ? AND app_group_role_id = ?", userID, roleID).
		Delete(&models.AppGroupUsersRoles{}).Error; err != nil {
		return fmt.Errorf("failed to revoke role from user: %w", err)
	}
	return nil
}

// GetUserAppGroupRoles 获取用户在应用组的所有角色
func (s *AppGroupService) GetUserAppGroupRoles(userID, clientID string) ([]models.AppGroupRole, error) {
	var assignments []models.AppGroupUsersRoles
	if err := s.db.Preload("AppGroupRole").
		Joins("JOIN app_group_roles ON app_group_users_roles.app_group_role_id = app_group_roles.id").
		Where("app_group_users_roles.user_id = ? AND app_group_roles.client_id = ?", userID, clientID).
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to get user app group roles: %w", err)
	}

	roles := make([]models.AppGroupRole, len(assignments))
	for i, assignment := range assignments {
		roles[i] = assignment.AppGroupRole
	}

	return roles, nil
}

// GetUserAppGroupPermissions 获取用户在应用组的所有权限
func (s *AppGroupService) GetUserAppGroupPermissions(userID, clientID string) ([]models.AppGroupPermission, error) {
	roles, err := s.GetUserAppGroupRoles(userID, clientID)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return []models.AppGroupPermission{}, nil
	}

	roleIDs := make([]string, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}

	var permissions []models.AppGroupPermission
	if err := s.db.Distinct().
		Joins("JOIN app_group_roles_permissions ON app_group_permissions.id = app_group_roles_permissions.app_group_permission_id").
		Where("app_group_roles_permissions.app_group_role_id IN ?", roleIDs).
		Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user app group permissions: %w", err)
	}

	return permissions, nil
}

// GetAppGroupRoleUsers 获取应用组角色的所有用户
func (s *AppGroupService) GetAppGroupRoleUsers(roleID string) ([]models.Users, error) {
	var assignments []models.AppGroupUsersRoles
	if err := s.db.Preload("User").Where("app_group_role_id = ?", roleID).Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to get role users: %w", err)
	}

	users := make([]models.Users, len(assignments))
	for i, assignment := range assignments {
		users[i] = assignment.User
	}

	return users, nil
}

// UserHasAppGroupPermission 检查用户是否拥有应用组权限
func (s *AppGroupService) UserHasAppGroupPermission(userID, clientID string, permissionCodes []string) (bool, error) {
	permissions, err := s.GetUserAppGroupPermissions(userID, clientID)
	if err != nil {
		return false, err
	}

	permMap := make(map[string]bool)
	for _, perm := range permissions {
		permMap[perm.Code] = true
	}

	for _, code := range permissionCodes {
		if !permMap[code] {
			return false, nil
		}
	}

	return true, nil
}

// ======================== 默认角色相关方法 ========================

// GetDefaultAppGroupRole 获取应用组的默认角色
func (s *AppGroupService) GetDefaultAppGroupRole(clientID string) (*models.AppGroupRole, error) {
	var role models.AppGroupRole
	if err := s.db.Where("client_id = ? AND is_default = ?", clientID, true).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get default role: %w", err)
	}
	return &role, nil
}

// SetDefaultAppGroupRole 设置应用组的默认角色
func (s *AppGroupService) SetDefaultAppGroupRole(clientID, roleID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 取消现有默认角色
		if err := tx.Model(&models.AppGroupRole{}).Where("client_id = ? AND is_default = ?", clientID, true).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to clear default role: %w", err)
		}

		// 设置新的默认角色
		if err := tx.Model(&models.AppGroupRole{}).Where("id = ? AND client_id = ?", roleID, clientID).
			Update("is_default", true).Error; err != nil {
			return fmt.Errorf("failed to set default role: %w", err)
		}

		return nil
	})
}

// AssignDefaultRoleToUser 为用户分配默认角色（如果存在）
func (s *AppGroupService) AssignDefaultRoleToUser(userID, clientID, assignedBy string) error {
	defaultRole, err := s.GetDefaultAppGroupRole(clientID)
	if err != nil {
		return err
	}

	if defaultRole == nil {
		return nil // 没有默认角色
	}

	return s.AssignAppGroupRoleToUser(userID, defaultRole.ID, assignedBy)
}
