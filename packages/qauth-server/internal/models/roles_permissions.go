package models

type RolesPermissions struct {
	RoleID       string `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID string `gorm:"type:uuid;primaryKey" json:"permission_id"`

	// Associations
	Role       Roles       `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Permission Permissions `gorm:"foreignKey:PermissionID;references:ID" json:"permission,omitempty"`
}

func (RolesPermissions) TableName() string {
	return "roles_permissions"
}
