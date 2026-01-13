package models

type Files struct {
	BaseModelWithUUID
	StorageKey  string  `gorm:"size:500;not null" json:"storage_key"`
	Bucket      string  `gorm:"size:100;not null" json:"bucket"`
	MimeType    string  `gorm:"size:100" json:"mime_type"`
	SizeBytes   int64   `gorm:"default:0" json:"size_bytes"`
	CreatorID   *string `gorm:"type:uuid;index" json:"creator_id,omitempty"`
	IsTemporary bool    `gorm:"default:false" json:"is_temporary"`

	// Associations
	Creator *Users `gorm:"foreignKey:CreatorID;references:ID" json:"creator,omitempty"`
}

func (Files) TableName() string {
	return "files"
}
