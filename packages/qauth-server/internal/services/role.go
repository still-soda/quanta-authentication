package services

import (
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// 先删除角色的权限关联
	if err := s.db.Where("roles_id = ?", roleID).Delete(&models.RolesPermissions{}).Error; err != nil {
		return err
	}
	// 删除用户角色关联
	if err := s.db.Where("roles_id = ?", roleID).Delete(&models.UsersRoles{}).Error; err != nil {
		return err
	}
	return s.db.Delete(&models.Roles{}, "id = ?", roleID).Error
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(roleID, name string, code string, description string) (*models.Roles, error) {
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}

	role.Name = name
	role.Code = code
	role.Description = description

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

// RoleWithStats 角色带统计信息
type RoleWithStats struct {
	models.Roles
	UserCount       int64 `json:"user_count"`
	PermissionCount int64 `json:"permission_count"`
}

// GetAllRolesWithStats 获取所有角色列表（带统计信息）
func (s *RoleService) GetAllRolesWithStats() ([]RoleWithStats, error) {
	var roles []models.Roles
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	result := make([]RoleWithStats, len(roles))
	for i, role := range roles {
		var userCount int64
		s.db.Model(&models.UsersRoles{}).Where("roles_id = ?", role.ID).Count(&userCount)

		var permCount int64
		s.db.Model(&models.RolesPermissions{}).Where("roles_id = ?", role.ID).Count(&permCount)

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
	var roles []models.Roles
	var total int64

	db := s.db.Model(&models.Roles{})

	// 搜索过滤
	if params.Search != "" {
		search := "%" + params.Search + "%"
		db = db.Where("name LIKE ? OR code LIKE ?", search, search)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := db.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, err
	}

	// 添加统计信息
	result := make([]RoleWithStats, len(roles))
	for i, role := range roles {
		var userCount int64
		s.db.Model(&models.UsersRoles{}).Where("roles_id = ?", role.ID).Count(&userCount)

		var permCount int64
		s.db.Model(&models.RolesPermissions{}).Where("roles_id = ?", role.ID).Count(&permCount)

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
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}

	var userCount int64
	s.db.Model(&models.UsersRoles{}).Where("roles_id = ?", role.ID).Count(&userCount)

	var permCount int64
	s.db.Model(&models.RolesPermissions{}).Where("roles_id = ?", role.ID).Count(&permCount)

	return &RoleWithStats{
		Roles:           role,
		UserCount:       userCount,
		PermissionCount: permCount,
	}, nil
}

// GetRoleByID 根据角色 ID 获取角色信息
func (s *RoleService) GetRoleByID(roleID string) (*models.Roles, error) {
	var role models.Roles
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetUserRole 获取用户的系统角色信息
func (s *RoleService) GetUserRole(userID string) (*models.Roles, error) {
	var userRole models.UsersRoles
	if err := s.db.Preload("Role", "is_system = ?", true).
		First(&userRole, "users_id = ?", userID).
		Error; err != nil {
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

	if err := s.db.Create(&userRoles).Clauses(clause.OnConflict{DoNothing: true}).Error; err != nil {
		return err
	}

	return nil
}

// RevokeRoleFromUser 从用户撤销角色
func (s *RoleService) RevokeRoleFromUser(userID, roleID string) error {
	return s.db.Delete(&models.UsersRoles{}, "users_id = ? AND roles_id = ?", userID, roleID).Error
}
