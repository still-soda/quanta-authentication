package services

import (
	"errors"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/providers"
	"qauth-server/internal/repository"

	"gorm.io/gorm"
)

// AppGroupService 应用组权限服务
type AppGroupService struct {
	db                   *gorm.DB
	userService          *UserService
	appGroupAdminRepo    *repository.AppGroupAdminRepository
	appGroupPermRepo     *repository.AppGroupPermissionRepository
	appGroupRoleRepo     *repository.AppGroupRoleRepository
	appGroupUserRoleRepo *repository.AppGroupUserRoleRepository
	logger               providers.ILogger
}

// NewAppGroupService 创建应用组权限服务
func NewAppGroupService(
	db *gorm.DB,
	userSrv *UserService,
	appGroupAdminRepo *repository.AppGroupAdminRepository,
	appGroupPermRepo *repository.AppGroupPermissionRepository,
	appGroupRoleRepo *repository.AppGroupRoleRepository,
	appGroupUserRoleRepo *repository.AppGroupUserRoleRepository,
	logger providers.ILogger,
) *AppGroupService {
	return &AppGroupService{
		db:                   db,
		userService:          userSrv,
		appGroupAdminRepo:    appGroupAdminRepo,
		appGroupPermRepo:     appGroupPermRepo,
		appGroupRoleRepo:     appGroupRoleRepo,
		appGroupUserRoleRepo: appGroupUserRoleRepo,
		logger:               logger.With("service", "AppGroupService"),
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

	if err := s.appGroupAdminRepo.Create(ownerAdmin); err != nil {
		return e.ErrFailedToCreateOwnerAdmin.Wrap(err)
	}

	s.logger.Info("initialized app group admins", "clientID", clientID, "creatorID", creatorID)
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

	if err := s.appGroupAdminRepo.Create(admin); err != nil {
		return e.ErrFailedToAddAppGroupAdmin.Wrap(err)
	}

	return nil
}

// RemoveAppGroupAdmin 移除应用组管理员
func (s *AppGroupService) RemoveAppGroupAdmin(clientID, userID string, adminType models.AppGroupAdminType) error {
	if err := s.appGroupAdminRepo.Delete(clientID, userID, adminType); err != nil {
		return e.ErrFailedToRemoveAppGroupAdmin.Wrap(err)
	}
	return nil
}

// GetAppGroupAdmins 获取应用组的所有管理员
func (s *AppGroupService) GetAppGroupAdmins(clientID string) ([]models.AppGroupAdmin, error) {
	admins, err := s.appGroupAdminRepo.FindByClientID(clientID)
	if err != nil {
		return nil, e.ErrFailedToGetAppGroupAdmins.Wrap(err)
	}
	return admins, nil
}

// GetAppGroupAdminsByUser 获取用户在所有应用组的管理员身份
func (s *AppGroupService) GetAppGroupAdminsByUser(userID string) ([]models.AppGroupAdmin, error) {
	admins, err := s.appGroupAdminRepo.FindByUserID(userID)
	if err != nil {
		return nil, e.ErrFailedToGetUserAppGroupAdmins.Wrap(err)
	}
	return admins, nil
}

// IsAppGroupAdmin 检查用户是否为应用组管理员
func (s *AppGroupService) IsAppGroupAdmin(clientID, userID string) (bool, models.AppGroupAdminType, error) {
	admin, err := s.appGroupAdminRepo.FindByClientIDAndUserID(clientID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", e.ErrFailedToCheckAppGroupAdmin.Wrap(err)
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

	if err := s.appGroupPermRepo.Create(permission); err != nil {
		return nil, e.ErrFailedToCreateAppGroupPermission.Wrap(err)
	}

	return permission, nil
}

// UpdateAppGroupPermission 更新应用组权限
func (s *AppGroupService) UpdateAppGroupPermission(permissionID, name, description string) (*models.AppGroupPermission, error) {
	permission, err := s.appGroupPermRepo.FindByID(permissionID)
	if err != nil {
		return nil, e.ErrFailedToFindPermission.Wrap(err)
	}

	permission.Name = name
	permission.Description = description

	if err := s.appGroupPermRepo.Update(permission); err != nil {
		return nil, e.ErrFailedToUpdateAppGroupPermission.Wrap(err)
	}

	return permission, nil
}

// DeleteAppGroupPermission 删除应用组权限
func (s *AppGroupService) DeleteAppGroupPermission(permissionID string) error {
	// 先删除角色权限关联
	if err := s.appGroupRoleRepo.DeletePermissionRoleAssociations(permissionID); err != nil {
		return e.ErrFailedToDeleteRolePermissionAssociations.Wrap(err)
	}

	// 删除权限
	if err := s.appGroupPermRepo.Delete(permissionID); err != nil {
		return e.ErrFailedToDeleteAppGroupPermission.Wrap(err)
	}

	return nil
}

// GetAppGroupPermission 获取应用组权限详情
func (s *AppGroupService) GetAppGroupPermission(permissionID string) (*models.AppGroupPermission, error) {
	permission, err := s.appGroupPermRepo.FindByID(permissionID)
	if err != nil {
		return nil, e.ErrFailedToGetAppGroupPermission.Wrap(err)
	}
	return permission, nil
}

// GetAppGroupPermissions 获取应用组的所有权限
func (s *AppGroupService) GetAppGroupPermissions(clientID string) ([]models.AppGroupPermission, error) {
	perms, err := s.appGroupPermRepo.FindByClientID(clientID)
	if err != nil {
		return nil, e.ErrFailedToGetAppGroupPermissions.Wrap(err)
	}
	return perms, nil
}

// GetAppGroupPermissionsByResource 获取应用组按资源分组的权限
func (s *AppGroupService) GetAppGroupPermissionsByResource(clientID string) (map[string][]models.AppGroupPermission, error) {
	perms, err := s.GetAppGroupPermissions(clientID)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.AppGroupPermission)
	for _, perm := range perms {
		grouped[perm.Resource] = append(grouped[perm.Resource], perm)
	}

	return grouped, nil
}

// GetAppGroupPermissionByCodes 根据权限代码获取权限列表
func (s *AppGroupService) GetAppGroupPermissionByCodes(clientID string, codes []string) ([]models.AppGroupPermission, error) {
	perms, err := s.appGroupPermRepo.FindByCodes(clientID, codes)
	if err != nil {
		return nil, e.ErrFailedToGetPermissionByCodes.Wrap(err)
	}
	return perms, nil
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

	if err := s.appGroupRoleRepo.Create(role); err != nil {
		return nil, e.ErrFailedToCreateAppGroupRole.Wrap(err)
	}

	return role, nil
}

// UpdateAppGroupRole 更新应用组角色
func (s *AppGroupService) UpdateAppGroupRole(roleID, name, description string, isDefault bool) (*models.AppGroupRole, error) {
	role, err := s.appGroupRoleRepo.FindByID(roleID)
	if err != nil {
		return nil, e.ErrFailedToFindRole.Wrap(err)
	}

	role.Name = name
	role.Description = description
	role.IsDefault = isDefault

	if err := s.appGroupRoleRepo.Update(role); err != nil {
		return nil, e.ErrFailedToUpdateAppGroupRole.Wrap(err)
	}

	return role, nil
}

// DeleteAppGroupRole 删除应用组角色
func (s *AppGroupService) DeleteAppGroupRole(roleID string) error {
	// 先删除角色权限关联
	if err := s.appGroupRoleRepo.DeleteRolePermissionAssociations(roleID); err != nil {
		return e.ErrFailedToDeleteRolePermissionAssociations.Wrap(err)
	}

	// 删除用户角色关联
	if err := s.appGroupRoleRepo.DeleteUserRoleAssociations(roleID); err != nil {
		return e.ErrFailedToDeleteUserRoleAssociations.Wrap(err)
	}

	// 删除角色
	if err := s.appGroupRoleRepo.Delete(roleID); err != nil {
		return e.ErrFailedToDeleteAppGroupRole.Wrap(err)
	}

	return nil
}

// GetAppGroupRole 获取应用组角色详情
func (s *AppGroupService) GetAppGroupRole(roleID string) (*models.AppGroupRole, error) {
	role, err := s.appGroupRoleRepo.FindByID(roleID)
	if err != nil {
		return nil, e.ErrFailedToGetAppGroupRole.Wrap(err)
	}
	return role, nil
}

// GetAppGroupRoles 获取应用组的所有角色
func (s *AppGroupService) GetAppGroupRoles(clientID string) ([]models.AppGroupRole, error) {
	roles, err := s.appGroupRoleRepo.FindByClientID(clientID)
	if err != nil {
		return nil, e.ErrFailedToGetAppGroupRoles.Wrap(err)
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
		userCount, _ := s.appGroupRoleRepo.CountUsers(role.ID)
		permCount, _ := s.appGroupRoleRepo.CountPermissions(role.ID)

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

	userCount, _ := s.appGroupRoleRepo.CountUsers(role.ID)
	permCount, _ := s.appGroupRoleRepo.CountPermissions(role.ID)

	return &AppGroupRoleWithStats{
		AppGroupRole:    *role,
		UserCount:       userCount,
		PermissionCount: permCount,
	}, nil
}

// GrantPermissionsToAppGroupRole 为应用组角色分配权限
func (s *AppGroupService) GrantPermissionsToAppGroupRole(roleID string, permissionIDs []string) error {
	role, err := s.appGroupRoleRepo.FindByID(roleID)
	if err != nil {
		return e.ErrFailedToFindRole.Wrap(err)
	}

	perms, err := s.appGroupPermRepo.FindByIDs(permissionIDs)
	if err != nil {
		return e.ErrFailedToFindPermissions.Wrap(err)
	}

	// 验证权限属于同一客户端
	for _, perm := range perms {
		if perm.ClientID != role.ClientID {
			return e.ErrPermissionNotBelongToClient.WithDetails("permissionID", perm.ID, "clientID", role.ClientID)
		}
	}

	// 提取权限 ID
	permIDs := make([]string, len(perms))
	for i, perm := range perms {
		permIDs[i] = perm.ID
	}

	if err := s.appGroupRoleRepo.GrantPermissions(roleID, permIDs); err != nil {
		return e.ErrFailedToGrantPermissionToRole.Wrap(err)
	}

	return nil
}

// RevokePermissionsFromAppGroupRole 从应用组角色撤销权限
func (s *AppGroupService) RevokePermissionsFromAppGroupRole(roleID string, permissionIDs []string) error {
	if err := s.appGroupRoleRepo.RevokePermissions(roleID, permissionIDs); err != nil {
		return e.ErrFailedToRevokePermissionsFromRole.Wrap(err)
	}
	return nil
}

// SetAppGroupRolePermissions 设置应用组角色的权限（替换现有权限）
func (s *AppGroupService) SetAppGroupRolePermissions(roleID string, permissionIDs []string) error {
	_, err := s.appGroupRoleRepo.FindByID(roleID)
	if err != nil {
		return e.ErrFailedToFindRole.Wrap(err)
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		repoWithTx := s.appGroupRoleRepo.WithTx(tx)

		// 删除现有关联
		if err := repoWithTx.ClearPermissions(roleID); err != nil {
			return e.ErrFailedToClearExistingPermissions.Wrap(err)
		}

		// 添加新关联
		if err := repoWithTx.GrantPermissions(roleID, permissionIDs); err != nil {
			return e.ErrFailedToGrantPermissionToRole.Wrap(err)
		}

		return nil
	})
}

// GetAppGroupRolePermissions 获取应用组角色的权限
func (s *AppGroupService) GetAppGroupRolePermissions(roleID string) ([]models.AppGroupPermission, error) {
	role, err := s.appGroupRoleRepo.FindWithPermissions(roleID)
	if err != nil {
		return nil, e.ErrFailedToGetRolePermissions.Wrap(err)
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

	if err := s.appGroupUserRoleRepo.Assign(assignment); err != nil {
		return e.ErrFailedToAssignRoleToUser.Wrap(err)
	}

	return nil
}

// RevokeAppGroupRoleFromUser 从用户撤销应用组角色
func (s *AppGroupService) RevokeAppGroupRoleFromUser(userID, roleID string) error {
	if err := s.appGroupUserRoleRepo.Revoke(userID, roleID); err != nil {
		return e.ErrFailedToRevokeRoleFromUser.Wrap(err)
	}
	return nil
}

// GetUserAppGroupRoles 获取用户在应用组的所有角色
func (s *AppGroupService) GetUserAppGroupRoles(userID, clientID string) ([]models.AppGroupRole, error) {
	assignments, err := s.appGroupUserRoleRepo.FindUserRoles(userID, clientID)
	if err != nil {
		return nil, e.ErrFailedToGetUserAppGroupRoles.Wrap(err)
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

	permissions, err := s.appGroupPermRepo.FindByRoleIDs(roleIDs)
	if err != nil {
		return nil, e.ErrFailedToGetUserAppGroupPermissions.Wrap(err)
	}

	return permissions, nil
}

// GetAppGroupRoleUsers 获取应用组角色的所有用户
func (s *AppGroupService) GetAppGroupRoleUsers(roleID string) ([]models.Users, error) {
	assignments, err := s.appGroupUserRoleRepo.FindRoleUsers(roleID)
	if err != nil {
		return nil, e.ErrFailedToGetRoleUsers.Wrap(err)
	}

	users := make([]models.Users, len(assignments))
	for i, assignment := range assignments {
		users[i] = assignment.User
	}

	return users, nil
}

// UserHasAppGroupPermission 检查用户是否拥有应用组权限
func (s *AppGroupService) UserHasAppGroupPermission(userID, clientID string, permissionCodes []string) (bool, error) {
	perms, err := s.GetUserAppGroupPermissions(userID, clientID)
	if err != nil {
		return false, err
	}

	permMap := make(map[string]bool)
	for _, perm := range perms {
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
	role, err := s.appGroupRoleRepo.FindDefaultRole(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, e.ErrFailedToGetDefaultRole.Wrap(err)
	}
	return role, nil
}

// SetDefaultAppGroupRole 设置应用组的默认角色
func (s *AppGroupService) SetDefaultAppGroupRole(clientID, roleID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		repoWithTx := s.appGroupRoleRepo.WithTx(tx)

		// 取消现有默认角色
		if err := repoWithTx.ClearDefaultRole(clientID); err != nil {
			return e.ErrFailedToClearDefaultRole.Wrap(err)
		}

		// 设置新的默认角色
		if err := repoWithTx.SetDefaultRole(roleID, clientID); err != nil {
			return e.ErrFailedToSetDefaultRole.Wrap(err)
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
