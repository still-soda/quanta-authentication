package models

type ErrorRecord struct {
	ID        string  `gorm:"type:char(36);primaryKey"`
	UserID    string  `gorm:"type:char(36);index"`
	ClientID  *string `gorm:"type:char(36);index"`
	IP        string  `gorm:"type:varchar(45)"`
	Location  string  `gorm:"type:varchar(100)"`
	ErrorType string  `gorm:"type:varchar(50)"`
	Message   string  `gorm:"type:text"`
	Timestamp int64   `gorm:"autoCreateTime"`
}

func (ErrorRecord) TableName() string {
	return "error_records"
}
