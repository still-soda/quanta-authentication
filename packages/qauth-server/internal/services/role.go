package services

import (
	"qauth-server/internal/models"

	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) GetRoleByID(roleID string) (*RoleService, error) {
	var role RoleService
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) GetUserRole(userID string) (*models.Roles, error) {
	var userRole models.UsersRoles
	if err := s.db.Preload("Role").First(&userRole, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &userRole.Role, nil
}
