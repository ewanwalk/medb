package audio

import (
	"fmt"
	"strings"
)

const (
	CodecUnknown Codec = -1
	CodecCopy    Codec = -2
	CodecAll     Codec = -3
	CodecAAC     Codec = 65536
	CodecAC3     Codec = 2048
	CodecDTS     Codec = 8192
	CodecDTSHD   Codec = 262144
	CodecFlac    Codec = 1048576
	CodecEAC3    Codec = 16777216
	CodecMP3     Codec = 524288
	CodecTrueHD  Codec = 33554432
)

type Codec int

// String
// the string representation of an audio codec as known to handbrake
func (a Codec) String() string {

	switch a {
	case CodecAll:
		return "all"
	case CodecCopy:
		return "copy"
	case CodecAAC:
		return "aac"
	case CodecAC3:
		return "ac3"
	case CodecDTS:
		return "dts"
	case CodecDTSHD:
		return "dtshd"
	case CodecFlac:
		return "flac"
	case CodecEAC3:
		return "eac3"
	case CodecMP3:
		return "mp3"
	case CodecTrueHD:
		return "truehd"
	}

	return ""
}

func (a Codec) QualityRank() int {

	switch a {
	case CodecMP3:
		return 1
	case CodecAAC:
		return 2
	case CodecAC3:
		return 3
	case CodecEAC3:
		return 4
	case CodecDTS:
		return 5
	case CodecDTSHD:
		return 6
	case CodecTrueHD:
		return 7
	case CodecFlac:
		return 8
	}

	return 0
}

// CopyString
// the string used to `copy` a source as known to handbrake
func (a Codec) CopyString() string {
	return fmt.Sprintf("copy:%s", a.String())
}

// AudioCodecID
// obtain an audio codec id from a string name
func CodecFromName(name string) Codec {

	switch strings.ToLower(name) {
	case "all":
		return CodecAll
	case "copy":
		return CodecCopy
	case "aac":
		return CodecAAC
	case "ac3":
		return CodecAC3
	case "dts":
		return CodecDTS
	case "dtshd":
		return CodecDTSHD
	case "flac":
		return CodecFlac
	case "eac3":
		return CodecEAC3
	case "mp3":
		return CodecMP3
	case "truehd":
		return CodecTrueHD
	}

	return CodecUnknown
}
