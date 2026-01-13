package models

import (
	"time"

	"gorm.io/datatypes"
)

type OperationRecord struct {
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OperatorID string         `gorm:"type:uuid;not null" json:"operator_id"`
	Module     string         `gorm:"size:50;not null" json:"module"`
	Action     string         `gorm:"size:50;not null" json:"action"`
	TargetID   string         `gorm:"size:255" json:"target_id"`
	Detail     datatypes.JSON `gorm:"type:jsonb" json:"detail"`
	IP         string         `gorm:"size:45" json:"ip"`
	Time       time.Time      `gorm:"autoCreateTime" json:"time"`
	DurationMs int            `gorm:"default:0" json:"duration_ms"`

	// Associations
	Operator Users `gorm:"foreignKey:OperatorID;references:ID" json:"operator,omitempty"`
}

func (OperationRecord) TableName() string {
	return "operation_record"
}
