package services

import (
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"

	"gorm.io/gorm"
)

type RoleService struct {
	db                *gorm.DB
	permissionService *PermissionService
	userService       *UserService
}

func NewRoleService(
	db *gorm.DB,
	permissionService *PermissionService,
	userService *UserService,
) *RoleService {
	return &RoleService{db: db, permissionService: permissionService, userService: userService}
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(roleID string) error {
	return s.db.Delete(&models.Roles{}, "id = ?", roleID).Error
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(roleID, name string, code string) (*models.Roles, error) {
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}

	role.Name = name
	role.Code = code

	if err := s.db.Save(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// CreateRole 创建新角色
func (s *RoleService) CreateRole(name, code string, permissionCodes []string) (*models.Roles, error) {
	role := &models.Roles{
		Name: name,
		Code: code,
	}

	if err := s.db.Create(role).Error; err != nil {
		return nil, err
	}

	if len(permissionCodes) > 0 {
		if err := s.GrantPermissionToRole(role.ID, permissionCodes); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// GetAllRoles 获取所有角色列表
func (s *RoleService) GetAllRoles() ([]models.Roles, error) {
	var roles []models.Roles
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID 根据角色 ID 获取角色信息
func (s *RoleService) GetRoleByID(roleID string) (*models.Roles, error) {
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetUserRole 获取用户的角色信息
func (s *RoleService) GetUserRole(userID string) (*models.Roles, error) {
	var userRole models.UsersRoles
	if err := s.db.Preload("Role").First(&userRole, "users_id = ?", userID).Error; err != nil {
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
func (s *RoleService) GrantPermissionToRoleByAction(roleID string, action permissions.Action) error {
	var permissions []models.Permissions
	if err := s.db.Where("action = ?", int8(action)).Find(&permissions).Error; err != nil {
		return err
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	return s.GrantPermissionToRole(roleID, permissionCodes)
}

// RevokePermissionFromRoleByAction 根据操作类型从角色撤销权限
func (s *RoleService) RevokePermissionFromRoleByAction(roleID string, action permissions.Action) error {
	var permissions []models.Permissions
	if err := s.db.Where("action = ?", int8(action)).Find(&permissions).Error; err != nil {
		return err
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	return s.RevokePermissionFromRole(roleID, permissionCodes)
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

// AssignRoleToUserByID 分配角色到用户
func (s *RoleService) AssignRoleToUserByID(userID, roleID string) error {
	user, err := s.userService.GetUserByID(userID, false)
	if err != nil {
		return err
	}

	role, err := s.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	userRole := &models.UsersRoles{
		UsersID: user.ID,
		RolesID: role.ID,
	}

	if err := s.db.Create(userRole).Error; err != nil {
		return err
	}

	return nil
}

// AssignRolesToUserByCode 分配多个角色到用户
func (s *RoleService) AssignRolesToUserByCode(userID string, roleCodes []string) error {
	user, err := s.userService.GetUserByID(userID, false)
	if err != nil {
		return err
	}

	var roles []models.Roles
	if err := s.db.Where("code IN ?", roleCodes).Find(&roles).Error; err != nil {
		return err
	}

	var userRoles []models.UsersRoles
	for _, role := range roles {
		userRoles = append(userRoles, models.UsersRoles{
			UsersID: user.ID,
			RolesID: role.ID,
		})
	}

	if err := s.db.FirstOrCreate(&userRoles).Error; err != nil {
		return err
	}

	return nil
}

// RevokeRoleFromUser 从用户撤销角色
func (s *RoleService) RevokeRoleFromUser(userID, roleID string) error {
	return s.db.Delete(&models.UsersRoles{}, "users_id = ? AND roles_id = ?", userID, roleID).Error
}
