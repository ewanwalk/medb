package models

import "time"

const (
	RevisionReasonReplaced = 0
)

type Revision struct {
	ID        int64      `gorm:"AUTO_INCREMENT;primary_key" json:"id"`
	FileID    int64      `gorm:"type:int(11);not null" json:"file_id"`
	PathID    int64      `gorm:"type:int(11);not null" json:"path_id"`
	Checksum  string     `gorm:"type:varchar(255);not null" json:"checksum"`
	Size      int64      `gorm:"type:bigint(20);not null;default:0" json:"size"`
	Encoded   int64      `gorm:"type:int(2);default:0" json:"encoded"`
	Reason    int64      `gorm:"type:int(2);default:0" json:"reason"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
