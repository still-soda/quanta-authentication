package services

import (
	"errors"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/providers"
	"qauth-server/internal/repository"
	"qauth-server/internal/utilities"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	repo      *repository.UserRepository
	auditRepo *repository.AuditRepository
	logger    providers.ILogger
}

func NewUserService(
	repo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
	logger providers.ILogger,
) *UserService {
	return &UserService{
		repo:      repo,
		auditRepo: auditRepo,
		logger:    logger.With("service", "UserService"),
	}
}

// GetUserByID 根据用户ID获取用户信息
func (s *UserService) GetUserByID(userID string, withRole bool) (*models.Users, error) {
	var user *models.Users
	var err error

	if withRole {
		user, err = s.repo.FindByIDWithRoles(userID)
	} else {
		user, err = s.repo.FindByID(userID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound.WithScope("GetUserByID")
		}
		return nil, e.ErrFailedToGetUser.WithScope("GetUserByID").Wrap(err)
	}
	return user, nil
}

// GetUserByStudentID 根据学生ID获取用户信息
func (s *UserService) GetUserByStudentID(studentID string) (*models.Users, error) {
	user, err := s.repo.FindByStudentID(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound.WithScope("GetUserByStudentID")
		}
		return nil, e.ErrFailedToGetUser.WithScope("GetUserByStudentID").Wrap(err)
	}
	return user, nil
}

// AuthenticateUser 验证用户凭据
func (s *UserService) AuthenticateUser(studentID, password string) (*models.Users, error) {
	user, err := s.repo.FindByStudentID(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound.WithScope("AuthenticateUser")
		}
		return nil, e.ErrFailedToGetUser.WithScope("AuthenticateUser").Wrap(err)
	}

	verified := utilities.VerifyPassword(password, user.Salt, user.PasswordHash)
	if !verified {
		return nil, e.ErrAuthenticationFailed.WithScope("AuthenticateUser")
	}

	return user, nil
}

type CreateUserParams struct {
	StudentID   string
	Password    string
	Email       string
	Name        string
	Phone       *string
	DisplayName *string
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(params *CreateUserParams) (*models.Users, error) {
	salt, err := utilities.GenerateSalt(16)
	if err != nil {
		return nil, e.ErrFailedToGenerateSalt.WithScope("CreateUser").Wrap(err)
	}

	hashedPassword := utilities.HashPassword(params.Password, salt)

	user := &models.Users{
		StudentID:    params.StudentID,
		Email:        params.Email,
		Name:         params.Name,
		Salt:         salt,
		PasswordHash: hashedPassword,
		Phone:        params.Phone,
		DisplayName:  params.DisplayName,
		Status:       models.UserStatusActive,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, e.ErrFailedToCreateUser.WithScope("CreateUser").Wrap(err)
	}

	// 重新查询以获取完整数据
	createdUser, err := s.repo.FindByIDWithRoles(user.ID)
	if err != nil {
		return nil, e.ErrFailedToGetUser.WithScope("CreateUser").Wrap(err)
	}

	return createdUser, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.Users) error {
	if err := s.repo.Update(user); err != nil {
		return e.ErrFailedToUpdateUser.WithScope("UpdateUser").Wrap(err)
	}
	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID string) error {
	// 先删除用户角色关联
	if err := s.repo.DeleteUserRoles(userID); err != nil {
		return e.ErrFailedToDeleteUser.WithScope("DeleteUser").Wrap(err)
	}
	if err := s.repo.Delete(userID); err != nil {
		return e.ErrFailedToDeleteUser.WithScope("DeleteUser").Wrap(err)
	}
	return nil
}

// UserCount 返回所有用户的数量
func (s *UserService) UserCount() (int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return -1, e.ErrFailedToCountUsersTotal.WithScope("UserCount").Wrap(err)
	}
	return count, nil
}

// GetUserCountByRole 返回按角色分类的用户数量统计
func (s *UserService) GetUserCountByRole() (map[string]int64, error) {
	roleCountMap, err := s.repo.CountByRole()
	if err != nil {
		return nil, e.ErrFailedToCountUsersByRole.WithScope("GetUserCountByRole").Wrap(err)
	}
	return roleCountMap, nil
}

