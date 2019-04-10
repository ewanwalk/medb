package handbrake

import "runtime"

func (h *Command) defaults() {

	if len(h.profile.VideoContainer) == 0 || !inSliceString([]string{"mkv", "mp4"}, h.profile.VideoContainer) {
		h.profile.VideoContainer = "mkv"
	}

	if len(h.profile.Codec) == 0 || !inSliceString([]string{"x264", "x264_10bit", "x265", "x265_10bit", "VP8", "VP9"}, h.profile.Codec) {
		h.profile.Codec = "x264"
	}

	if h.profile.QualityLevel <= 0 || h.profile.QualityLevel > 50 {
		h.profile.QualityLevel = 21
	}

	if len(h.profile.VideoTune) != 0 {

		switch h.profile.Codec {
		case "x264":
			fallthrough
		case "x264_10bit":
			if !inSliceString([]string{"film", "animation", "grain", "stillimage", "psnr", "ssim", "fastdecode", "zerolatency"}, h.profile.VideoTune) {
				h.profile.VideoTune = ""
			}
		case "x265":
			fallthrough
		case "x265_10bit":
			if !inSliceString([]string{"psnr", "ssim", "fastdecode", "zerolatency"}, h.profile.VideoTune) {
				h.profile.VideoTune = ""
			}
		default:
			h.profile.VideoTune = ""
		}

	}

	if len(h.profile.AudioContainer) == 0 {
		h.profile.AudioContainer = "copy"
	}

	if h.profile.AudioBitRate == 0 {
		h.profile.AudioBitRate = 128
	}

	if h.profile.Threads == 0 {
		h.profile.Threads = 1.0
	}

	if h.profile.Threads < 1 {
		h.profile.Threads = float64(runtime.NumCPU()) * h.profile.Threads
	}

}
