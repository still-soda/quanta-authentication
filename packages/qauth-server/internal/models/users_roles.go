package models

import "time"

type UsersRoles struct {
	UserID     string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	RoleID     string    `gorm:"type:uuid;primaryKey" json:"role_id"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`

	// Associations
	User Users `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Role Roles `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

func (UsersRoles) TableName() string {
	return "users_roles"
}
