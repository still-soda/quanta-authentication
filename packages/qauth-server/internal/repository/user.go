package repository

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByID 根据用户ID查询用户
func (r *UserRepository) FindByID(userID string) (*models.Users, error) {
	var user models.Users
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDWithRoles 根据用户ID查询用户（预加载角色）
func (r *UserRepository) FindByIDWithRoles(userID string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByStudentID 根据学生ID查询用户
func (r *UserRepository) FindByStudentID(studentID string) (*models.Users, error) {
	var user models.Users
	if err := r.db.First(&user, "student_id = ?", studentID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *models.Users) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *models.Users) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(userID string) error {
	return r.db.Delete(&models.Users{}, "id = ?", userID).Error
}

// Count 统计用户总数
func (r *UserRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Users{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByRole 统计按角色分类的用户数量
func (r *UserRepository) CountByRole() (map[string]int64, error) {
	type RoleCount struct {
		RoleName string
		Count    int64
	}

	var results []RoleCount
	if err := r.db.Table("users").
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

// CountByStatus 统计各状态用户数量
func (r *UserRepository) CountByStatus() (map[models.UserStatus]int64, error) {
	type StatusCount struct {
		Status models.UserStatus
		Count  int64
	}

	var results []StatusCount
	if err := r.db.Model(&models.Users{}).
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

// ListUsersParams 用户列表查询参数
type ListUsersParams struct {
	Page     int
	PageSize int
	Search   string
	Status   string
	RoleID   string
	SortBy   string
	SortDesc bool
}

// List 分页查询用户列表
func (r *UserRepository) List(params ListUsersParams) ([]models.Users, int64, error) {
	var users []models.Users
	var total int64

	db := r.db.Model(&models.Users{})

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
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
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
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	offset := (params.Page - 1) * params.PageSize
	db = db.Offset(offset).Limit(params.PageSize)

	// 查询用户
	if err := db.Preload("Roles").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeleteUserRoles 删除用户的所有角色关联
func (r *UserRepository) DeleteUserRoles(userID string) error {
	return r.db.Where("users_id = ?", userID).Delete(&models.UsersRoles{}).Error
}

// GetUserRoles 获取用户的所有角色
func (r *UserRepository) GetUserRoles(userID string) ([]models.Roles, error) {
	var user models.Users
	if err := r.db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// ReplaceUserRoles 替换用户的角色
func (r *UserRepository) ReplaceUserRoles(userID string, roles []models.Roles) error {
	var user models.Users
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	return r.db.Model(&user).Association("Roles").Replace(roles)
}

// AppendUserRoles 追加用户角色
func (r *UserRepository) AppendUserRoles(userID string, roles []models.Roles) error {
	var user models.Users
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	return r.db.Model(&user).Association("Roles").Append(roles)
}

// DeleteUserRolesByIDs 删除用户的指定角色
func (r *UserRepository) DeleteUserRolesByIDs(userID string, roles []models.Roles) error {
	var user models.Users
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	return r.db.Model(&user).Association("Roles").Delete(roles)
}

// FindRolesByIDs 根据角色ID批量查询角色
func (r *UserRepository) FindRolesByIDs(roleIDs []string) ([]models.Roles, error) {
	var roles []models.Roles
	if err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
