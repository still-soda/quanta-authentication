package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

// FindByCodes 根据权限代码批量查询权限
func (r *PermissionRepository) FindByCodes(codes []string) ([]*models.Permissions, error) {
	var perms []*models.Permissions
	if err := r.db.Where("code IN ?", codes).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// CountByCodes 统计指定代码的权限数量
func (r *PermissionRepository) CountByCodes(codes []string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Permissions{}).Where("code IN ?", codes).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindAll 获取所有权限
func (r *PermissionRepository) FindAll() ([]models.Permissions, error) {
	var perms []models.Permissions
	if err := r.db.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// ListPermissionsParams 权限列表查询参数
type ListPermissionsParams struct {
	Page      int
	PageSize  int
	Search    string
	Resource  string
	SortField string
	SortOrder string // "asc" 或 "desc"
}

// ListPermissionsResult 权限列表查询结果
type ListPermissionsResult struct {
	Items    []models.Permissions
	Total    int64
	Page     int
	PageSize int
}

// List 获取权限列表（带分页和过滤）
func (r *PermissionRepository) List(params ListPermissionsParams) (*ListPermissionsResult, error) {
	var perms []models.Permissions
	var total int64

	db := r.db.Model(&models.Permissions{})

	// 搜索过滤
	if params.Search != "" {
		search := "%" + params.Search + "%"
		db = db.Where("code LIKE ? OR description LIKE ? OR resource LIKE ?", search, search, search)
	}

	// 资源过滤
	if params.Resource != "" {
		db = db.Where("resource = ?", params.Resource)
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
		params.PageSize = 15
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// 权限表允许排序的字段白名单
	sortableFields := map[string]string{
		"code":       "code",
		"resource":   "resource",
		"action":     "action",
		"created_at": "created_at",
	}

	// 构建排序语句
	orderClause := "created_at DESC" // 默认排序
	if params.SortField != "" {
		// 验证排序字段是否在白名单中，防止 SQL 注入
		if dbField, ok := sortableFields[params.SortField]; ok {
			sortDirection := "ASC"
			if params.SortOrder == "desc" {
				sortDirection = "DESC"
			}
			orderClause = dbField + " " + sortDirection
		}
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := db.Offset(offset).Limit(params.PageSize).Order(orderClause).Find(&perms).Error; err != nil {
		return nil, err
	}

	return &ListPermissionsResult{
		Items:    perms,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// FindByID 根据权限 ID 获取权限信息
func (r *PermissionRepository) FindByID(permID string) (*models.Permissions, error) {
	var perm models.Permissions
	if err := r.db.First(&perm, "id = ?", permID).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

// FindByResource 根据资源获取权限列表
func (r *PermissionRepository) FindByResource(resource string) ([]models.Permissions, error) {
	var perms []models.Permissions
	if err := r.db.Where("resource = ?", resource).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// Create 创建权限
func (r *PermissionRepository) Create(perm *models.Permissions) error {
	return r.db.Create(perm).Error
}

// Update 更新权限
func (r *PermissionRepository) Update(perm *models.Permissions) error {
	return r.db.Save(perm).Error
}

// Delete 删除权限
func (r *PermissionRepository) Delete(permID string) error {
	return r.db.Delete(&models.Permissions{}, "id = ?", permID).Error
}

// DeleteRolePermissions 删除权限关联的所有角色权限关系
func (r *PermissionRepository) DeleteRolePermissions(permID string) error {
	return r.db.Where("permissions_id = ?", permID).Delete(&models.RolesPermissions{}).Error
}

// GetRolePermissions 获取角色拥有的权限
func (r *PermissionRepository) GetRolePermissions(roleID string) ([]models.Permissions, error) {
	var role models.Roles
	if err := r.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
