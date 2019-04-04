package models

import (
	"encoder-backend/pkg/file"
	"github.com/Ewan-Walker/gorm"
	"path/filepath"
	"time"
)

const (
	// Status codes
	FileStatusDeleted int64 = 0
	FileStatusEnabled int64 = 1
	// Encoder status codes
	FileEncodeStatusNotDone int64 = 0
	FileEncodeStatusDone    int64 = 1
	FileEncodeStatusPending int64 = 2
	FileEncodeStatusErrored int64 = 3
	FileEncodeStatusRunning int64 = 10
)

type File struct {
	ID            int64     `gorm:"AUTO_INCREMENT;primary_key"`
	PathID        int64     `gorm:"type:int(11);not null"`
	Name          string    `gorm:"type:varchar(255);not null"`
	Size          int64     `gorm:"type:bigint(20);not null;default:0"`
	Checksum      string    `gorm:"type:varchar(255);not null"`
	Source        string    `gorm:"type:varchar(512);not null"`
	Status        int64     `gorm:"type:int(2);default:1;index:media_status_idx"`
	StatusEncoder int64     `gorm:"type:int(2);default:0"`
	Created       time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Updated       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Encodes []Encode
}

func (f *File) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("Updated", time.Now().UTC())
}

// CurrentChecksum
// obtains the raw checksum by checking the file directly
func (f File) CurrentChecksum() (string, error) {
	return file.Checksum(filepath.Join(f.Source, f.Name))
}
