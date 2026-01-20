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

// ValidatePermissionCodes 验证用户是否拥有指定权限代码
func ValidatePermissionCodes(c *gin.Context, roleService *RoleService, codes []string) error {
	var userInfo *jwt.UserJWTClaims
	if info, exists := c.Get("userInfo"); exists {
		userInfo = info.(*jwt.UserJWTClaims)
	}

	userRole, err := roleService.GetRoleByCode(userInfo.Role)
	if err != nil {
		utilities.GetLogger().Error("Get user role error", "error", err, "userRole", userInfo.Role)
		return app_error.ErrInternalServerError
	}

	hasPermission, err := roleService.RoleHasPermissions(userRole.ID, codes)
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
