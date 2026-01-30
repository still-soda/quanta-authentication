package services

import (
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/providers"
	"qauth-server/internal/repository"
	"qauth-server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type PermissionService struct {
	repo   *repository.PermissionRepository
	logger providers.ILogger
}

func NewPermissionService(repo *repository.PermissionRepository, logger providers.ILogger) *PermissionService {
	return &PermissionService{
		repo:   repo,
		logger: logger.With("service", "PermissionService"),
	}
}

func (s *PermissionService) GetPermissionByCodes(code []string) ([]*models.Permissions, error) {
	perms, err := s.repo.FindByCodes(code)
	if err != nil {
		return nil, e.ErrFailedToFindPermissionsByCodes.Wrap(err).WithScope("GetPermissionByCodes")
	}
	return perms, nil
}

// CheckCodesExists 验证权限代码是否存在
func (s *PermissionService) CheckCodesExists(codes []string) (bool, error) {
	count, err := s.repo.CountByCodes(codes)
	if err != nil {
		return false, e.ErrFailedToCountPermissionsByCodes.Wrap(err).WithScope("CheckCodesExists")
	}
	return count == int64(len(codes)), nil
}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]models.Permissions, error) {
	perms, err := s.repo.FindAll()
	if err != nil {
		return nil, e.ErrFailedToFindAllPermissions.Wrap(err).WithScope("GetAllPermissions")
	}
	return perms, nil
}

// ListPermissions 获取权限列表（带分页）
func (s *PermissionService) ListPermissions(params repository.ListPermissionsParams) (*repository.ListPermissionsResult, error) {
	result, err := s.repo.List(params)
	if err != nil {
		return nil, e.ErrFailedToListPermissions.Wrap(err).WithScope("ListPermissions")
	}
	return result, nil
}

// GetPermissionByID 根据权限 ID 获取权限信息
func (s *PermissionService) GetPermissionByID(permID string) (*models.Permissions, error) {
	perm, err := s.repo.FindByID(permID)
	if err != nil {
		return nil, e.ErrFailedToFindPermissionByID.Wrap(err).WithScope("GetPermissionByID")
	}
	return perm, nil
}

// GetPermissionsByResource 根据资源获取权限列表
func (s *PermissionService) GetPermissionsByResource(resource string) ([]models.Permissions, error) {
	perms, err := s.repo.FindByResource(resource)
	if err != nil {
		return nil, e.ErrFailedToFindPermissionsByResource.Wrap(err).WithScope("GetPermissionsByResource")
	}
	return perms, nil
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(resource string, action int8, code string, description string) (*models.Permissions, error) {
	perm := &models.Permissions{
		Resource:    resource,
		Action:      action,
		Code:        code,
		Description: description,
	}

	if err := s.repo.Create(perm); err != nil {
		return nil, e.ErrFailedToCreatePermission.Wrap(err).WithScope("CreatePermission")
	}

	return perm, nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(permID string, resource string, action int8, code string, description string) (*models.Permissions, error) {
	perm, err := s.repo.FindByID(permID)
	if err != nil {
		return nil, e.ErrFailedToFindPermissionByID.Wrap(err).WithScope("UpdatePermission")
	}

	perm.Resource = resource
	perm.Action = action
	perm.Code = code
	perm.Description = description

	if err := s.repo.Update(perm); err != nil {
		return nil, e.ErrFailedToUpdatePermission.Wrap(err).WithScope("UpdatePermission")
	}

	return perm, nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(permID string) error {
	// 先删除关联的角色权限关系
	if err := s.repo.DeleteRolePermissions(permID); err != nil {
		return e.ErrFailedToDeleteRolePermissions.Wrap(err).WithScope("DeletePermission")
	}
	if err := s.repo.Delete(permID); err != nil {
		return e.ErrFailedToDeletePermission.Wrap(err).WithScope("DeletePermission")
	}
	return nil
}

// GetPermissionsGroupedByResource 获取按资源分组的权限
func (s *PermissionService) GetPermissionsGroupedByResource() (map[string][]models.Permissions, error) {
	perms, err := s.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.Permissions)
	for _, perm := range perms {
		grouped[perm.Resource] = append(grouped[perm.Resource], perm)
	}

	return grouped, nil
}

// GetRolePermissions 获取角色拥有的权限
func (s *PermissionService) GetRolePermissions(roleID string) ([]models.Permissions, error) {
	perms, err := s.repo.GetRolePermissions(roleID)
	if err != nil {
		return nil, e.ErrFailedToGetRolePermissions.Wrap(err).WithScope("GetRolePermissions")
	}
	return perms, nil
}

// VerifyPermissions 验证用户是否拥有指定权限代码
func VerifyPermissions(c *gin.Context, roleService *RoleService, codes []string) error {
	var userInfo *jwt.UserJWTClaims
	if info, exists := c.Get("userInfo"); exists {
		userInfo = info.(*jwt.UserJWTClaims)
	}

	userRole, err := roleService.GetRoleByCode(userInfo.Role)
	if err != nil {
		return e.ErrInternalServerError.Wrap(err).WithScope("VerifyPermissions")
	}

	hasPermission, err := roleService.RoleHasPermissions(userRole.BaseModelWithUUID.ID, codes)
	if err != nil {
		return e.ErrInternalServerError.Wrap(err).WithScope("VerifyPermissions")
	}
	if !hasPermission {
		return e.ErrNoPermission.Wrap(nil).WithScope("VerifyPermissions")
	}

	return nil
}
