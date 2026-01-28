package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	roleService       *services.RoleService
	permissionService *services.PermissionService
	auditService      *services.AuditService
}

func NewPermissionHandler(
	roleService *services.RoleService,
	permissionService *services.PermissionService,
	auditService *services.AuditService,
) *PermissionHandler {
	return &PermissionHandler{
		roleService:       roleService,
		permissionService: permissionService,
		auditService:      auditService,
	}
}

// GetPermissions 获取所有权限
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	// 验证是否有查看角色权限（查看权限需要角色查看权限）
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleView}); err != nil {
		utilities.GetLogger().Error("permission denied when getting permissions", "error", err.Error())
		response.HandlerError(c, err)
		return
	}

	// 获取查询参数
	var params services.ListPermissionsParams
	params.Page = utilities.ParseIntParam(c.Query("page"), 1)
	params.PageSize = utilities.ParseIntParam(c.Query("page_size"), 15)
	params.Search = c.Query("search")
	params.Resource = c.Query("resource")
	params.SortField = c.Query("sort_field")
	params.SortOrder = c.Query("sort_order")

	// 如果明确请求所有数据（用于兼容），则返回所有权限
	if c.Query("all") == "true" {
		perms, err := h.permissionService.GetAllPermissions()
		if err != nil {
			utilities.GetLogger().Error("failed to get permissions", "error", err.Error())
			response.HandlerError(c, app_error.ErrInternalServerError)
			return
		}
		response.HandlerSuccess(c, perms)
		return
	}

	// 分页查询
	result, err := h.permissionService.ListPermissions(params)
	if err != nil {
		utilities.GetLogger().Error("failed to list permissions", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, result)
}

// GetPermissionsGrouped 获取按资源分组的权限
func (h *PermissionHandler) GetPermissionsGrouped(c *gin.Context) {
	// 验证是否有查看角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleView}); err != nil {
		utilities.GetLogger().Error("permission denied when getting grouped permissions", "error", err.Error())
		response.HandlerError(c, err)
		return
	}

	// 获取按资源分组的权限
	grouped, err := h.permissionService.GetPermissionsGroupedByResource()
	if err != nil {
		utilities.GetLogger().Error("failed to get grouped permissions", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, grouped)
}

// GetPermission 获取单个权限
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	// 验证是否有查看角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	permissionID := c.Param("id")

	// 获取权限
	perm, err := h.permissionService.GetPermissionByID(permissionID)
	if err != nil {
		utilities.GetLogger().Error("failed to get permission", "error", err.Error())
		response.HandlerError(c, app_error.ErrPermissionNoExist)
		return
	}

	response.HandlerSuccess(c, perm)
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	startTime := time.Now()

	// 验证是否有创建角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleCreate}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 绑定请求参数
	var req struct {
		Resource    string `json:"resource" binding:"required"`
		Action      int8   `json:"action" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 创建权限
	perm, err := h.permissionService.CreatePermission(req.Resource, req.Action, req.Code, req.Description)
	if err != nil {
		utilities.GetLogger().Error("failed to create permission", "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModulePermission,
			Action:       models.AuditActionPermissionCreate,
			TargetType:   "permission",
			TargetName:   req.Code,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录权限创建
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModulePermission,
		Action:     models.AuditActionPermissionCreate,
		TargetID:   perm.ID,
		TargetType: "permission",
		TargetName: perm.Code,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"resource":    req.Resource,
				"action":      req.Action,
				"code":        req.Code,
				"description": req.Description,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, perm)
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	startTime := time.Now()

	// 验证是否有更新权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleUpdate}); err != nil {
		response.HandlerError(c, err)
		return
	}

	permissionID := c.Param("id")

	// 绑定请求参数
	var req struct {
		Resource    string `json:"resource" binding:"required"`
		Action      int8   `json:"action" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 获取旧权限信息
	oldPerm, err := h.permissionService.GetPermissionByID(permissionID)
	if err != nil {
		utilities.GetLogger().Error("failed to get permission before update", "error", err.Error())
		response.HandlerError(c, app_error.ErrPermissionNoExist)
		return
	}

	// 更新权限
	perm, err := h.permissionService.UpdatePermission(permissionID, req.Resource, req.Action, req.Code, req.Description)
	if err != nil {
		utilities.GetLogger().Error("failed to update permission", "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModulePermission,
			Action:       models.AuditActionPermissionUpdate,
			TargetID:     permissionID,
			TargetType:   "permission",
			TargetName:   oldPerm.Code,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录权限更新
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModulePermission,
		Action:     models.AuditActionPermissionUpdate,
		TargetID:   perm.ID,
		TargetType: "permission",
		TargetName: perm.Code,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"resource":    oldPerm.Resource,
				"action":      oldPerm.Action,
				"code":        oldPerm.Code,
				"description": oldPerm.Description,
			},
			NewValue: map[string]any{
				"resource":    req.Resource,
				"action":      req.Action,
				"code":        req.Code,
				"description": req.Description,
			},
			ChangedFields: []string{"resource", "action", "code", "description"},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, perm)
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	startTime := time.Now()

	// 验证是否有删除权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleDelete}); err != nil {
		response.HandlerError(c, err)
		return
	}

	permissionID := c.Param("id")

	// 获取权限信息
	perm, err := h.permissionService.GetPermissionByID(permissionID)
	if err != nil {
		utilities.GetLogger().Error("failed to get permission before delete", "error", err.Error())
		response.HandlerError(c, app_error.ErrPermissionNoExist)
		return
	}

	// 删除权限
	if err := h.permissionService.DeletePermission(permissionID); err != nil {
		utilities.GetLogger().Error("failed to delete permission", "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModulePermission,
			Action:       models.AuditActionPermissionDelete,
			TargetID:     permissionID,
			TargetType:   "permission",
			TargetName:   perm.Code,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 记录权限删除
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModulePermission,
		Action:     models.AuditActionPermissionDelete,
		TargetID:   permissionID,
		TargetType: "permission",
		TargetName: perm.Code,
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"message": "permission deleted successfully"})
}

