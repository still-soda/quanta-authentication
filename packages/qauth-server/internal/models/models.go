package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// File 文件模型
type File struct {
	BaseModel
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	FileName    string `gorm:"size:255;not null" json:"file_name"`
	FileKey     string `gorm:"size:255;not null;uniqueIndex" json:"file_key"`
	FileSize    int64  `gorm:"not null" json:"file_size"`
	ContentType string `gorm:"size:100" json:"content_type"`
	URL         string `gorm:"size:500" json:"url"`
	User        User   `gorm:"foreignKey:UserID" json:"user,omitzero"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}
