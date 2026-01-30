package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// Create 创建角色
func (r *RoleRepository) Create(role *models.Roles) error {
	return r.db.Create(role).Error
}

// FindByID 根据ID查询角色
func (r *RoleRepository) FindByID(roleID string) (*models.Roles, error) {
	var role models.Roles
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByIDWithPermissions 根据ID查询角色（预加载权限）
func (r *RoleRepository) FindByIDWithPermissions(roleID string) (*models.Roles, error) {
	var role models.Roles
	if err := r.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode 根据代码查询角色
func (r *RoleRepository) FindByCode(roleCode string) (*models.Roles, error) {
	var role models.Roles
	if err := r.db.First(&role, "code = ?", roleCode).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCodes 根据代码批量查询角色
func (r *RoleRepository) FindByCodes(roleCodes []string) ([]models.Roles, error) {
	var roles []models.Roles
	if err := r.db.Where("code IN ?", roleCodes).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindAll 查询所有角色
func (r *RoleRepository) FindAll() ([]models.Roles, error) {
	var roles []models.Roles
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Update 更新角色
func (r *RoleRepository) Update(role *models.Roles) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(roleID string) error {
	return r.db.Delete(&models.Roles{}, "id = ?", roleID).Error
}

// ListRolesParams 角色列表查询参数
type ListRolesParams struct {
	Page     int
	PageSize int
	Search   string
}

// List 分页查询角色列表
func (r *RoleRepository) List(params ListRolesParams) ([]models.Roles, int64, error) {
	var roles []models.Roles
	var total int64

	db := r.db.Model(&models.Roles{})

	// 搜索过滤
	if params.Search != "" {
		search := "%" + params.Search + "%"
		db = db.Where("name LIKE ? OR code LIKE ?", search, search)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}

	return roles, total, nil
}

// CountUsersByRoleID 统计角色下的用户数量
func (r *RoleRepository) CountUsersByRoleID(roleID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.UsersRoles{}).Where("roles_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountPermissionsByRoleID 统计角色下的权限数量
func (r *RoleRepository) CountPermissionsByRoleID(roleID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.RolesPermissions{}).Where("roles_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// AddPermissions 为角色添加权限
func (r *RoleRepository) AddPermissions(roleID string, permissions []*models.Permissions) error {
	var role models.Roles
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}
	return r.db.Model(&role).Association("Permissions").Append(permissions)
}

// RemovePermissions 从角色移除权限
func (r *RoleRepository) RemovePermissions(roleID string, permissions []*models.Permissions) error {
	var role models.Roles
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}
	return r.db.Model(&role).Association("Permissions").Delete(permissions)
}
