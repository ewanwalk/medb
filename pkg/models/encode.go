package models

import (
	"time"
)

const (
	EncodeErrored   int64 = 0
	EncodeDone      int64 = 1
	EncodeCancelled int64 = 2
	EncodeRunning   int64 = 99
	EncodeQueued    int64 = 98
)

type Encode struct {
	ID               int64     `gorm:"AUTO_INCREMENT;primary_key" json:"id"`
	QualityProfileID int64     `gorm:"type:int(11);not null;index" json:"quality_profile_id"`
	FileID           int64     `gorm:"type:int(11);not null;index" json:"file_id"`
	File             File      `json:"-"`
	TimeStart        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"time_start"`
	TimeEnd          time.Time `gorm:"type:timestamp" json:"time_end"`
	Duration         int64     `gorm:"type:int(11);not null;default:0" json:"duration"`
	Status           int64     `gorm:"type:int(11);not null;default:0;index" json:"status"`
	SizeAtStart      int64     `gorm:"type:bigint(20);default:0" json:"size_at_start"`
	SizeAtEnd        int64     `gorm:"type:bigint(20);default:0" json:"size_at_end"`
	ChecksumAtStart  string    `gorm:"type:varchar(255)" json:"checksum_at_start"`
	ChecksumAtEnd    string    `gorm:"type:varchar(255)" json:"checksum_at_end"`
	NameAtStart      string    `gorm:"type:varchar(255)" json:"name_at_start"`
	NameAtEnd        string    `gorm:"type:varchar(255)" json:"name_at_end"`
	Error            string    `gorm:"type:varchar(255)" json:"error"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
