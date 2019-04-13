package models

import (
	"encoder-backend/pkg/file"
	"github.com/Ewan-Walker/gorm"
	"os"
	"path/filepath"
	"time"
)

// TODO [integrity] on boot clear any "pending" or "running" encode status files

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
	ID            int64     `gorm:"AUTO_INCREMENT;primary_key" json:"id"`
	PathID        int64     `gorm:"type:int(11);not null;index" json:"path_id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Size          int64     `gorm:"type:bigint(20);not null;default:0" json:"size"`
	Checksum      string    `gorm:"type:varchar(255);not null" json:"checksum"`
	Source        string    `gorm:"type:varchar(512);not null" json:"source"`
	Status        int64     `gorm:"type:int(2);default:1;index" json:"status"`
	StatusEncoder int64     `gorm:"type:int(2);default:0" json:"status_encoder"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	// TODO potentially allow marshalling but will need a custom marshaller as we dont want this with a request for
	// worker status
	Encodes []Encode `json:"-"`
	Path    Path     `gorm:"association_autoupdate:false" json:"-"`
}

// CurrentChecksum
// obtains the raw checksum by checking the file directly
func (f File) CurrentChecksum() (string, error) {
	return file.Checksum(filepath.Join(f.Source, f.Name))
}

// Exists
// checks the filesystem to ensure the file still exists
func (f File) Exists() bool {

	_, err := os.Stat(filepath.Join(f.Source, f.Name))
	if err != nil {
		return false
	}

	sum, err := f.CurrentChecksum()
	if err != nil || sum != f.Checksum {
		return false
	}

	return true
}

// FileNeedsEncode
// a gorm scope to find files which still need encoding
func FileNeedsEncode(db *gorm.DB) *gorm.DB {

	return db.Joins("left join paths on paths.id = files.path_id").
		Joins("left join encodes as e on e.file_id = files.id").
		Where("e.id = (?) OR e.id is null",
			db.Table("encodes").Select("MAX(id)").Where("file_id = files.id").QueryExpr(),
		).
		Where("files.size >= paths.minimum_file_size").
		Where("files.status = ?", FileStatusEnabled).
		Where(
			"((files.status_encoder = ? AND files.checksum != e.checksum_at_end) OR e.status = ?)",
			FileEncodeStatusNotDone, EncodeCancelled,
		).
		Where("files.status_encoder <> ?", FileEncodeStatusErrored).
		Order("paths.priority DESC")
}
