package models

import "time"

type LoginType string

const (
	LoginTypePassword     LoginType = "PASSWORD"
	LoginTypeOAuth2       LoginType = "OAUTH2"
	LoginTypeRefreshToken LoginType = "REFRESH_TOKEN"
)

type LoginState struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string    `gorm:"type:uuid;not null" json:"user_id"`
	ClientID   *string   `gorm:"type:uuid" json:"client_id,omitempty"`
	Time       time.Time `gorm:"autoCreateTime" json:"time"`
	Type       LoginType `gorm:"type:varchar(20);not null" json:"type"`
	IP         string    `gorm:"size:45" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Location   string    `gorm:"size:200" json:"location"`
	IsSuccess  bool      `gorm:"default:false" json:"is_success"`
	FailReason string    `gorm:"size:255" json:"fail_reason"`

	// Associations
	User   Users         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Client *OAuth2Client `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
}

func (LoginState) TableName() string {
	return "login_state"
}
