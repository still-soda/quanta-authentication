package handlers

import (
	"crypto/rand"
	"encoding/base64"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  *services.UserService
	roleService  *services.RoleService
	auditService *services.AuditService
}

func NewUserHandler(
	userService *services.UserService,
	roleService *services.RoleService,
	auditService *services.AuditService,
) *UserHandler {
	return &UserHandler{
		userService:  userService,
		roleService:  roleService,
		auditService: auditService,
	}
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserList}); err != nil {
		utilities.GetLogger().Error("permission denied when listing users", "error", err.Error())
		response.HandlerError(c, err)
		return
	}

	// 获取查询参数
	var params services.ListUsersParams
	params.Page = utilities.ParseIntParam(c.Query("page"), 1)
	params.PageSize = utilities.ParseIntParam(c.Query("page_size"), 10)
	params.Search = c.Query("search")
	params.Status = c.Query("status")
	params.RoleID = c.Query("role_id")
	params.SortBy = c.Query("sort_by")
	params.SortDesc = c.Query("sort_desc") == "true"

	result, err := h.userService.ListUsers(params)
	if err != nil {
		utilities.GetLogger().Error("failed to list users", "error", err.Error())
		response.HandlerError(c, app_error.ErrFailedToGetUsers)
		return
	}

	response.HandlerSuccess(c, result)
}

// GetUser 获取单个用户
func (h *UserHandler) GetUser(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	user, err := h.userService.GetUserWithStats(userID)
	if err != nil {
		utilities.GetLogger().Error("failed to get user", "user_id", userID, "error", err.Error())
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	response.HandlerSuccess(c, user)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserCreate}); err != nil {
		response.HandlerError(c, err)
		return
	}

	var req struct {
		StudentID   string   `json:"student_id" binding:"required"`
		Name        string   `json:"name" binding:"required"`
		Email       string   `json:"email" binding:"required,email"`
		Password    string   `json:"password"` // 如果为空，则生成随机密码
		Phone       *string  `json:"phone"`
		DisplayName *string  `json:"display_name"`
		RoleIDs     []string `json:"role_ids"` // 可选的角色 ID 列表
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 如果没有提供密码，生成随机密码
	password := req.Password
	if password == "" {
		randomBytes := make([]byte, 12)
		if _, err := rand.Read(randomBytes); err != nil {
			response.HandlerError(c, app_error.ErrInternalServerError)
			return
		}
		password = base64.URLEncoding.EncodeToString(randomBytes)[:12]
	}

	// 创建用户
	user, err := h.userService.CreateUser(&services.CreateUserParams{
		StudentID:   req.StudentID,
		Password:    password,
		Email:       req.Email,
		Name:        req.Name,
		Phone:       req.Phone,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		utilities.GetLogger().Error("failed to create user", "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionUserCreate,
			TargetType:   "user",
			TargetName:   req.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})

		// 检查是否是唯一约束错误
		if app_error.IsErrorWithPgsqlCode(err, "23505") {
			response.HandlerError(c, app_error.ErrCreateUserConflict)
			return
		}
		response.HandlerError(c, app_error.ErrFailedToCreateUser)
		return
	}

	// 如果提供了角色 ID，分配角色
	if len(req.RoleIDs) > 0 {
		if err := h.userService.SetUserRoles(user.ID, req.RoleIDs); err != nil {
			utilities.GetLogger().Error("failed to assign roles to user", "user_id", user.ID, "error", err.Error())
			// 不中断流程，用户已创建成功
		}
	}

	// 重新获取用户以包含角色信息
	userWithStats, _ := h.userService.GetUserWithStats(user.ID)

	// 记录审计日志
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionUserCreate,
		TargetID:   user.ID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"student_id": req.StudentID,
				"name":       req.Name,
				"email":      req.Email,
				"role_ids":   req.RoleIDs,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, userWithStats)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserUpdate}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	var req struct {
		Name        *string `json:"name"`
		Email       *string `json:"email"`
		Phone       *string `json:"phone"`
		DisplayName *string `json:"display_name"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 获取旧用户信息
	oldUser, err := h.userService.GetUserByID(userID, false)
	if err != nil {
		utilities.GetLogger().Error("user not found", "user_id", userID)
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	// 构建更新参数
	updateParams := services.UpdateUserParams{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		DisplayName: req.DisplayName,
	}
	if req.Status != nil {
		status := models.UserStatus(*req.Status)
		updateParams.Status = &status
	}

	// 更新用户
	user, err := h.userService.UpdateUserByID(userID, updateParams)
	if err != nil {
		utilities.GetLogger().Error("failed to update user", "user_id", userID, "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionUserUpdate,
			TargetID:     userID,
			TargetType:   "user",
			TargetName:   oldUser.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrFailedToUpdateUser)
		return
	}

	// 记录审计日志
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionUserUpdate,
		TargetID:   userID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"name":         oldUser.Name,
				"email":        oldUser.Email,
				"phone":        oldUser.Phone,
				"display_name": oldUser.DisplayName,
				"status":       oldUser.Status,
			},
			NewValue: map[string]any{
				"name":         user.Name,
				"email":        user.Email,
				"phone":        user.Phone,
				"display_name": user.DisplayName,
				"status":       user.Status,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	// 获取带统计信息的用户
	userWithStats, _ := h.userService.GetUserWithStats(userID)
	response.HandlerSuccess(c, userWithStats)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserDelete}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID, false)
	if err != nil {
		utilities.GetLogger().Error("user not found", "user_id", userID)
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	// 删除用户
	if err := h.userService.DeleteUser(userID); err != nil {
		utilities.GetLogger().Error("failed to delete user", "user_id", userID, "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionUserDelete,
			TargetID:     userID,
			TargetType:   "user",
			TargetName:   user.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrFailedToDeleteUser)
		return
	}

	// 记录审计日志
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionUserDelete,
		TargetID:   userID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{
				"student_id": user.StudentID,
				"name":       user.Name,
				"email":      user.Email,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, nil)
}

// GetUserRoles 获取用户角色
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	roles, err := h.userService.GetUserRoles(userID)
	if err != nil {
		utilities.GetLogger().Error("failed to get user roles", "user_id", userID, "error", err.Error())
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	response.HandlerSuccess(c, roles)
}

// SetUserRoles 设置用户角色（替换）
func (h *UserHandler) SetUserRoles(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserAssignRoles}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	var req struct {
		RoleIDs []string `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID, true)
	if err != nil {
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	// 获取旧角色
	oldRoleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		oldRoleNames[i] = role.Name
	}

	// 设置新角色
	if err := h.userService.SetUserRoles(userID, req.RoleIDs); err != nil {
		utilities.GetLogger().Error("failed to set user roles", "user_id", userID, "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionUserUpdate,
			TargetID:     userID,
			TargetType:   "user",
			TargetName:   user.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrFailedToSetUserRoles)
		return
	}

	// 获取新角色
	newRoles, _ := h.userService.GetUserRoles(userID)
	newRoleNames := make([]string, len(newRoles))
	for i, role := range newRoles {
		newRoleNames[i] = role.Name
	}

	// 记录审计日志
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionUserUpdate,
		TargetID:   userID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			OldValue: map[string]any{"roles": oldRoleNames},
			NewValue: map[string]any{"roles": newRoleNames},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, newRoles)
}

