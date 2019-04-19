package models

import (
	"github.com/ewanwalk/gorm"
	"os"
	"path/filepath"
	"time"
)

const (
	PathStatusDisabled = iota
	PathStatusEnabled
)

const (
	pathPriorityMax     = 999
	PathAbsoluteMinimum = 2 << 25 // ~67mb
)

type Path struct {
	ID        int64      `gorm:"AUTO_INCREMENT;primary_key;" json:"id,omitempty"`
	Name      string     `gorm:"type:varchar(255)" json:"name,omitempty"`
	Directory string     `gorm:"type:varchar(1024);not null" json:"directory,omitempty"`
	Type      int64      `gorm:"type:int(11);not null;default:1" json:"type,omitempty"`
	Status    int64      `gorm:"type:int(3);not null;default:1" json:"status,omitempty"`
	Priority  int64      `gorm:"type:int(11);not null;default:1" json:"priority,omitempty"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"type:timestamp" json:"updated_at,omitempty"`

	// Option Specific
	QualityProfileID  int64           `gorm:"type:int(11);not null;index" json:"quality_profile_id,omitempty"`
	QualityProfile    *QualityProfile `json:"quality_profile,omitempty"`
	EventScanInterval int64           `gorm:"type:int(11);not null;default:500" json:"event_scan_interval,omitempty"`
	MinimumFileSize   int64           `gorm:"not null;default:250000000" json:"minimum_file_size,omitempty"`
}

func (p *Path) IsValid() map[string]string {

	errs := make(map[string]string)

	if len(p.Name) == 0 {
		errs["name"] = "The name field is required"
	}

	if len(p.Directory) == 0 {
		errs["directory"] = "The directory field is required"
	} else {

		dir, err := filepath.Abs(p.Directory)
		if err != nil {
			errs["directory"] = err.Error()
		} else {

			p.Directory = dir

			info, err := os.Stat(p.Directory)
			if err != nil {
				errs["directory"] = err.Error()
			} else if !info.IsDir() {
				errs["directory"] = "The directory provided is not a directory"
			}
		}
	}

	if p.EventScanInterval != 0 && p.EventScanInterval < 250 {
		errs["event_scan_interval"] = "The event scan interval minimum is 250ms"
	}

	if p.Status != PathStatusEnabled && p.Status != PathStatusDisabled {
		errs["status"] = "The status provided is invalid"
	}

	if p.Priority != 0 && p.Priority < 0 {
		errs["priority"] = "The priority minimum is 1"
	}

	// roughly 67mb
	if p.MinimumFileSize != 0 && p.MinimumFileSize <= PathAbsoluteMinimum {
		errs["minimum_file_size"] = "The minimum file"
	}

	if p.QualityProfileID == 0 {
		errs["quality_profile"] = "The quality profile field is required"
	}

	return errs
}

// Scopes

func PathEnabled(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", PathStatusEnabled)
}
