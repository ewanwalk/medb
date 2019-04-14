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
	ID               int64           `gorm:"AUTO_INCREMENT;primary_key" json:"id,omitempty"`
	QualityProfileID int64           `gorm:"type:int(11);not null;index" json:"quality_profile_id,omitempty"`
	QualityProfile   *QualityProfile `json:"quality_profile,omitempty"`
	FileID           int64           `gorm:"type:int(11);not null;index" json:"file_id,omitempty"`
	File             *File           `json:"file,omitempty"`
	TimeStart        *time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"time_start,omitempty"`
	TimeEnd          *time.Time      `gorm:"type:timestamp" json:"time_end,omitempty"`
	Duration         int64           `gorm:"type:int(11);not null;default:0" json:"duration,omitempty"`
	Status           int64           `gorm:"type:int(11);not null;default:0;index" json:"status,omitempty"`
	SizeAtStart      int64           `gorm:"type:bigint(20);default:0" json:"size_at_start,omitempty"`
	SizeAtEnd        int64           `gorm:"type:bigint(20);default:0" json:"size_at_end,omitempty"`
	ChecksumAtStart  string          `gorm:"type:varchar(255)" json:"checksum_at_start,omitempty"`
	ChecksumAtEnd    string          `gorm:"type:varchar(255)" json:"checksum_at_end,omitempty"`
	NameAtStart      string          `gorm:"type:varchar(255)" json:"name_at_start,omitempty"`
	NameAtEnd        string          `gorm:"type:varchar(255)" json:"name_at_end,omitempty"`
	Error            string          `gorm:"type:varchar(255)" json:"error,omitempty"`
	CreatedAt        *time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
}
