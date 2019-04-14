package models

import (
	"encoder-backend/pkg/slice"
	"runtime"
	"time"
)

type QualityProfile struct {
	ID   int64  `gorm:"AUTO_INCREMENT;primary_key" json:"id,omitempty"`
	Name string `gorm:"type:varchar(255);not null" json:"name,omitempty"`

	// Constant quality (0-50) recommended: 21
	QualityLevel int64 `gorm:"type:int(4);not null;default:21" json:"quality_level,omitempty"`
	// The output container, e.g. mkv, mp4, webm
	VideoContainer string `gorm:"type:varchar(50);not null;default:\"mkv\"" json:"video_container,omitempty"`
	VideoTune      string `gorm:"type:varchar(50);not null;default:\"\"" json:"video_tune,omitempty"`
	// The number of subtitle tracks to keep / encode
	SubtitleTracks int64 `gorm:"type:int(2);not null;default:5" json:"subtitle_tracks,omitempty"`
	// The audio container (--aencoder, copy, acc, etc)
	AudioContainer string `gorm:"type:varchar(512);not null;default:\"copy\"" json:"audio_container,omitempty"`
	// The audio bitrate (--ab 128)
	AudioBitRate int64 `gorm:"type:int(11);not null;default:128" json:"audio_bit_rate,omitempty"`
	// The number of audio tracks to keep / encode
	AudioTracks int64 `gorm:"type:int(2);not null;default:5" json:"audio_tracks,omitempty"`
	// The compression codec to use (e.g. x264, x265
	Codec string `gorm:"type:varchar(50);not null;default:\"x264\"" json:"codec,omitempty"`
	// The number (or percentage) of threads to use for this profile
	Threads float64 `gorm:"type:double(4,2);not null;default:0.33" json:"threads,omitempty"`

	CreatedAt *time.Time `gorm:"timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"timestamp" json:"updated_at,omitempty"`

	Paths   []Path   `json:"paths,omitempty"`
	Encodes []Encode `json:"encodes,omitempty"`
}

func (q *QualityProfile) IsValid() map[string]string {

	errs := make(map[string]string)

	if q.SubtitleTracks < 0 {
		q.SubtitleTracks = 0
	}

	if q.AudioTracks < 0 {
		q.AudioTracks = 1
	}

	if q.QualityLevel != 0 {
		if q.QualityLevel < 1 {
			errs["quality_level"] = "The quality level field must be greater than or equal to 1"
		}

		if q.QualityLevel > 50 {
			errs["quality_level"] = "The quality level fields must be less than or equal to 50"
		}
	}

	if len(q.Codec) != 0 && !slice.InString(q.Codec, []string{"x264", "x264_10bit", "x265", "x265_10bit", "x265_12bit", "VP8", "VP9"}) {
		errs["codec"] = "The video codec provided is invalid"
	}

	// containers
	if len(q.VideoContainer) != 0 && !slice.InString(q.VideoContainer, []string{"mp4", "mkv", "webm"}) {
		errs["video_container"] = "The video container provided is invalid"
	}

	// video tune
	if len(q.VideoTune) != 0 {

		if q.Codec != "x264" {
			errs["video_tune"] = "The video tune field may only be used with the x264 codec"
		} else if !slice.InString(q.VideoTune, []string{"animation", "film", "grain"}) {
			errs["video_tune"] = "The video tune field provided is invalid"
		}

	}

	// TODO validate audio container
	if q.AudioBitRate != 0 && q.AudioBitRate < 68 {
		errs["audio_bit_rate"] = "The audio bit rate field must be greater than 68"
	}

	// if someone provides a high cpu usage
	available := float64(runtime.NumCPU()) - q.Threads*float64(runtime.NumCPU())
	if q.Threads >= float64(runtime.NumCPU()) || available < 1 {
		errs["threads"] = "The threads field exceeds or uses 100% of cores available, one core must be left available"
	} else if q.Threads < 0 {
		errs["threads"] = "The threads field must be above 0"
	}

	return errs
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
