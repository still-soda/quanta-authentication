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
