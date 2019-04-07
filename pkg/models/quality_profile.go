package models

import (
	"time"
)

type QualityProfile struct {
	ID           int64  `gorm:"AUTO_INCREMENT;primary_key"`
	Name         string `gorm:"type:varchar(255);not null"`
	QualityLevel int64  `gorm:"type:int(4);not null;default:21"`
	// The output container, e.g. mkv, mp4
	VideoContainer string `gorm:"type:varchar(50);not null;default:\"mkv\""`
	// The number of subtitle tracks to keep / encode
	SubtitleTracks int64 `gorm:"type:int(2);not null;default:5"`
	// The audio container (--aencoder, copy, acc, etc)
	AudioContainer string `gorm:"type:varchar(50);not null;default:\"copy\""`
	// The number of audio tracks to keep / encode
	AudioTracks int64 `gorm:"type:int(2);not null;default:5"`
	// The compression codec to use (e.g. x264, x265
	Codec string `gorm:"type:varchar(50);not null;default:\"x264\""`
	// The number (or percentage) of threads to use for this profile
	Threads   float64   `gorm:"type:double(4,2);not null;default:0.33"`
	CreatedAt time.Time `gorm:"timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"timestamp"`

	Paths   []Path
	Encodes []Encode
}
