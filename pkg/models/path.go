package models

import (
	"github.com/Ewan-Walker/gorm"
	"time"
)

const (
	PathStatusDisabled = iota
	PathStatusEnabled
)

const (
	pathPriorityMax = 999
)

type Path struct {
	ID        int64     `gorm:"AUTO_INCREMENT;primary_key;" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Directory string    `gorm:"type:varchar(1024);not null" json:"directory"`
	Type      int64     `gorm:"type:int(11);not null;default:1" json:"type"`
	Status    int64     `gorm:"type:int(3);not null;default:1" json:"status"`
	Priority  int64     `gorm:"type:int(11);not null;default:1" json:"priority"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`

	// Option Specific
	QualityProfileID  int64          `gorm:"type:int(11);not null;index" json:"quality_profile_id"`
	QualityProfile    QualityProfile `json:"-"`
	EventScanInterval int64          `gorm:"type:int(11);not null;default:500" json:"event_scan_interval"`
	MinimumFileSize   int64          `gorm:"not null;default:250000000" json:"minimum_file_size"`
}

func PathEnabled(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", PathStatusEnabled)
}
