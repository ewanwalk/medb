package handbrake

import (
	"context"
	"encoder-backend/pkg/encoder/handbrake/audio"
	"encoding/json"
	"errors"
	"os/exec"
	"regexp"
)

type Scan struct {
	MainFeature int
	TitleList   []Title
}

type Title struct {
	AngleCount        int
	AudioList         []AudioTrack
	ChapterList       []Chapter
	Color             Color
	Container         string
	Crop              []int
	Duration          Duration
	FrameRate         FrameRate
	Geometry          Geometry
	Index             int
	InterlaceDetected bool
	Name              string
	Path              string
	Playlist          int
	SubtitleList      []SubtitleTrack
	Type              int
}

type AudioTrack struct {
	Attributes        AudioAttributes
	BitRate           int
	ChannelCount      int
	ChannelLayout     int
	ChannelLayoutName string
	Codec             int
	CodecName         string
	Description       string
	LFECount          int
	Language          string
	LanguageCode      string
	SampleRate        int
}

type AudioAttributes struct {
	AltCommentary    bool
	Commentary       bool
	Default          bool
	Normal           bool
	Secondary        bool
	VisuallyImpaired bool
}

type SubtitleTrack struct {
	Attributes   SubtitleAttributes
	Format       string
	Language     string
	LanguageCode string
	Source       int
	SourceName   string
}

type SubtitleAttributes struct {
	FourByThree   bool `json:"4By3"`
	Children      bool
	ClosedCaption bool
	Commentary    bool
	Default       bool
	Force         bool
	Large         bool
	Letterbox     bool
	Normal        bool
	PanScan       bool
	Wide          bool
}

type Chapter struct {
	Duration Duration
	Name     string
}

type Duration struct {
	Hours   int
	Minutes int
	Seconds int
	Ticks   int64
}

type FrameRate struct {
	Den int
	Num int
}

type Geometry struct {
	Height int
	Width  int
	Par    FrameRate
}

type Color struct {
	Matrix   int
	Primary  int
	Transfer int
}

var (
	ErrScanFailed    = errors.New("handbrake: failed to parse scan output")
	ErrTitleNotFound = errors.New("handbrake: failed to find title in file")
	jsonRegex        = regexp.MustCompile(`(?ms)(?:JSON Title Set: )({.*?^})`)
)

// scan
// utilizes the handbrake --scan command to obtain detailed information on the file and allows us to make
// better decisions on what to do with a video
func (h *Command) scan(ctx context.Context, file string) (*Title, error) {

	args := []string{
		"--json", "--scan", "-i", file,
	}

	cmd := exec.CommandContext(ctx, h.binary, args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	match := jsonRegex.FindSubmatch(output)
	if len(match) != 2 {
		return nil, ErrScanFailed
	}

	output = match[1]

	scan := &Scan{}

	err = json.Unmarshal(output, scan)
	if err != nil {
		return nil, err
	}

	if len(scan.TitleList) == 0 {
		return nil, ErrTitleNotFound
	}

	return &scan.TitleList[0], nil
}

// IsKnown
// whether or not the codec is a known source or not
func (a AudioTrack) IsKnown() bool {
	for _, codec := range []audio.Codec{
		audio.CodecAAC,
		audio.CodecAC3,
		audio.CodecDTS,
		audio.CodecDTSHD,
		audio.CodecFlac,
		audio.CodecEAC3,
		audio.CodecMP3,
		audio.CodecTrueHD,
	} {
		if audio.Codec(a.Codec) == codec {
			return true
		}
	}

	return false
}

func (a AudioTrack) AudioCodec() audio.Codec {
	return audio.Codec(a.Codec)
}

// MultiLanguageAudio
// whether or not we have more than one language audio track (foreign audio)
func (t Title) MultiLanguageAudio() bool {
	last := ""
	for _, track := range t.AudioList {
		if len(last) != 0 && track.Language != last {
			return true
		}

		last = track.Language
	}

	return false
}
