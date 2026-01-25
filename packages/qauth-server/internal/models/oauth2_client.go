package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ClientStatus OAuth 客户端状态
type ClientStatus string

const (
	ClientStatusActive      ClientStatus = "active"      // 生产环境
	ClientStatusDevelopment ClientStatus = "development" // 开发中
	ClientStatusDeprecated  ClientStatus = "deprecated"  // 已弃用
)

type ClientData struct {
	ID     string `gorm:"primaryKey" json:"id"`
	Domain string `gorm:"not null" json:"domain"`
	Public bool   `json:"public"`
	Secret string `gorm:"not null" json:"-"`
	UserID string `gorm:"not null" json:"user_id"`
}

// Value implements the driver.Valuer interface for database serialization
func (cd ClientData) Value() (driver.Value, error) {
	return json.Marshal(cd)
}

// Scan implements the sql.Scanner interface for database deserialization
func (cd *ClientData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, cd)
}

// StringArray 字符串数组类型，用于存储 JSON 数组
type StringArray []string

// Value implements the driver.Valuer interface for database serialization
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return "[]", nil
	}
	return json.Marshal(sa)
}

// Scan implements the sql.Scanner interface for database deserialization
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, sa)
}

type OAuth2Client struct {
	BaseModelWithUUID
	Secret       string       `gorm:"size:255;not null" json:"-"`
	Name         string       `gorm:"size:100;not null" json:"name"`
	Description  string       `gorm:"size:500" json:"description"`
	Domain       string       `gorm:"size:500" json:"domain"`
	RedirectURIs StringArray  `gorm:"type:jsonb;default:'[]'" json:"redirect_uris"`
	Scopes       StringArray  `gorm:"type:jsonb;default:'[]'" json:"scopes"`
	GrantTypes   StringArray  `gorm:"type:jsonb;default:'[]'" json:"grant_types"`
	Status       ClientStatus `gorm:"size:20;default:'development'" json:"status"`
	Trusted      bool         `gorm:"default:false" json:"trusted"`
	Logo         string       `gorm:"size:500" json:"logo"`
	Icon         string       `gorm:"size:100;default:'pi pi-box'" json:"icon"`
	IconBg       string       `gorm:"size:200;default:'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)'" json:"icon_bg"`
	LastUsedAt   *time.Time   `json:"last_used_at"`
	RequestCount int64        `gorm:"default:0" json:"request_count"`
	Data         ClientData   `gorm:"type:jsonb" json:"data"`

	// Associations
	LoginStates []LoginState `gorm:"foreignKey:ClientID" json:"login_states,omitempty"`
}

func (OAuth2Client) TableName() string {
	return "oauth2_client"
}

// OAuth2ClientResponse 用于 API 响应的客户端数据结构
type OAuth2ClientResponse struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	ClientID     string       `json:"client_id"` // 同 ID，前端兼容
	Description  string       `json:"description"`
	Domain       string       `json:"domain"`
	RedirectURIs []string     `json:"redirect_uris"`
	Scopes       []string     `json:"scopes"`
	GrantTypes   []string     `json:"grant_types"`
	Status       ClientStatus `json:"status"`
	Trusted      bool         `json:"trusted"`
	Logo         string       `json:"logo"`
	Icon         string       `json:"icon"`
	IconBg       string       `json:"icon_bg"`
	LastUsedAt   *time.Time   `json:"last_used_at"`
	RequestCount int64        `json:"request_count"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// ToResponse 将 OAuth2Client 转换为 API 响应格式
func (c *OAuth2Client) ToResponse() *OAuth2ClientResponse {
	return &OAuth2ClientResponse{
		ID:           c.ID,
		Name:         c.Name,
		ClientID:     c.ID,
		Description:  c.Description,
		Domain:       c.Domain,
		RedirectURIs: c.RedirectURIs,
		Scopes:       c.Scopes,
		GrantTypes:   c.GrantTypes,
		Status:       c.Status,
		Trusted:      c.Trusted,
		Logo:         c.Logo,
		Icon:         c.Icon,
		IconBg:       c.IconBg,
		LastUsedAt:   c.LastUsedAt,
		RequestCount: c.RequestCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
