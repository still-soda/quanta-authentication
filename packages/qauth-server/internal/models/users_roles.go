package models

import "time"

type UsersRoles struct {
	UsersID    string    `gorm:"type:uuid;primaryKey" json:"users_id"`
	RolesID    string    `gorm:"type:uuid;primaryKey" json:"roles_id"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`

	// Associations
	User Users `gorm:"foreignKey:UsersID;references:ID" json:"user,omitempty"`
	Role Roles `gorm:"foreignKey:RolesID;references:ID" json:"role,omitempty"`
}

func (UsersRoles) TableName() string {
	return "users_roles"
}
