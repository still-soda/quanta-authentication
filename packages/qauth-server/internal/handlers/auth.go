package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/services"
	"qauth-server/pkg/jwt"
	"qauth-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService  *services.UserService
	roleService  *services.RoleService
	auditService *services.AuditService
}

func NewAuthHandler(
	userService *services.UserService,
	roleService *services.RoleService,
	auditService *services.AuditService,
) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		roleService:  roleService,
		auditService: auditService,
	}
}

// Register 处理用户注册请求
func (h *AuthHandler) Register(c *gin.Context) {
	startTime := time.Now()

	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Name      string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	user, err := h.userService.CreateUser(&services.CreateUserParams{
		StudentID: req.StudentID,
		Password:  req.Password,
		Email:     req.Email,
		Name:      req.Name,
	})
	if err != nil {
		// 记录注册失败的审计日志
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:     models.AuditModuleAuth,
			Action:     models.AuditActionRegister,
			TargetType: "user",
			Detail: &models.AuditLogDetail{
				Metadata: map[string]any{
					"student_id": req.StudentID,
					"email":      req.Email,
				},
				FailReason: err.Error(),
			},
			Status:       models.AuditStatusError,
			ErrorMessage: err.Error(),
			DurationMs:   time.Since(startTime).Milliseconds(),
		})

		// 检查是否为唯一性冲突错误
		if app_error.IsErrorWithPgsqlCode(err, "23505") {
			response.HandlerError(c, app_error.ErrCreateUserConflict)
			return
		}
		response.HandlerError(c, app_error.ErrFailedToCreateUser)
		return
	}

	// 记录注册成功的审计日志
	h.auditService.Log(&services.AuditContext{
		OperatorID:   user.ID,
		OperatorName: user.Name,
		IP:           c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
	}, &services.AuditEntry{
		Module:     models.AuditModuleAuth,
		Action:     models.AuditActionRegister,
		TargetID:   user.ID,
		TargetType: "user",
		TargetName: user.Name,
		Detail: &models.AuditLogDetail{
			Metadata: map[string]any{
				"student_id": req.StudentID,
				"email":      req.Email,
			},
		},
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, user)
}

// Login 处理用户登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	startTime := time.Now()

	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	user, err := h.userService.AuthenticateUser(req.StudentID, req.Password)
	if err != nil {
		// 记录登录失败
		h.auditService.LogLogin(c, "", req.StudentID, "PASSWORD", false, "authentication failed", time.Since(startTime).Milliseconds())
		response.HandlerError(c, app_error.ErrAuthenticationFailed)
		return
	}

	role, err := h.roleService.GetUserRole(user.ID)
	if err != nil {
		h.auditService.LogLogin(c, user.ID, user.Name, "PASSWORD", false, "failed to get user role", time.Since(startTime).Milliseconds())
		response.HandlerError(c, app_error.ErrFailedToGetUserRole)
		return
	}

	info := &jwt.JWTInfo{
		UserID:    user.ID,
		StudentID: user.StudentID,
		Role:      role.Code,
	}
	accessToken, err := jwt.GenerateAccessToken(info)
	refreshToken, err := jwt.GenerateRefreshToken(info)
	if err != nil {
		h.auditService.LogLogin(c, user.ID, user.Name, "PASSWORD", false, "failed to generate token", time.Since(startTime).Milliseconds())
		response.HandlerError(c, app_error.ErrFailedToGenerateToken)
		return
	}

	// 记录登录成功
	h.auditService.LogLogin(c, user.ID, user.Name, "PASSWORD", true, "", time.Since(startTime).Milliseconds())

	response.HandlerSuccess(c, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken 用于滑动续签令牌，生成新的访问令牌和刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	startTime := time.Now()

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	userInfo, err := jwt.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		h.auditService.LogWithGinContext(c, &services.AuditEntry{
			Module:       models.AuditModuleAuth,
			Action:       models.AuditActionTokenRefresh,
			Status:       models.AuditStatusError,
			ErrorMessage: "invalid refresh token",
			DurationMs:   time.Since(startTime).Milliseconds(),
		})
		response.HandlerError(c, app_error.ErrInvalidRefreshToken)
		return
	}

	info := &jwt.JWTInfo{
		UserID:    userInfo.UserID,
		StudentID: userInfo.StudentID,
		Role:      userInfo.Role,
	}
	newAccessToken, err := jwt.GenerateAccessToken(info)
	newRefreshToken, err := jwt.GenerateRefreshToken(info)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGenerateToken)
		return
	}

	// 记录令牌刷新
	h.auditService.Log(&services.AuditContext{
		OperatorID: userInfo.UserID,
		IP:         c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}, &services.AuditEntry{
		Module:     models.AuditModuleAuth,
		Action:     models.AuditActionTokenRefresh,
		TargetID:   userInfo.UserID,
		TargetType: "user",
		Status:     models.AuditStatusSuccess,
		DurationMs: time.Since(startTime).Milliseconds(),
	})

	response.HandlerSuccess(c, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
