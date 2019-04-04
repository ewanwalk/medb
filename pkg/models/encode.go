package models

import (
	"github.com/Ewan-Walker/gorm"
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
	ID               int64 `gorm:"AUTO_INCREMENT;primary_key"`
	QualityProfileID int64 `gorm:"type:int(11);not null;index:quality_profile_idx"`
	FileID           int64 `gorm:"type:int(11);not null;index:encodes_media_id_idx"`
	File             File
	TimeStart        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	TimeEnd          time.Time `gorm:"type:timestamp"`
	Duration         int64     `gorm:"type:int(11);not null;default:0"`
	Status           int64     `gorm:"type:int(11);not null;default:0;index:encodes_status_idx"`
	SizeAtStart      int64     `gorm:"type:bigint(20);default:0"`
	SizeAtEnd        int64     `gorm:"type:bigint(20);default:0"`
	ChecksumAtStart  string    `gorm:"type:varchar(255)"`
	ChecksumAtEnd    string    `gorm:"type:varchar(255)"`
	NameAtStart      string    `gorm:"type:varchar(255)"`
	NameAtEnd        string    `gorm:"type:varchar(255)"`
	Error            string    `gorm:"type:varchar(255)"`
}

func (e *Encode) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("TimeEnd", time.Now().UTC())
}
