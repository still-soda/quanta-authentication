package models

type Counters struct {
	BaseModelWithID
	Key       string `gorm:"type:varchar(100);index;not null" json:"key"`
	Count     int64  `gorm:"type:bigint;not null" json:"count"`
	Timestamp int64  `gorm:"autoCreateTime" json:"timestamp"`
}

func (Counters) TableName() string {
	return "counters"
}
