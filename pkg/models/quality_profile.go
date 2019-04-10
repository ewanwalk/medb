package models

import (
	"time"
)

type QualityProfile struct {
	ID   int64  `gorm:"AUTO_INCREMENT;primary_key"`
	Name string `gorm:"type:varchar(255);not null"`

	// Constant quality (0-50) recommended: 21
	QualityLevel int64 `gorm:"type:int(4);not null;default:21"`
	// The output container, e.g. mkv, mp4, webm
	VideoContainer string `gorm:"type:varchar(50);not null;default:\"mkv\""`
	VideoTune      string `gorm:"type:varchar(50);not null;default:\"\""`
	// The number of subtitle tracks to keep / encode
	SubtitleTracks int64 `gorm:"type:int(2);not null;default:5"`
	// The audio container (--aencoder, copy, acc, etc)
	AudioContainer string `gorm:"type:varchar(512);not null;default:\"copy\""`
	// The audio bitrate (--ab 128)
	AudioBitRate int64 `gorm:"type:int(11);not null;default:128"`
	// The number of audio tracks to keep / encode
	AudioTracks int64 `gorm:"type:int(2);not null;default:5"`
	// The compression codec to use (e.g. x264, x265
	Codec string `gorm:"type:varchar(50);not null;default:\"x264\""`
	// The number (or percentage) of threads to use for this profile
	Threads float64 `gorm:"type:double(4,2);not null;default:0.33"`

	CreatedAt time.Time `gorm:"timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"timestamp"`

	Paths   []Path
	Encodes []Encode
}

// TODO default quality profiles

/**
 x264 video tune options
film
animation
*/

// Movie:
/*
AudioContainer ac3
AudioBitRate   640
VideoTune      film (dependant on if animated ?)
*/
// Show:
/*
AudioContainer aac
AudioBitRate   128
*/
