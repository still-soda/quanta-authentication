package models

type Roles struct {
	BaseModelWithUUID
	Resource    string `gorm:"size:50;not null" json:"resource"`
	Code        string `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name        string `gorm:"size:50;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsSystem    bool   `gorm:"default:false" json:"is_system"`

	// Associations
	Users       []Users       `gorm:"many2many:users_roles;" json:"users,omitempty"`
	Permissions []Permissions `gorm:"many2many:roles_permissions;" json:"permissions,omitempty"`
}

func (Roles) TableName() string {
	return "roles"
}
