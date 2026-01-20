package services

import (
	"errors"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"

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
	StudentID string
	Password  string
	Email     string
	Name      string
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(params *CreateUserParams) (*models.Users, error) {
	salt, err := utilities.GenerateSalt(16)
	if err != nil {
		return nil, err
	}

	hashedPassword := utilities.HashPassword(params.Password, salt)

	if err := s.db.Create(&models.Users{
		StudentID:    params.StudentID,
		Email:        params.Email,
		Name:         params.Name,
		Salt:         salt,
		PasswordHash: hashedPassword,
	}).Error; err != nil {
		return nil, err
	}

	var user models.Users
	if err := s.db.First(&user, "student_id = ?", params.StudentID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.Users) error {
	return s.db.Save(user).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID string) error {
	return s.db.Delete(&models.Users{}, "id = ?", userID).Error
}
