package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/services"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
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

	response.HandlerSuccess(c, user)
}
