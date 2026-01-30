package services

import (
	"qauth-server/internal/config/permissions"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/repository"
)

type RoleService struct {
	roleRepo          *repository.RoleRepository
	usersRolesRepo    *repository.UsersRolesRepository
	rolesPermsRepo    *repository.RolesPermissionsRepository
	permissionService *PermissionService
	userService       *UserService
}

func NewRoleService(
	roleRepo *repository.RoleRepository,
	usersRolesRepo *repository.UsersRolesRepository,
	rolesPermsRepo *repository.RolesPermissionsRepository,
	permissionService *PermissionService,
	userService *UserService,
) *RoleService {
	return &RoleService{
		roleRepo:          roleRepo,
		usersRolesRepo:    usersRolesRepo,
		rolesPermsRepo:    rolesPermsRepo,
		permissionService: permissionService,
		userService:       userService,
	}
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(roleID string) error {
	// 先删除角色的权限关联
	if err := s.rolesPermsRepo.DeleteByRoleID(roleID); err != nil {
		return e.ErrFailedToDeleteRolePerms.WithScope("RoleService.DeleteRole").Wrap(err)
	}
	// 删除用户角色关联
	if err := s.usersRolesRepo.DeleteByRoleID(roleID); err != nil {
		return e.ErrFailedToDeleteUserRoles.WithScope("RoleService.DeleteRole").Wrap(err)
	}
	if err := s.roleRepo.Delete(roleID); err != nil {
		return e.ErrFailedToDeleteRole.WithScope("RoleService.DeleteRole").Wrap(err)
	}
	return nil
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(roleID, name string, code string, description string) (*models.Roles, error) {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, e.ErrRoleNotFound.WithScope("RoleService.UpdateRole").Wrap(err)
	}

	role.Name = name
	role.Code = code
	role.Description = description

	if err := s.roleRepo.Update(role); err != nil {
		return nil, e.ErrFailedToUpdateRole.WithScope("RoleService.UpdateRole").Wrap(err)
	}

	return role, nil
}

// CreateRole 创建新角色
func (s *RoleService) CreateRole(name, code string, permissionCodes []string) (*models.Roles, error) {
	role := &models.Roles{
		Name: name,
		Code: code,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, e.ErrFailedToCreateRole.WithScope("RoleService.CreateRole").Wrap(err)
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
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, e.ErrFailedToListRoles.WithScope("RoleService.GetAllRoles").Wrap(err)
	}
	return roles, nil
}

// RoleWithStats 角色带统计信息
type RoleWithStats struct {
	models.Roles
	UserCount       int64 `json:"user_count"`
	PermissionCount int64 `json:"permission_count"`
}

// GetAllRolesWithStats 获取所有角色列表（带统计信息）
func (s *RoleService) GetAllRolesWithStats() ([]RoleWithStats, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, e.ErrFailedToListRoles.WithScope("RoleService.GetAllRolesWithStats").Wrap(err)
	}

	result := make([]RoleWithStats, len(roles))
	for i, role := range roles {
		userCount, err := s.roleRepo.CountUsersByRoleID(role.ID)
		if err != nil {
			return nil, e.ErrFailedToCountUsers.WithScope("RoleService.GetAllRolesWithStats").Wrap(err)
		}

		permCount, err := s.roleRepo.CountPermissionsByRoleID(role.ID)
		if err != nil {
			return nil, e.ErrFailedToCountPermissions.WithScope("RoleService.GetAllRolesWithStats").Wrap(err)
		}

		result[i] = RoleWithStats{
			Roles:           role,
			UserCount:       userCount,
			PermissionCount: permCount,
		}
	}

	return result, nil
}

// ListRolesParams 角色列表查询参数
type ListRolesParams struct {
	Page     int
	PageSize int
	Search   string
}

