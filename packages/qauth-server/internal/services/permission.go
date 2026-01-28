package services

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwt"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

func (s *PermissionService) GetPermissionByCodes(code []string) ([]*models.Permissions, error) {
	var permissions []*models.Permissions
	if err := s.db.Where("code IN ?", code).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// CheckCodesExists 验证权限代码是否存在
func (s *PermissionService) CheckCodesExists(codes []string) (bool, error) {
	var count int64
	if err := s.db.Model(&models.Permissions{}).Where("code IN ?", codes).Count(&count).Error; err != nil {
		return false, err
	}
	return count == int64(len(codes)), nil
}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]models.Permissions, error) {
	var permissions []models.Permissions
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
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
	Items    []models.Permissions `json:"items"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

// 权限表允许排序的字段白名单
var permissionSortableFields = map[string]string{
	"code":       "code",
	"resource":   "resource",
	"action":     "action",
	"created_at": "created_at",
}

// ListPermissions 获取权限列表（带分页）
func (s *PermissionService) ListPermissions(params ListPermissionsParams) (*ListPermissionsResult, error) {
	var permissions []models.Permissions
	var total int64

	db := s.db.Model(&models.Permissions{})

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

	// 构建排序语句
	orderClause := "created_at DESC" // 默认排序
	if params.SortField != "" {
		// 验证排序字段是否在白名单中，防止 SQL 注入
		if dbField, ok := permissionSortableFields[params.SortField]; ok {
			sortDirection := "ASC"
			if params.SortOrder == "desc" {
				sortDirection = "DESC"
			}
			orderClause = dbField + " " + sortDirection
		}
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := db.Offset(offset).Limit(params.PageSize).Order(orderClause).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return &ListPermissionsResult{
		Items:    permissions,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetPermissionByID 根据权限 ID 获取权限信息
func (s *PermissionService) GetPermissionByID(permissionID string) (*models.Permissions, error) {
	var permission models.Permissions
	if err := s.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionsByResource 根据资源获取权限列表
func (s *PermissionService) GetPermissionsByResource(resource string) ([]models.Permissions, error) {
	var permissions []models.Permissions
	if err := s.db.Where("resource = ?", resource).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(resource string, action int8, code string, description string) (*models.Permissions, error) {
	permission := &models.Permissions{
		Resource:    resource,
		Action:      action,
		Code:        code,
		Description: description,
	}

	if err := s.db.Create(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(permissionID string, resource string, action int8, code string, description string) (*models.Permissions, error) {
	var permission models.Permissions
	if err := s.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return nil, err
	}

	permission.Resource = resource
	permission.Action = action
	permission.Code = code
	permission.Description = description

	if err := s.db.Save(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(permissionID string) error {
	// 先删除关联的角色权限关系
	if err := s.db.Where("permissions_id = ?", permissionID).Delete(&models.RolesPermissions{}).Error; err != nil {
		return err
	}
	return s.db.Delete(&models.Permissions{}, "id = ?", permissionID).Error
}

// GetPermissionsGroupedByResource 获取按资源分组的权限
func (s *PermissionService) GetPermissionsGroupedByResource() (map[string][]models.Permissions, error) {
	permissions, err := s.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.Permissions)
	for _, perm := range permissions {
		grouped[perm.Resource] = append(grouped[perm.Resource], perm)
	}

	return grouped, nil
}

// GetRolePermissions 获取角色拥有的权限
func (s *PermissionService) GetRolePermissions(roleID string) ([]models.Permissions, error) {
	var role models.Roles
	if err := s.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// VerifyPermissions 验证用户是否拥有指定权限代码
func VerifyPermissions(c *gin.Context, roleService *RoleService, codes []string) error {
	var userInfo *jwt.UserJWTClaims
	if info, exists := c.Get("userInfo"); exists {
		userInfo = info.(*jwt.UserJWTClaims)
	}

	userRole, err := roleService.GetRoleByCode(userInfo.Role)
	if err != nil {
		utilities.GetLogger().Error("Get user role error", "error", err, "userRole", userInfo.Role)
		return app_error.ErrInternalServerError
	}

	hasPermission, err := roleService.RoleHasPermissions(userRole.BaseModelWithUUID.ID, codes)
	if err != nil {
		utilities.GetLogger().Error("Check role permissions error", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return app_error.ErrInternalServerError
	}
	if !hasPermission {
		response.HandlerError(c, app_error.ErrNoPermission)
		return app_error.ErrNoPermission
	}

	return nil
}
