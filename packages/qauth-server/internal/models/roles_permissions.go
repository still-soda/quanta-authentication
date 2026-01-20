package models

type RolesPermissions struct {
	RolesID       string `gorm:"type:uuid;primaryKey" json:"roles_id"`
	PermissionsID string `gorm:"type:uuid;primaryKey" json:"permissions_id"`

	// Associations
	Role       Roles       `gorm:"foreignKey:RolesID;references:ID" json:"role,omitempty"`
	Permission Permissions `gorm:"foreignKey:PermissionsID;references:ID" json:"permission,omitempty"`
}

func (RolesPermissions) TableName() string {
	return "roles_permissions"
}
