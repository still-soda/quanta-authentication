package models

type Images struct {
	BaseModelWithUUID
	Width     int     `gorm:"default:0" json:"width"`
	Height    int     `gorm:"default:0" json:"height"`
	FileID    *string `gorm:"type:uuid;index" json:"file_id,omitempty"`
	CreatorID *string `gorm:"type:uuid;index" json:"creator_id,omitempty"`

	// Associations
	File    *Files `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
	Creator *Users `gorm:"foreignKey:CreatorID;references:ID" json:"creator,omitempty"`
}

func (Images) TableName() string {
	return "images"
}
