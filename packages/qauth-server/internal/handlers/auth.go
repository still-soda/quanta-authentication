package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/services"
	"qauth-server/pkg/jwt"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
	roleService *services.RoleService
}

func NewAuthHandler(
	userService *services.UserService,
	roleService *services.RoleService,
) *AuthHandler {
	return &AuthHandler{userService: userService, roleService: roleService}
}

// Register 处理用户注册请求
func (h *AuthHandler) Register(c *gin.Context) {
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
		// 检查是否为唯一性冲突错误
		if app_error.IsErrorWithPgsqlCode(err, "23505") {
			response.HandlerError(c, app_error.ErrCreateUserConflict)
			return
		}
		response.HandlerError(c, app_error.ErrFailedToCreateUser)
		return
	}

	response.HandlerSuccess(c, user)

}

// Login 处理用户登录请求
func (h *AuthHandler) Login(c *gin.Context) {
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
		response.HandlerError(c, app_error.ErrAuthenticationFailed)
		return
	}

	role, err := h.roleService.GetUserRole(user.ID)
	if err != nil {
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
		response.HandlerError(c, app_error.ErrFailedToGenerateToken)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken 用于滑动续签令牌，生成新的访问令牌和刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	userInfo, err := jwt.ParseRefreshToken(req.RefreshToken)
	if err != nil {
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

	response.HandlerSuccess(c, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
