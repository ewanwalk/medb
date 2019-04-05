package models

import (
	"errors"
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
	Priority  int64     `gorm:"type:int(4);not null;default:1"`
	Created   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Updated   time.Time `gorm:"type:timestamp"`

	// Option Specific
	QualityProfileID  int64 `gorm:"type:int(11);not null;default:0"`
	QualityProfile    QualityProfile
	EventScanInterval int64 `gorm:"type:int(11);not null;default:500"`
	MinimumFileSize   int64 `gorm:"not null;default:250000000"`
}

func (p *Path) BeforeUpdate(scope *gorm.Scope) error {

	if p.Priority > pathPriorityMax {
		return errors.New("model.models: priority exceeds maximum allowed")
	}

	return scope.SetColumn("Updated", time.Now().UTC())
}

func PathEnabled(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", PathStatusEnabled)
}
