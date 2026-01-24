package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService       *services.RoleService
	permissionService *services.PermissionService
	auditService      *services.AuditService
}

func NewRoleHandler(
	roleService *services.RoleService,
	permissionService *services.PermissionService,
	auditService *services.AuditService,
) *RoleHandler {
	return &RoleHandler{
		roleService:       roleService,
		permissionService: permissionService,
		auditService:      auditService,
	}
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
	startTime := time.Now()

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
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleRole,
			Action:       models.AuditActionRoleCreate,
			TargetType:   "role",
			TargetName:   req.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, err)
		return
	}

	// 记录角色创建
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleRole,
		Action:     models.AuditActionRoleCreate,
		TargetID:   role.ID,
		TargetType: "role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"name":        req.Name,
				"code":        req.Code,
				"permissions": req.Permissions,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, role)
}

// UpdateRole 处理更新角色请求
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	startTime := time.Now()

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
	oldRole, err := h.roleService.GetRoleByID(req.RoleID)
	if err != nil {
		response.HandlerError(c, err)
		return
	}
	if oldRole == nil {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 更新角色
	role, err := h.roleService.UpdateRole(req.RoleID, req.Name, req.Code)
	if err != nil {
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleRole,
			Action:       models.AuditActionRoleUpdate,
			TargetID:     req.RoleID,
			TargetType:   "role",
			TargetName:   oldRole.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, err)
		return
	}

	// 记录角色更新
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleRole,
		Action:     models.AuditActionRoleUpdate,
		TargetID:   role.ID,
		TargetType: "role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"name": oldRole.Name,
				"code": oldRole.Code,
			},
			NewValue: map[string]any{
				"name": req.Name,
				"code": req.Code,
			},
			ChangedFields: []string{"name", "code"},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, role)
}

// DeleteRole 处理删除角色请求
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	startTime := time.Now()

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

	// 验证角色是否存在并获取信息
	role, err := h.roleService.GetRoleByID(req.RoleID)
	if err != nil {
		response.HandlerError(c, err)
		return
	}
	if role == nil {
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 删除角色
	if err := h.roleService.DeleteRole(req.RoleID); err != nil {
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleRole,
			Action:       models.AuditActionRoleDelete,
			TargetID:     req.RoleID,
			TargetType:   "role",
			TargetName:   role.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, err)
		return
	}

	// 记录角色删除
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleRole,
		Action:     models.AuditActionRoleDelete,
		TargetID:   req.RoleID,
		TargetType: "role",
		TargetName: role.Name,
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"message": "role deleted successfully"})
}

// GrantPermissionsToRole 处理为角色授予权限请求
func (h *RoleHandler) GrantPermissionsToRole(c *gin.Context) {
	startTime := time.Now()

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
	role, err := h.roleService.GetRoleByID(req.RoleID)
	if err != nil {
		response.HandlerError(c, err)
		return
	}
	if role == nil {
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
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModulePermission,
			Action:       models.AuditActionPermissionGrant,
			TargetID:     req.RoleID,
			TargetType:   "role",
			TargetName:   role.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, err)
		return
	}

	// 记录权限授予
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModulePermission,
		Action:     models.AuditActionPermissionGrant,
		TargetID:   req.RoleID,
		TargetType: "role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"permissions": req.Permissions,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"message": "permissions granted to role successfully"})
}

// RevokePermissionsFromRole 处理从角色撤销权限请求
func (h *RoleHandler) RevokePermissionsFromRole(c *gin.Context) {
	startTime := time.Now()

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
	role, err := h.roleService.GetRoleByID(req.RoleID)
	if err != nil {
		response.HandlerError(c, err)
		return
	}
	if role == nil {
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
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModulePermission,
			Action:       models.AuditActionPermissionRevoke,
			TargetID:     req.RoleID,
			TargetType:   "role",
			TargetName:   role.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, err)
		return
	}

	// 记录权限撤销
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModulePermission,
		Action:     models.AuditActionPermissionRevoke,
		TargetID:   req.RoleID,
		TargetType: "role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"permissions": req.Permissions,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"message": "permissions revoked from role successfully"})
}
