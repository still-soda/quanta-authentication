package services

import (
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"

	"gorm.io/gorm"
)

type RoleService struct {
	db                *gorm.DB
	permissionService *PermissionService
}

func NewRoleService(db *gorm.DB, permissionService *PermissionService) *RoleService {
	return &RoleService{db: db, permissionService: permissionService}
}

// GetRoleByID 根据角色 ID 获取角色信息
func (s *RoleService) GetRoleByID(roleID string) (*RoleService, error) {
	var role RoleService
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetUserRole 获取用户的角色信息
func (s *RoleService) GetUserRole(userID string) (*models.Roles, error) {
	var userRole models.UsersRoles
	if err := s.db.Preload("Role").First(&userRole, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &userRole.Role, nil
}

// GetRoleByCode 根据角色代码获取角色信息
func (s *RoleService) GetRoleByCode(roleCode string) (*models.Roles, error) {
	var role models.Roles
	if err := s.db.First(&role, "code = ?", roleCode).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GrantPermissionToRole 分配权限到角色
func (s *RoleService) GrantPermissionToRole(roleID string, permissionCodes []string) error {
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	permissions, err := s.permissionService.GetPermissionByCodes(permissionCodes)
	if err != nil {
		return err
	}

	if err := s.db.Model(&role).Association("Permissions").Append(permissions); err != nil {
		return err
	}

	return nil
}

// RevokePermissionFromRole 从角色撤销权限
func (s *RoleService) RevokePermissionFromRole(roleID string, permissionCodes []string) error {
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	permissions, err := s.permissionService.GetPermissionByCodes(permissionCodes)
	if err != nil {
		return err
	}

	if err := s.db.Model(&role).Association("Permissions").Delete(permissions); err != nil {
		return err
	}

	return nil
}

// GrantPermissionToRoleByAction 根据操作类型分配权限到角色
func (s *RoleService) GrantPermissionToRoleByAction(roleId string, action permissions.Action) error {
	var permissions []models.Permissions
	if err := s.db.Where("action = ?", int8(action)).Find(&permissions).Error; err != nil {
		return err
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	return s.GrantPermissionToRole(roleId, permissionCodes)
}

// RevokePermissionFromRoleByAction 根据操作类型从角色撤销权限
func (s *RoleService) RevokePermissionFromRoleByAction(roleId string, action permissions.Action) error {
	var permissions []models.Permissions
	if err := s.db.Where("action = ?", int8(action)).Find(&permissions).Error; err != nil {
		return err
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	return s.RevokePermissionFromRole(roleId, permissionCodes)
}

// RoleHasPermissions 检查角色是否拥有指定权限代码
func (s *RoleService) RoleHasPermissions(roleID string, permissionCodes []string) (bool, error) {
	var role models.Roles
	if err := s.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return false, err
	}

	permissionMap := make(map[string]bool)
	for _, perm := range role.Permissions {
		permissionMap[perm.Code] = true
	}

	for _, code := range permissionCodes {
		if !permissionMap[code] {
			return false, nil
		}
	}

	return true, nil
}