// UserWithStats 用户带统计信息
type UserWithStats struct {
	models.Users
	RoleNames   []string   `json:"role_names"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// ListUsersParams 用户列表查询参数
type ListUsersParams struct {
	Page     int
	PageSize int
	Search   string // 搜索关键词（姓名、邮箱、学号）
	Status   string // 用户状态筛选
	RoleID   string // 角色 ID 筛选
	SortBy   string // 排序字段
	SortDesc bool   // 是否降序
}

// ListUsersResult 用户列表结果
type ListUsersResult struct {
	Users      []UserWithStats `json:"users"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ListUsers 获取用户列表（带分页和筛选）
func (s *UserService) ListUsers(params ListUsersParams) (*ListUsersResult, error) {
	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// 构建 repository 参数
	repoParams := repository.ListUsersParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Search:   params.Search,
		Status:   params.Status,
		RoleID:   params.RoleID,
		SortBy:   params.SortBy,
		SortDesc: params.SortDesc,
	}

	users, total, err := s.repo.List(repoParams)
	if err != nil {
		return nil, e.ErrFailedToGetUsers.WithScope("ListUsers").Wrap(err)
	}

	// 转换为带统计信息的结果
	usersWithStats := make([]UserWithStats, len(users))
	for i, user := range users {
		roleNames := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roleNames[j] = role.Name
		}

		// 获取最后登录时间（从审计日志）
		lastLogin, err := s.auditRepo.FindLastLoginByUserID(user.ID)
		if err != nil {
			// 记录错误但不中断流程
			lastLogin = nil
		}

		usersWithStats[i] = UserWithStats{
			Users:       user,
			RoleNames:   roleNames,
			LastLoginAt: lastLogin,
		}
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &ListUsersResult{
		Users:      usersWithStats,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserWithStats 获取单个用户（带统计信息）
func (s *UserService) GetUserWithStats(userID string) (*UserWithStats, error) {
	user, err := s.repo.FindByIDWithRoles(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound.WithScope("GetUserWithStats")
		}
		return nil, e.ErrFailedToGetUser.WithScope("GetUserWithStats").Wrap(err)
	}

	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// 获取最后登录时间
	lastLogin, err := s.auditRepo.FindLastLoginByUserID(user.ID)
	if err != nil {
		// 记录错误但不中断流程
		lastLogin = nil
	}

	return &UserWithStats{
		Users:       *user,
		RoleNames:   roleNames,
		LastLoginAt: lastLogin,
	}, nil
}

// UpdateUserParams 更新用户参数
type UpdateUserParams struct {
	Name        *string
	Email       *string
	Phone       *string
	DisplayName *string
	Status      *models.UserStatus
}

// UpdateUserByID 根据 ID 更新用户信息
func (s *UserService) UpdateUserByID(userID string, params UpdateUserParams) (*models.Users, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound.WithScope("UpdateUserByID")
		}
		return nil, e.ErrFailedToGetUser.WithScope("UpdateUserByID").Wrap(err)
	}

	// 只更新非空字段
	if params.Name != nil {
		user.Name = *params.Name
	}
	if params.Email != nil {
		user.Email = *params.Email
	}
	if params.Phone != nil {
		user.Phone = params.Phone
	}
	if params.DisplayName != nil {
		user.DisplayName = params.DisplayName
	}
	if params.Status != nil {
		user.Status = *params.Status
	}

	if err := s.repo.Update(user); err != nil {
		return nil, e.ErrFailedToUpdateUser.WithScope("UpdateUserByID").Wrap(err)
	}

	// 重新查询以获取完整数据
	updatedUser, err := s.repo.FindByIDWithRoles(userID)
	if err != nil {
		return nil, e.ErrFailedToGetUser.WithScope("UpdateUserByID").Wrap(err)
	}

	return updatedUser, nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(userID string, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrUserNotFound.WithScope("ResetPassword")
		}
		return e.ErrFailedToGetUser.WithScope("ResetPassword").Wrap(err)
	}

	salt, err := utilities.GenerateSalt(16)
	if err != nil {
		return e.ErrFailedToGenerateSalt.WithScope("ResetPassword").Wrap(err)
	}

	hashedPassword := utilities.HashPassword(newPassword, salt)
	user.Salt = salt
	user.PasswordHash = hashedPassword

	if err := s.repo.Update(user); err != nil {
		return e.ErrFailedToResetPassword.WithScope("ResetPassword").Wrap(err)
	}
	return nil
}

// GetUserRoles 获取用户的所有角色
func (s *UserService) GetUserRoles(userID string) ([]models.Roles, error) {
	roles, err := s.repo.GetUserRoles(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound.WithScope("GetUserRoles")
		}
		return nil, e.ErrFailedToGetUserRoles.WithScope("GetUserRoles").Wrap(err)
	}
	return roles, nil
}

// SetUserRoles 设置用户的角色（替换现有角色）
func (s *UserService) SetUserRoles(userID string, roleIDs []string) error {
	// 查询角色
	var roles []models.Roles
	if len(roleIDs) > 0 {
		var err error
		roles, err = s.repo.FindRolesByIDs(roleIDs)
		if err != nil {
			return e.ErrFailedToFindRoles.WithScope("SetUserRoles").Wrap(err)
		}
	}

	// 替换用户角色
	if err := s.repo.ReplaceUserRoles(userID, roles); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrUserNotFound.WithScope("SetUserRoles")
		}
		return e.ErrFailedToSetUserRoles.WithScope("SetUserRoles").Wrap(err)
	}

	return nil
}

// AssignRolesToUser 为用户分配角色（追加）
func (s *UserService) AssignRolesToUser(userID string, roleIDs []string) error {
	roles, err := s.repo.FindRolesByIDs(roleIDs)
	if err != nil {
		return e.ErrFailedToFindRoles.WithScope("AssignRolesToUser").Wrap(err)
	}

	if err := s.repo.AppendUserRoles(userID, roles); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrUserNotFound.WithScope("AssignRolesToUser")
		}
		return e.ErrFailedToAssignRoles.WithScope("AssignRolesToUser").Wrap(err)
	}

	return nil
}

// RevokeRolesFromUser 从用户撤销角色
func (s *UserService) RevokeRolesFromUser(userID string, roleIDs []string) error {
	roles, err := s.repo.FindRolesByIDs(roleIDs)
	if err != nil {
		return e.ErrFailedToFindRoles.WithScope("RevokeRolesFromUser").Wrap(err)
	}

	if err := s.repo.DeleteUserRolesByIDs(userID, roles); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrUserNotFound.WithScope("RevokeRolesFromUser")
		}
		return e.ErrFailedToRevokeRoles.WithScope("RevokeRolesFromUser").Wrap(err)
	}

	return nil
}

// GetUserStatusCounts 获取各状态用户数量统计
func (s *UserService) GetUserStatusCounts() (map[models.UserStatus]int64, error) {
	counts, err := s.repo.CountByStatus()
	if err != nil {
		return nil, e.ErrFailedToCountUsersTotal.WithScope("GetUserStatusCounts").Wrap(err)
	}
	return counts, nil
}
