package models

type Permissions struct {
	BaseModelWithUUID
	Resource    string `gorm:"size:50;not null" json:"resource"`
	Action      string `gorm:"size:50;not null" json:"action"`
	Code        string `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Description string `gorm:"size:255" json:"description"`

	// Associations
	Roles []Roles `gorm:"many2many:roles_permissions;" json:"roles,omitempty"`
}

func (Permissions) TableName() string {
	return "permissions"
}
