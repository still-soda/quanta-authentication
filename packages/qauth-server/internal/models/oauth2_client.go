package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type OAuth2Client struct {
	BaseModelWithUUID
	Secret string     `gorm:"size:255;not null" json:"-"`
	Name   string     `gorm:"size:100;not null" json:"name"`
	Domain string     `gorm:"size:500" json:"domain"`
	Data   ClientData `gorm:"type:jsonb" json:"data"`

	// Associations
	LoginStates []LoginState `gorm:"foreignKey:ClientID" json:"login_states,omitempty"`
}

func (OAuth2Client) TableName() string {
	return "oauth2_client"
}