// AssignRolesToUser 为用户分配角色（追加）
func (h *UserHandler) AssignRolesToUser(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserAssignRoles}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	var req struct {
		RoleIDs []string `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID, false)
	if err != nil {
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	// 分配角色
	if err := h.userService.AssignRolesToUser(userID, req.RoleIDs); err != nil {
		utilities.GetLogger().Error("failed to assign roles to user", "user_id", userID, "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionUserUpdate,
			TargetID:     userID,
			TargetType:   "user",
			TargetName:   user.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrFailedToSetUserRoles)
		return
	}

	// 获取更新后的角色
	roles, _ := h.userService.GetUserRoles(userID)

	// 记录审计日志
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionUserUpdate,
		TargetID:   userID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"assigned_role_ids": req.RoleIDs,
				"current_roles":     roleNames,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, roles)
}

// RevokeRolesFromUser 从用户撤销角色
func (h *UserHandler) RevokeRolesFromUser(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserRevokeRoles}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	var req struct {
		RoleIDs []string `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID, false)
	if err != nil {
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	// 撤销角色
	if err := h.userService.RevokeRolesFromUser(userID, req.RoleIDs); err != nil {
		utilities.GetLogger().Error("failed to revoke roles from user", "user_id", userID, "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionUserUpdate,
			TargetID:     userID,
			TargetType:   "user",
			TargetName:   user.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrFailedToSetUserRoles)
		return
	}

	// 获取更新后的角色
	roles, _ := h.userService.GetUserRoles(userID)

	// 记录审计日志
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionUserUpdate,
		TargetID:   userID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			NewValue: map[string]any{
				"revoked_role_ids": req.RoleIDs,
				"current_roles":    roleNames,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, roles)
}

// ResetUserPassword 重置用户密码
func (h *UserHandler) ResetUserPassword(c *gin.Context) {
	startTime := time.Now()

	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserResetPassword}); err != nil {
		response.HandlerError(c, err)
		return
	}

	userID := c.Param("id")

	var req struct {
		NewPassword string `json:"new_password"` // 如果为空，生成随机密码
	}
	c.ShouldBindJSON(&req)

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID, false)
	if err != nil {
		response.HandlerError(c, app_error.ErrUserNotFound)
		return
	}

	// 生成新密码
	newPassword := req.NewPassword
	if newPassword == "" {
		randomBytes := make([]byte, 12)
		if _, err := rand.Read(randomBytes); err != nil {
			response.HandlerError(c, app_error.ErrInternalServerError)
			return
		}
		newPassword = base64.URLEncoding.EncodeToString(randomBytes)[:12]
	}

	// 重置密码
	if err := h.userService.ResetPassword(userID, newPassword); err != nil {
		utilities.GetLogger().Error("failed to reset user password", "user_id", userID, "error", err.Error())
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleUser,
			Action:       models.AuditActionPasswordReset,
			TargetID:     userID,
			TargetType:   "user",
			TargetName:   user.Name,
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrFailedToResetPassword)
		return
	}

	// 记录审计日志
	h.auditService.LogWithGinContext(c, &services.AuditEntry{
		Module:     models.AuditModuleUser,
		Action:     models.AuditActionPasswordReset,
		TargetID:   userID,
		TargetType: "user",
		TargetName: user.Name,
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	// 返回新密码（只在请求中没有提供密码时返回）
	result := map[string]any{
		"success": true,
	}
	if req.NewPassword == "" {
		result["new_password"] = newPassword
	}

	response.HandlerSuccess(c, result)
}

// GetUserStatusCounts 获取用户状态统计
func (h *UserHandler) GetUserStatusCounts(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.UserList}); err != nil {
		response.HandlerError(c, err)
		return
	}

	counts, err := h.userService.GetUserStatusCounts()
	if err != nil {
		utilities.GetLogger().Error("failed to get user status counts", "error", err.Error())
		response.HandlerError(c, app_error.ErrFailedToGetUsers)
		return
	}

	response.HandlerSuccess(c, counts)
}
