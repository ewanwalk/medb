package models

import (
	"time"
)

type QualityProfile struct {
	ID   int64  `gorm:"AUTO_INCREMENT;primary_key" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	// Constant quality (0-50) recommended: 21
	QualityLevel int64 `gorm:"type:int(4);not null;default:21" json:"quality_level"`
	// The output container, e.g. mkv, mp4, webm
	VideoContainer string `gorm:"type:varchar(50);not null;default:\"mkv\"" json:"video_container"`
	VideoTune      string `gorm:"type:varchar(50);not null;default:\"\"" json:"video_tune"`
	// The number of subtitle tracks to keep / encode
	SubtitleTracks int64 `gorm:"type:int(2);not null;default:5" json:"subtitle_tracks"`
	// The audio container (--aencoder, copy, acc, etc)
	AudioContainer string `gorm:"type:varchar(512);not null;default:\"copy\"" json:"audio_container"`
	// The audio bitrate (--ab 128)
	AudioBitRate int64 `gorm:"type:int(11);not null;default:128" json:"audio_bit_rate"`
	// The number of audio tracks to keep / encode
	AudioTracks int64 `gorm:"type:int(2);not null;default:5" json:"audio_tracks"`
	// The compression codec to use (e.g. x264, x265
	Codec string `gorm:"type:varchar(50);not null;default:\"x264\"" json:"codec"`
	// The number (or percentage) of threads to use for this profile
	Threads float64 `gorm:"type:double(4,2);not null;default:0.33" json:"threads"`

	CreatedAt time.Time `gorm:"timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"timestamp" json:"updated_at"`

	Paths   []Path   `json:"-"`
	Encodes []Encode `json:"-"`
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
