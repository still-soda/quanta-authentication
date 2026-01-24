package models

import "time"

type UsersRoles struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UsersID    string    `gorm:"type:uuid;index;uniqueIndex:idx_users_roles,priority:1" json:"users_id"`
	RolesID    string    `gorm:"type:uuid;index;uniqueIndex:idx_users_roles,priority:2" json:"roles_id"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`

	// Associations
	User Users `gorm:"foreignKey:UsersID;references:ID" json:"user"`
	Role Roles `gorm:"foreignKey:RolesID;references:ID" json:"role"`
}

func (UsersRoles) TableName() string {
	return "users_roles"
}
