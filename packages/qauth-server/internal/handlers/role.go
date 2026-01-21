package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService       *services.RoleService
	permissionService *services.PermissionService
}

func NewRoleHandler(roleService *services.RoleService, permissionService *services.PermissionService) *RoleHandler {
	return &RoleHandler{roleService: roleService, permissionService: permissionService}
}

// getRoleExistByID 检查角色是否存在
func (h *RoleHandler) getRoleExistByID(roleID string) (bool, error) {
	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// GetRoles 处理获取角色列表请求
func (h *RoleHandler) GetRoles(c *gin.Context) {
	// 验证是否有查看角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 获取所有角色
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetRoles)
		return
	}

	response.HandlerSuccess(c, roles)
}

// GetRole 处理获取单个角色请求
func (h *RoleHandler) GetRole(c *gin.Context) {
	// 验证是否有查看角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	roleID := c.Param("id")

	// 获取角色
	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetRoles)
		return
	}
	if role == nil {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	response.HandlerSuccess(c, role)
}

// CreateRole 处理创建角色请求
func (h *RoleHandler) CreateRole(c *gin.Context) {
	// 验证是否有创建角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleCreate}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 绑定请求参数
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Code        string   `json:"code" binding:"required"`
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证权限是否存在
	if exist, err := h.permissionService.CheckCodesExists(req.Permissions); err != nil || !exist {
		response.HandlerError(c, app_error.ErrPermissionNoExist)
		return
	}

	// 创建角色
	role, err := h.roleService.CreateRole(req.Name, req.Code, req.Permissions)
	if err != nil {
		response.HandlerError(c, err)
		return
	}

	response.HandlerSuccess(c, role)
}

// UpdateRole 处理更新角色请求
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	// 验证是否有更新权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleUpdate}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 绑定请求参数
	var req struct {
		RoleID string `json:"role_id" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Code   string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证角色是否存在
	if exist, err := h.getRoleExistByID(req.RoleID); err != nil {
		response.HandlerError(c, err)
		return
	} else if !exist {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 更新角色
	role, err := h.roleService.UpdateRole(req.RoleID, req.Name, req.Code)
	if err != nil {
		response.HandlerError(c, err)
		return
	}

	response.HandlerSuccess(c, role)
}

// DeleteRole 处理删除角色请求
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	// 验证是否有删除权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleDelete}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 绑定请求参数
	var req struct {
		RoleID string `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证角色是否存在
	if exist, err := h.getRoleExistByID(req.RoleID); err != nil {
		response.HandlerError(c, err)
		return
	} else if !exist {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 删除角色
	if err := h.roleService.DeleteRole(req.RoleID); err != nil {
		response.HandlerError(c, err)
		return
	}
	response.HandlerSuccess(c, gin.H{"message": "role deleted successfully"})
}

// GrantPermissionsToRole 处理为角色授予权限请求
func (h *RoleHandler) GrantPermissionsToRole(c *gin.Context) {
	// 验证是否有授权权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleAssignPermissions}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 绑定请求参数
	var req struct {
		RoleID      string   `json:"role_id" binding:"required"`
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证角色是否存在
	if exist, err := h.getRoleExistByID(req.RoleID); err != nil {
		response.HandlerError(c, err)
		return
	} else if !exist {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 验证权限是否存在
	if exist, err := h.permissionService.CheckCodesExists(req.Permissions); err != nil || !exist {
		response.HandlerError(c, app_error.ErrPermissionNoExist)
		return
	}

	// 为角色授予权限
	if err := h.roleService.GrantPermissionToRole(req.RoleID, req.Permissions); err != nil {
		response.HandlerError(c, err)
		return
	}

	response.HandlerSuccess(c, gin.H{"message": "permissions granted to role successfully"})
}

// RevokePermissionsFromRole 处理从角色撤销权限请求
func (h *RoleHandler) RevokePermissionsFromRole(c *gin.Context) {
	// 验证是否有撤销权限的权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleRevokePermissions}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 绑定请求参数
	var req struct {
		RoleID      string   `json:"role_id" binding:"required"`
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证角色是否存在
	if exist, err := h.getRoleExistByID(req.RoleID); err != nil {
		response.HandlerError(c, err)
		return
	} else if !exist {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 验证权限是否存在
	if exist, err := h.permissionService.CheckCodesExists(req.Permissions); err != nil || !exist {
		response.HandlerError(c, app_error.ErrPermissionNoExist)
		return
	}

	// 从角色撤销权限
	if err := h.roleService.RevokePermissionFromRole(req.RoleID, req.Permissions); err != nil {
		response.HandlerError(c, err)
		return
	}

	response.HandlerSuccess(c, gin.H{"message": "permissions revoked from role successfully"})
}
