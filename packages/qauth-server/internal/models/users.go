package models

type UserStatus string

const (
	UserStatusActive UserStatus = "ACTIVE"
	UserStatusLocked UserStatus = "LOCKED"
	UserStatusBanned UserStatus = "BANNED"
)

type Users struct {
	BaseModelWithUUID
	Email         string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone         *string    `gorm:"uniqueIndex;size:20" json:"phone,omitempty"`
	PasswordHash  string     `gorm:"size:255;not null" json:"-"`
	Salt          string     `gorm:"size:255" json:"-"`
	Name          string     `gorm:"size:50;not null" json:"name"`
	DisplayName   string     `gorm:"size:50" json:"display_name"`
	AvatarID      *string    `gorm:"type:uuid" json:"avatar_id,omitempty"`
	Status        UserStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`

	// Associations
	Avatar           *Images           `gorm:"foreignKey:AvatarID;references:ID" json:"avatar,omitempty"`
	Roles            []Roles           `gorm:"many2many:users_roles;" json:"roles,omitempty"`
	Organizations    []Organization    `gorm:"foreignKey:UserID" json:"organizations,omitempty"`
	LoginStates      []LoginState      `gorm:"foreignKey:UserID" json:"login_states,omitempty"`
	OperationRecords []OperationRecord `gorm:"foreignKey:OperatorID" json:"operation_records,omitempty"`
	UploadedFiles    []Files           `gorm:"foreignKey:CreatorID;references:ID" json:"uploaded_files,omitempty"`
}

func (Users) TableName() string {
	return "users"
}
