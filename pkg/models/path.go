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
	ID        int64     `gorm:"AUTO_INCREMENT;primary_key;"`
	Name      string    `gorm:"type:varchar(255)"`
	Directory string    `gorm:"type:varchar(1024);not null"`
	Type      int64     `gorm:"type:int(11);not null;default:1"`
	Status    int64     `gorm:"type:int(3);not null;default:1"`
	Priority  int64     `gorm:"type:int(11);not null;default:1"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp"`

	// Option Specific
	QualityProfileID  int64 `gorm:"type:int(11);not null;index"`
	QualityProfile    QualityProfile
	EventScanInterval int64 `gorm:"type:int(11);not null;default:500"`
	MinimumFileSize   int64 `gorm:"not null;default:250000000"`
}

func PathEnabled(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", PathStatusEnabled)
}
