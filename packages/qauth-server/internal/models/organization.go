package models

type Organization struct {
	BaseModelWithUUID
	UserID       string  `gorm:"type:uuid;not null" json:"user_id"`
	SuperiorID   *string `gorm:"type:uuid" json:"superior_id,omitempty"`
	AncestorPath string  `gorm:"size:500" json:"ancestor_path"`
	Depth        int     `gorm:"default:0" json:"depth"`
	OrgRole      string  `gorm:"size:50" json:"org_role"`
	Class        string  `gorm:"size:50" json:"class"`

	// Associations
	User     Users          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Superior *Organization  `gorm:"foreignKey:SuperiorID;references:ID" json:"superior,omitempty"`
	Children []Organization `gorm:"foreignKey:SuperiorID" json:"children,omitempty"`
}

func (Organization) TableName() string {
	return "organization"
}