// GetRolePermissions 获取角色的权限
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	// 验证是否有查看角色权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	roleID := c.Param("id")

	// 获取角色权限
	perms, err := h.permissionService.GetRolePermissions(roleID)
	if err != nil {
		utilities.GetLogger().Error("failed to get role permissions", "error", err.Error())
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	response.HandlerSuccess(c, perms)
}

// SetRolePermissions 设置角色的权限（替换所有权限）
func (h *PermissionHandler) SetRolePermissions(c *gin.Context) {
	startTime := time.Now()

	// 验证是否有授权权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.RoleAssignPermissions}); err != nil {
		response.HandlerError(c, err)
		return
	}

	roleID := c.Param("id")

	// 绑定请求参数
	var req struct {
		PermissionCodes []string `json:"permission_codes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证角色是否存在
	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		utilities.GetLogger().Error("failed to get role", "error", err.Error())
		response.HandlerError(c, app_error.ErrRoleNoExist)
		return
	}

	// 获取当前权限
	currentPerms, _ := h.permissionService.GetRolePermissions(roleID)
	currentCodes := make([]string, len(currentPerms))
	for i, p := range currentPerms {
		currentCodes[i] = p.Code
	}

	// 验证新权限是否存在（如果有新权限的话）
	if len(req.PermissionCodes) > 0 {
		if exist, err := h.permissionService.CheckCodesExists(req.PermissionCodes); err != nil || !exist {
			utilities.GetLogger().Error("permission does not exist when setting role permissions", "permissions", req.PermissionCodes)
			response.HandlerError(c, app_error.ErrPermissionNoExist)
			return
		}
	}

	// 先撤销所有权限
	if len(currentCodes) > 0 {
		if err := h.roleService.RevokePermissionFromRole(roleID, currentCodes); err != nil {
			utilities.GetLogger().Error("failed to revoke permissions from role", "error", err.Error())
			response.HandlerError(c, app_error.ErrInternalServerError)
			return
		}
	}

	// 授予新权限
	if len(req.PermissionCodes) > 0 {
		if err := h.roleService.GrantPermissionToRole(roleID, req.PermissionCodes); err != nil {
			utilities.GetLogger().Error("failed to grant permissions to role", "error", err.Error())
			h.auditService.LogWithGinContext(c, &services.AuditEntry{
				Module:       models.AuditModulePermission,
				Action:       models.AuditActionPermissionGrant,
				TargetID:     roleID,
				TargetType:   "role",
				TargetName:   role.Name,
				Status:       models.AuditStatusError,
				ErrorMessage: err.Error(),
				DurationMs:   time.Since(startTime).Milliseconds(),
			})
			response.HandlerError(c, app_error.ErrInternalServerError)
			return
		}
	}

	// 记录权限设置
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModulePermission,
		Action:     models.AuditActionPermissionGrant,
		TargetID:   roleID,
		TargetType: "role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"permissions": currentCodes,
			},
			NewValue: map[string]any{
				"permissions": req.PermissionCodes,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{"message": "role permissions updated successfully"})
}