// ListRolesResult 角色列表查询结果
type ListRolesResult struct {
	Items    []RoleWithStats `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// ListRolesWithStats 获取角色列表（带分页和统计信息）
func (s *RoleService) ListRolesWithStats(params ListRolesParams) (*ListRolesResult, error) {
	roles, total, err := s.roleRepo.List(repository.ListRolesParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Search:   params.Search,
	})
	if err != nil {
		return nil, e.ErrFailedToListRoles.WithScope("RoleService.ListRolesWithStats").Wrap(err)
	}

	// 添加统计信息
	result := make([]RoleWithStats, len(roles))
	for i, role := range roles {
		userCount, err := s.roleRepo.CountUsersByRoleID(role.ID)
		if err != nil {
			return nil, e.ErrFailedToCountUsers.WithScope("RoleService.ListRolesWithStats").Wrap(err)
		}

		permCount, err := s.roleRepo.CountPermissionsByRoleID(role.ID)
		if err != nil {
			return nil, e.ErrFailedToCountPermissions.WithScope("RoleService.ListRolesWithStats").Wrap(err)
		}

		result[i] = RoleWithStats{
			Roles:           role,
			UserCount:       userCount,
			PermissionCount: permCount,
		}
	}

	return &ListRolesResult{
		Items:    result,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetRoleWithStats 获取单个角色（带统计信息）
func (s *RoleService) GetRoleWithStats(roleID string) (*RoleWithStats, error) {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, e.ErrRoleNotFound.WithScope("RoleService.GetRoleWithStats").Wrap(err)
	}

	userCount, err := s.roleRepo.CountUsersByRoleID(role.ID)
	if err != nil {
		return nil, e.ErrFailedToCountUsers.WithScope("RoleService.GetRoleWithStats").Wrap(err)
	}

	permCount, err := s.roleRepo.CountPermissionsByRoleID(role.ID)
	if err != nil {
		return nil, e.ErrFailedToCountPermissions.WithScope("RoleService.GetRoleWithStats").Wrap(err)
	}

	return &RoleWithStats{
		Roles:           *role,
		UserCount:       userCount,
		PermissionCount: permCount,
	}, nil
}

// GetRoleByID 根据角色 ID 获取角色信息
func (s *RoleService) GetRoleByID(roleID string) (*models.Roles, error) {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, e.ErrRoleNotFound.WithScope("RoleService.GetRoleByID").Wrap(err)
	}
	return role, nil
}

// GetUserRole 获取用户的系统角色信息
func (s *RoleService) GetUserRole(userID string) (*models.Roles, error) {
	userRole, err := s.usersRolesRepo.FindByUserID(userID)
	if err != nil {
		return nil, e.ErrFailedToGetUserRole.WithScope("RoleService.GetUserRole").Wrap(err)
	}
	return &userRole.Role, nil
}

// GetRoleByCode 根据角色代码获取角色信息
func (s *RoleService) GetRoleByCode(roleCode string) (*models.Roles, error) {
	role, err := s.roleRepo.FindByCode(roleCode)
	if err != nil {
		return nil, e.ErrRoleNotFound.WithScope("RoleService.GetRoleByCode").Wrap(err)
	}
	return role, nil
}

// GrantPermissionToRole 分配权限到角色
func (s *RoleService) GrantPermissionToRole(roleID string, permissionCodes []string) error {
	permissions, err := s.permissionService.GetPermissionByCodes(permissionCodes)
	if err != nil {
		return e.ErrPermissionNoExist.WithScope("RoleService.GrantPermissionToRole").Wrap(err)
	}

	if err := s.roleRepo.AddPermissions(roleID, permissions); err != nil {
		return e.ErrFailedToGrantPermissions.WithScope("RoleService.GrantPermissionToRole").Wrap(err)
	}
	return nil
}

// RevokePermissionFromRole 从角色撤销权限
func (s *RoleService) RevokePermissionFromRole(roleID string, permissionCodes []string) error {
	permissions, err := s.permissionService.GetPermissionByCodes(permissionCodes)
	if err != nil {
		return e.ErrPermissionNoExist.WithScope("RoleService.RevokePermissionFromRole").Wrap(err)
	}

	if err := s.roleRepo.RemovePermissions(roleID, permissions); err != nil {
		return e.ErrFailedToRevokePermissions.WithScope("RoleService.RevokePermissionFromRole").Wrap(err)
	}
	return nil
}

// GrantPermissionToRoleByAction 根据操作类型分配权限到角色
func (s *RoleService) GrantPermissionToRoleByAction(roleID string, action permissions.Action) error {
	permissions, err := s.rolesPermsRepo.FindPermissionsByAction(int8(action))
	if err != nil {
		return e.ErrPermissionNoExist.WithScope("RoleService.GrantPermissionToRoleByAction").Wrap(err)
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	return s.GrantPermissionToRole(roleID, permissionCodes)
}

// RevokePermissionFromRoleByAction 根据操作类型从角色撤销权限
func (s *RoleService) RevokePermissionFromRoleByAction(roleID string, action permissions.Action) error {
	permissions, err := s.rolesPermsRepo.FindPermissionsByAction(int8(action))
	if err != nil {
		return e.ErrPermissionNoExist.WithScope("RoleService.RevokePermissionFromRoleByAction").Wrap(err)
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	return s.RevokePermissionFromRole(roleID, permissionCodes)
}

// RoleHasPermissions 检查角色是否拥有指定权限代码
func (s *RoleService) RoleHasPermissions(roleID string, permissionCodes []string) (bool, error) {
	role, err := s.roleRepo.FindByIDWithPermissions(roleID)
	if err != nil {
		return false, e.ErrFailedToCheckPermissions.WithScope("RoleService.RoleHasPermissions").Wrap(err)
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

	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return e.ErrRoleNotFound.WithScope("RoleService.AssignRoleToUserByID").Wrap(err)
	}

	userRole := &models.UsersRoles{
		UsersID: user.ID,
		RolesID: role.ID,
	}

	if err := s.usersRolesRepo.Create(userRole); err != nil {
		return e.ErrFailedToAssignRole.WithScope("RoleService.AssignRoleToUserByID").Wrap(err)
	}
	return nil
}

// AssignRolesToUserByCode 分配多个角色到用户
func (s *RoleService) AssignRolesToUserByCode(userID string, roleCodes []string) error {
	user, err := s.userService.GetUserByID(userID, false)
	if err != nil {
		return err
	}

	roles, err := s.roleRepo.FindByCodes(roleCodes)
	if err != nil {
		return e.ErrRoleNotFound.WithScope("RoleService.AssignRolesToUserByCode").Wrap(err)
	}

	var userRoles []models.UsersRoles
	for _, role := range roles {
		userRoles = append(userRoles, models.UsersRoles{
			UsersID: user.ID,
			RolesID: role.ID,
		})
	}

	if err := s.usersRolesRepo.BatchCreate(userRoles); err != nil {
		return e.ErrFailedToAssignRole.WithScope("RoleService.AssignRolesToUserByCode").Wrap(err)
	}
	return nil
}

// RevokeRoleFromUser 从用户撤销角色
func (s *RoleService) RevokeRoleFromUser(userID, roleID string) error {
	if err := s.usersRolesRepo.DeleteByUserIDAndRoleID(userID, roleID); err != nil {
		return e.ErrFailedToRevokeRole.WithScope("RoleService.RevokeRoleFromUser").Wrap(err)
	}
	return nil
}
