package services

import (
	"errors"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUserByID 根据用户ID获取用户信息
func (s *UserService) GetUserByID(userID string, withRole bool) (*models.Users, error) {
	var user models.Users
	db := s.db

	if withRole {
		db = db.Preload("Roles")
	}

	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByStudentID 根据学生ID获取用户信息
func (s *UserService) GetUserByStudentID(studentID string) (*models.Users, error) {
	var user models.Users
	if err := s.db.First(&user, "student_id = ?", studentID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AuthenticateUser 验证用户凭据
func (s *UserService) AuthenticateUser(studentID, password string) (*models.Users, error) {
	var user models.Users
	if err := s.db.First(&user, "student_id = ?", studentID).Error; err != nil {
		return nil, err
	}

	verified := utilities.VerifyPassword(password, user.Salt, user.PasswordHash)
	if !verified {
		return nil, errors.New("authentication failed")
	}

	return &user, nil
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
		return nil, err
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

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// 重新查询以获取完整数据
	var createdUser models.Users
	if err := s.db.Preload("Roles").First(&createdUser, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.Users) error {
	return s.db.Save(user).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID string) error {
	// 先删除用户角色关联
	if err := s.db.Where("users_id = ?", userID).Delete(&models.UsersRoles{}).Error; err != nil {
		return err
	}
	return s.db.Delete(&models.Users{}, "id = ?", userID).Error
}

// UserCount 返回所有用户的数量
func (s *UserService) UserCount() (int64, error) {
	var count int64
	if err := s.db.Model(&models.Users{}).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}

// GetUserCountByRole 返回按角色分类的用户数量统计
func (s *UserService) GetUserCountByRole() (map[string]int64, error) {
	type RoleCount struct {
		RoleName string
		Count    int64
	}

	var results []RoleCount
	if err := s.db.Table("users").
		Joins("LEFT JOIN users_roles ON users.id = users_roles.users_id").
		Joins("LEFT JOIN roles ON users_roles.roles_id = roles.id").
		Where("users.deleted_at IS NULL AND roles.deleted_at IS NULL").
		Select("COALESCE(roles.name, '未分配角色') as role_name, COUNT(DISTINCT users.id) as count").
		Group("roles.name").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	roleCountMap := make(map[string]int64)
	for _, result := range results {
		roleCountMap[result.RoleName] = result.Count
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

	db := s.db.Model(&models.Users{})

	// 搜索过滤
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ? OR student_id ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 状态过滤
	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	// 角色过滤
	if params.RoleID != "" {
		db = db.Joins("JOIN users_roles ON users.id = users_roles.users_id").
			Where("users_roles.roles_id = ?", params.RoleID)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	sortColumn := "created_at"
	if params.SortBy != "" {
		allowedSorts := map[string]bool{
			"name": true, "email": true, "created_at": true, "updated_at": true, "status": true,
		}
		if allowedSorts[params.SortBy] {
			sortColumn = params.SortBy
		}
	}
	sortOrder := "DESC"
	if !params.SortDesc {
		sortOrder = "ASC"
	}
	db = db.Order(sortColumn + " " + sortOrder)

	// 分页
	offset := (params.Page - 1) * params.PageSize
	db = db.Offset(offset).Limit(params.PageSize)

	// 查询用户
	var users []models.Users
	if err := db.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}

	// 转换为带统计信息的结果
	usersWithStats := make([]UserWithStats, len(users))
	for i, user := range users {
		roleNames := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roleNames[j] = role.Name
		}

		// 获取最后登录时间（从审计日志）
		var lastLogin *time.Time
		var auditLog models.AuditLog
		if err := s.db.Where("operator_id = ? AND action = ?", user.ID, models.AuditActionLogin).
			Order("created_at DESC").
			First(&auditLog).Error; err == nil {
			lastLogin = &auditLog.CreatedAt
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
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
	var user models.Users
	if err := s.db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// 获取最后登录时间
	var lastLogin *time.Time
	var auditLog models.AuditLog
	if err := s.db.Where("operator_id = ? AND action = ?", user.ID, models.AuditActionLogin).
		Order("created_at DESC").
		First(&auditLog).Error; err == nil {
		lastLogin = &auditLog.CreatedAt
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &UserWithStats{
		Users:       user,
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
	var user models.Users
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
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

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// 重新查询以获取完整数据
	if err := s.db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(userID string, newPassword string) error {
	var user models.Users
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	salt, err := utilities.GenerateSalt(16)
	if err != nil {
		return err
	}

	hashedPassword := utilities.HashPassword(newPassword, salt)
	user.Salt = salt
	user.PasswordHash = hashedPassword

	return s.db.Save(&user).Error
}

// GetUserRoles 获取用户的所有角色
func (s *UserService) GetUserRoles(userID string) ([]models.Roles, error) {
	var user models.Users
	if err := s.db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// SetUserRoles 设置用户的角色（替换现有角色）
func (s *UserService) SetUserRoles(userID string, roleIDs []string) error {
	var user models.Users
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	// 查询角色
	var roles []models.Roles
	if len(roleIDs) > 0 {
		if err := s.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}

	// 替换用户角色
	if err := s.db.Model(&user).Association("Roles").Replace(roles); err != nil {
		return err
	}

	return nil
}

// AssignRolesToUser 为用户分配角色（追加）
func (s *UserService) AssignRolesToUser(userID string, roleIDs []string) error {
	var user models.Users
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var roles []models.Roles
	if err := s.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}

	if err := s.db.Model(&user).Association("Roles").Append(roles); err != nil {
		return err
	}

	return nil
}

// RevokeRolesFromUser 从用户撤销角色
func (s *UserService) RevokeRolesFromUser(userID string, roleIDs []string) error {
	var user models.Users
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var roles []models.Roles
	if err := s.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}

	if err := s.db.Model(&user).Association("Roles").Delete(roles); err != nil {
		return err
	}

	return nil
}

// GetUserStatusCounts 获取各状态用户数量统计
func (s *UserService) GetUserStatusCounts() (map[models.UserStatus]int64, error) {
	type StatusCount struct {
		Status models.UserStatus
		Count  int64
	}

	var results []StatusCount
	if err := s.db.Model(&models.Users{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[models.UserStatus]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}
