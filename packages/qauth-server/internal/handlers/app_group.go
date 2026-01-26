package handlers

import (
	"fmt"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwt"
	"qauth-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// AppGroupHandler 应用组权限处理器
type AppGroupHandler struct {
	appGroupService *services.AppGroupService
	oauthService    *services.OAuthService
	roleService     *services.RoleService
	auditService    *services.AuditService
}

// NewAppGroupHandler 创建应用组权限处理器
func NewAppGroupHandler(
	appGroupService *services.AppGroupService,
	oauthService *services.OAuthService,
	roleService *services.RoleService,
	auditService *services.AuditService,
) *AppGroupHandler {
	return &AppGroupHandler{
		appGroupService: appGroupService,
		oauthService:    oauthService,
		roleService:     roleService,
		auditService:    auditService,
	}
}

// getUserInfo 从上下文获取用户信息
func (h *AppGroupHandler) getUserInfo(c *gin.Context) *jwt.UserJWTClaims {
	if info, exists := c.Get("userInfo"); exists {
		return info.(*jwt.UserJWTClaims)
	}
	return nil
}

// checkAppGroupAdminPermission 检查用户是否有应用组管理权限
func (h *AppGroupHandler) checkAppGroupAdminPermission(c *gin.Context, clientID string, requiredTypes ...models.AppGroupAdminType) error {
	userInfo := h.getUserInfo(c)
	if userInfo == nil {
		return app_error.ErrUnauthorized
	}

	// 超级管理员拥有所有权限
	if userInfo.Role == "system_super_admin" {
		return nil
	}

	// 检查是否是应用组管理员
	hasPermission, err := h.appGroupService.HasAppGroupAdminPermission(clientID, userInfo.UserID, requiredTypes...)
	if err != nil {
		return app_error.ErrInternalServerError
	}
	if !hasPermission {
		return app_error.ErrNoPermission
	}

	return nil
}

// checkClientExists 检查客户端是否存在
func (h *AppGroupHandler) checkClientExists(clientID string) error {
	_, err := h.oauthService.GetClientByID(clientID)
	if err != nil {
		return app_error.ErrOAuthClientNotFound
	}
	return nil
}

// ======================== 应用组管理员相关 API ========================

// GetAppGroupAdmins 获取应用组管理员列表
// GET /_/v1/clients/:client_id/admins
func (h *AppGroupHandler) GetAppGroupAdmins(c *gin.Context) {
	clientID := c.Param("id")

	// 检查权限（owner、admin 和超级管理员可查看）
	if err := h.checkAppGroupAdminPermission(c, clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	admins, err := h.appGroupService.GetAppGroupAdmins(clientID)
	if err != nil {
		utilities.GetLogger().Error("failed to get app group admins", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	result := make([]*models.AppGroupAdminResponse, len(admins))
	for i, admin := range admins {
		result[i] = admin.ToResponse()
	}

	response.HandlerSuccess(c, result)
}

// AddAppGroupAdmin 添加应用组管理员
// POST /_/v1/clients/:client_id/admins
func (h *AppGroupHandler) AddAppGroupAdmin(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")

	// 检查权限（只有 owner、admin 和超级管理员可添加）
	if err := h.checkAppGroupAdminPermission(c, clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		UserID    string                   `json:"user_id" binding:"required"`
		AdminType models.AppGroupAdminType `json:"admin_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 验证管理员类型
	validTypes := map[models.AppGroupAdminType]bool{
		models.AppGroupAdminTypeAdmin:             true,
		models.AppGroupAdminTypeRoleManager:       true,
		models.AppGroupAdminTypePermissionManager: true,
	}
	if !validTypes[req.AdminType] {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	userInfo := h.getUserInfo(c)
	if err := h.appGroupService.AddAppGroupAdmin(clientID, req.UserID, req.AdminType, userInfo.UserID); err != nil {
		utilities.GetLogger().Error("failed to add app group admin", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupAdminAdd,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupAdminAdd,
		TargetID:   clientID,
		TargetType: "oauth_client",
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"user_id":    req.UserID,
				"admin_type": req.AdminType,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// RemoveAppGroupAdmin 移除应用组管理员
// DELETE /_/v1/clients/:client_id/admins/:user_id
func (h *AppGroupHandler) RemoveAppGroupAdmin(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	userID := c.Param("user_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	adminType := models.AppGroupAdminType(c.Query("admin_type"))
	if adminType == "" {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 不能移除 owner
	if adminType == models.AppGroupAdminTypeOwner {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	if err := h.appGroupService.RemoveAppGroupAdmin(clientID, userID, adminType); err != nil {
		utilities.GetLogger().Error("failed to remove app group admin", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupAdminRemove,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupAdminRemove,
		TargetID:   clientID,
		TargetType: "oauth_client",
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"user_id":    userID,
				"admin_type": adminType,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// ======================== 应用组权限相关 API ========================

// GetAppGroupPermissions 获取应用组权限列表
// GET /_/v1/clients/:client_id/permissions
func (h *AppGroupHandler) GetAppGroupPermissions(c *gin.Context) {
	clientID := c.Param("id")

	// 检查客户端是否存在
	if err := h.checkClientExists(clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 检查权限（所有应用组管理员都可以查看）
	if err := h.checkAppGroupAdminPermission(c, clientID,
		models.AppGroupAdminTypePermissionManager,
		models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	permissions, err := h.appGroupService.GetAppGroupPermissions(clientID)
	if err != nil {
		utilities.GetLogger().Error("failed to get app group permissions", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	result := make([]*models.AppGroupPermissionResponse, len(permissions))
	for i, perm := range permissions {
		result[i] = perm.ToResponse()
	}

	response.HandlerSuccess(c, result)
}

// GetAppGroupPermissionsGrouped 获取按资源分组的应用组权限
// GET /_/v1/clients/:client_id/permissions/grouped
func (h *AppGroupHandler) GetAppGroupPermissionsGrouped(c *gin.Context) {
	clientID := c.Param("id")

	// 检查客户端是否存在
	if err := h.checkClientExists(clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID,
		models.AppGroupAdminTypePermissionManager,
		models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	grouped, err := h.appGroupService.GetAppGroupPermissionsByResource(clientID)
	if err != nil {
		utilities.GetLogger().Error("failed to get grouped permissions", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	result := make(map[string][]*models.AppGroupPermissionResponse)
	for resource, perms := range grouped {
		result[resource] = make([]*models.AppGroupPermissionResponse, len(perms))
		for i, perm := range perms {
			result[resource][i] = perm.ToResponse()
		}
	}

	response.HandlerSuccess(c, result)
}

// CreateAppGroupPermission 创建应用组权限
// POST /_/v1/clients/:client_id/permissions
func (h *AppGroupHandler) CreateAppGroupPermission(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")

	// 检查客户端是否存在
	if err := h.checkClientExists(clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 检查权限（只有 owner、admin、permission_manager 和超级管理员可创建）
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		Resource    string `json:"resource" binding:"required"`
		Action      int8   `json:"action" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 生成完整的权限代码（防止全局冲突）
	fullCode := fmt.Sprintf("app_%s_%s", clientID[:8], req.Code)

	permission, err := h.appGroupService.CreateAppGroupPermission(
		clientID, req.Resource, req.Action, fullCode, req.Name, req.Description,
	)
	if err != nil {
		utilities.GetLogger().Error("failed to create app group permission", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupPermissionCreate,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupPermissionCreate,
		TargetID:   permission.ID,
		TargetType: "app_group_permission",
		TargetName: permission.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"resource":    req.Resource,
				"action":      req.Action,
				"code":        fullCode,
				"name":        req.Name,
				"description": req.Description,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, permission.ToResponse())
}

// UpdateAppGroupPermission 更新应用组权限
// PUT /_/v1/clients/:client_id/permissions/:permission_id
func (h *AppGroupHandler) UpdateAppGroupPermission(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	permissionID := c.Param("permission_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	permission, err := h.appGroupService.UpdateAppGroupPermission(permissionID, req.Name, req.Description)
	if err != nil {
		utilities.GetLogger().Error("failed to update app group permission", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupPermissionUpdate,
			TargetID:     permissionID,
			TargetType:   "app_group_permission",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupPermissionUpdate,
		TargetID:   permissionID,
		TargetType: "app_group_permission",
		TargetName: permission.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"name":        req.Name,
				"description": req.Description,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, permission.ToResponse())
}

// DeleteAppGroupPermission 删除应用组权限
// DELETE /_/v1/clients/:client_id/permissions/:permission_id
func (h *AppGroupHandler) DeleteAppGroupPermission(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	permissionID := c.Param("permission_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 获取权限信息用于审计
	permission, err := h.appGroupService.GetAppGroupPermission(permissionID)
	if err != nil {
		response.HandlerError(c, app_error.ErrNotFound)
		return
	}

	if err := h.appGroupService.DeleteAppGroupPermission(permissionID); err != nil {
		utilities.GetLogger().Error("failed to delete app group permission", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupPermissionDelete,
			TargetID:     permissionID,
			TargetType:   "app_group_permission",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupPermissionDelete,
		TargetID:   permissionID,
		TargetType: "app_group_permission",
		TargetName: permission.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"code": permission.Code,
				"name": permission.Name,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// ======================== 应用组角色相关 API ========================

// GetAppGroupRoles 获取应用组角色列表
// GET /_/v1/clients/:client_id/roles
func (h *AppGroupHandler) GetAppGroupRoles(c *gin.Context) {
	clientID := c.Param("id")

	// 检查客户端是否存在
	if err := h.checkClientExists(clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID,
		models.AppGroupAdminTypeRoleManager,
		models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	roles, err := h.appGroupService.GetAppGroupRolesWithStats(clientID)
	if err != nil {
		utilities.GetLogger().Error("failed to get app group roles", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, roles)
}

// GetAppGroupRole 获取应用组角色详情
// GET /_/v1/clients/:client_id/roles/:role_id
func (h *AppGroupHandler) GetAppGroupRole(c *gin.Context) {
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID,
		models.AppGroupAdminTypeRoleManager,
		models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	role, err := h.appGroupService.GetAppGroupRoleWithStats(roleID)
	if err != nil {
		utilities.GetLogger().Error("failed to get app group role", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	response.HandlerSuccess(c, role)
}

// CreateAppGroupRole 创建应用组角色
// POST /_/v1/clients/:client_id/roles
func (h *AppGroupHandler) CreateAppGroupRole(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")

	// 检查客户端是否存在
	if err := h.checkClientExists(clientID); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 检查权限（只有 owner、admin、role_manager 和超级管理员可创建）
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		Code        string `json:"code" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsDefault   bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 生成完整的角色代码（防止全局冲突）
	fullCode := fmt.Sprintf("app_%s_%s", clientID[:8], req.Code)

	role, err := h.appGroupService.CreateAppGroupRole(clientID, fullCode, req.Name, req.Description, req.IsDefault)
	if err != nil {
		utilities.GetLogger().Error("failed to create app group role", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupRoleCreate,
			TargetID:     clientID,
			TargetType:   "oauth_client",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupRoleCreate,
		TargetID:   role.ID,
		TargetType: "app_group_role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"code":        fullCode,
				"name":        req.Name,
				"description": req.Description,
				"is_default":  req.IsDefault,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, role.ToResponse())
}

// UpdateAppGroupRole 更新应用组角色
// PUT /_/v1/clients/:client_id/roles/:role_id
func (h *AppGroupHandler) UpdateAppGroupRole(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsDefault   bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	role, err := h.appGroupService.UpdateAppGroupRole(roleID, req.Name, req.Description, req.IsDefault)
	if err != nil {
		utilities.GetLogger().Error("failed to update app group role", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupRoleUpdate,
			TargetID:     roleID,
			TargetType:   "app_group_role",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupRoleUpdate,
		TargetID:   roleID,
		TargetType: "app_group_role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"name":        req.Name,
				"description": req.Description,
				"is_default":  req.IsDefault,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, role.ToResponse())
}

// DeleteAppGroupRole 删除应用组角色
// DELETE /_/v1/clients/:client_id/roles/:role_id
func (h *AppGroupHandler) DeleteAppGroupRole(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 获取角色信息用于审计
	role, err := h.appGroupService.GetAppGroupRole(roleID)
	if err != nil {
		response.HandlerError(c, app_error.ErrNotFound)
		return
	}

	if err := h.appGroupService.DeleteAppGroupRole(roleID); err != nil {
		utilities.GetLogger().Error("failed to delete app group role", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupRoleDelete,
			TargetID:     roleID,
			TargetType:   "app_group_role",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupRoleDelete,
		TargetID:   roleID,
		TargetType: "app_group_role",
		TargetName: role.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"code": role.Code,
				"name": role.Name,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// GetAppGroupRolePermissions 获取应用组角色的权限
// GET /_/v1/clients/:client_id/roles/:role_id/permissions
func (h *AppGroupHandler) GetAppGroupRolePermissions(c *gin.Context) {
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID,
		models.AppGroupAdminTypeRoleManager,
		models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	permissions, err := h.appGroupService.GetAppGroupRolePermissions(roleID)
	if err != nil {
		utilities.GetLogger().Error("failed to get role permissions", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	result := make([]*models.AppGroupPermissionResponse, len(permissions))
	for i, perm := range permissions {
		result[i] = perm.ToResponse()
	}

	response.HandlerSuccess(c, result)
}

// SetAppGroupRolePermissions 设置应用组角色的权限
// PUT /_/v1/clients/:client_id/roles/:role_id/permissions
func (h *AppGroupHandler) SetAppGroupRolePermissions(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限（需要有 role_manager 和 permission_manager 权限）
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		PermissionIDs []string `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	if err := h.appGroupService.SetAppGroupRolePermissions(roleID, req.PermissionIDs); err != nil {
		utilities.GetLogger().Error("failed to set role permissions", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupRoleAssignPermissions,
			TargetID:     roleID,
			TargetType:   "app_group_role",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupRoleAssignPermissions,
		TargetID:   roleID,
		TargetType: "app_group_role",
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"permission_ids": req.PermissionIDs,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// ======================== 应用组用户角色相关 API ========================

// GetAppGroupRoleUsers 获取应用组角色的用户列表
// GET /_/v1/clients/:client_id/roles/:role_id/users
func (h *AppGroupHandler) GetAppGroupRoleUsers(c *gin.Context) {
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	users, err := h.appGroupService.GetAppGroupRoleUsers(roleID)
	if err != nil {
		utilities.GetLogger().Error("failed to get role users", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 返回简化的用户信息
	type UserInfo struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		StudentID string `json:"student_id"`
	}

	result := make([]UserInfo, len(users))
	for i, user := range users {
		result[i] = UserInfo{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			StudentID: user.StudentID,
		}
	}

	response.HandlerSuccess(c, result)
}

// AssignAppGroupRoleToUser 为用户分配应用组角色
// POST /_/v1/clients/:client_id/roles/:role_id/users
func (h *AppGroupHandler) AssignAppGroupRoleToUser(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	roleID := c.Param("role_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	userInfo := h.getUserInfo(c)
	if err := h.appGroupService.AssignAppGroupRoleToUser(req.UserID, roleID, userInfo.UserID); err != nil {
		utilities.GetLogger().Error("failed to assign role to user", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupUserAssignRoles,
			TargetID:     roleID,
			TargetType:   "app_group_role",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupUserAssignRoles,
		TargetID:   roleID,
		TargetType: "app_group_role",
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"user_id": req.UserID,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// RevokeAppGroupRoleFromUser 从用户撤销应用组角色
// DELETE /_/v1/clients/:client_id/roles/:role_id/users/:user_id
func (h *AppGroupHandler) RevokeAppGroupRoleFromUser(c *gin.Context) {
	startTime := time.Now()
	clientID := c.Param("id")
	roleID := c.Param("role_id")
	userID := c.Param("user_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	if err := h.appGroupService.RevokeAppGroupRoleFromUser(userID, roleID); err != nil {
		utilities.GetLogger().Error("failed to revoke role from user", "error", err)
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleOAuth,
			Action:       models.AuditActionAppGroupUserRevokeRoles,
			TargetID:     roleID,
			TargetType:   "app_group_role",
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionAppGroupUserRevokeRoles,
		TargetID:   roleID,
		TargetType: "app_group_role",
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"user_id": userID,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// GetUserAppGroupRoles 获取用户在应用组的角色
// GET /_/v1/clients/:client_id/users/:user_id/roles
func (h *AppGroupHandler) GetUserAppGroupRoles(c *gin.Context) {
	clientID := c.Param("id")
	userID := c.Param("user_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID, models.AppGroupAdminTypeRoleManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	roles, err := h.appGroupService.GetUserAppGroupRoles(userID, clientID)
	if err != nil {
		utilities.GetLogger().Error("failed to get user roles", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	result := make([]*models.AppGroupRoleResponse, len(roles))
	for i, role := range roles {
		result[i] = role.ToResponse()
	}

	response.HandlerSuccess(c, result)
}

// GetUserAppGroupPermissions 获取用户在应用组的权限
// GET /_/v1/clients/:client_id/users/:user_id/permissions
func (h *AppGroupHandler) GetUserAppGroupPermissions(c *gin.Context) {
	clientID := c.Param("id")
	userID := c.Param("user_id")

	// 检查权限
	if err := h.checkAppGroupAdminPermission(c, clientID,
		models.AppGroupAdminTypeRoleManager,
		models.AppGroupAdminTypePermissionManager); err != nil {
		response.HandlerError(c, err)
		return
	}

	permissions, err := h.appGroupService.GetUserAppGroupPermissions(userID, clientID)
	if err != nil {
		utilities.GetLogger().Error("failed to get user permissions", "error", err)
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 转换为响应格式
	result := make([]*models.AppGroupPermissionResponse, len(permissions))
	for i, perm := range permissions {
		result[i] = perm.ToResponse()
	}

	response.HandlerSuccess(c, result)
}
