package handbrake

import (
	"context"
	"encoder-backend/pkg/config"
	"encoder-backend/pkg/encoder/handbrake/audio"
	"encoder-backend/pkg/models"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	binary  string
	profile models.QualityProfile
	title   *Title
	args    []string

	// reference
	staged string
}

var (
	ErrStagingDirectory = errors.New("handbrake.run: failed to create or find staging path")
)

func New(profile models.QualityProfile) (*Command, error) {

	// set default arguments
	e := &Command{
		profile: profile,
		args: []string{
			"-v0",
		},
	}

	e.defaults()
	e.parse()

	return e, e.validate()
}

func (h *Command) StagedFile() string {
	return h.staged
}

// Get
// obtain the command to run this handbrake instance
func (h *Command) Get(ctx context.Context, file string) (*exec.Cmd, error) {

	err := h.get(ctx, file)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"binary":  h.binary,
		"staging": h.staged,
		"args":    strings.Join(h.args, " "),
	}).Debug("handbrake.get: exec")

	return exec.CommandContext(ctx, h.binary, h.args...), nil
}

func (h *Command) get(ctx context.Context, file string) error {

	staging := "./staging"
	if env := os.Getenv(config.EnvEncoderStage); len(env) != 0 {
		staging = env
	}

	info, err := os.Stat(staging)
	if err != nil {

		err := os.MkdirAll(staging, 0666)
		if err != nil {
			return ErrStagingDirectory
		}
	}

	if info != nil && !info.IsDir() {
		return ErrStagingDirectory
	}

	title, err := h.scan(ctx, file)
	if err != nil {
		return err
	}

	h.title = title

	h.parseAudioTracks()

	h.staged = filepath.Join(staging, filepath.Base(file))

	h.args = append(h.args, "-i", file, "-o", h.staged)

	return nil
}

func (h *Command) parse() {

	h.args = append(h.args,
		"--encoder-level", "4.1",
		"--format", h.profile.VideoContainer,
		"--encoder", h.profile.Codec,
		//"--encoder-preset", "medium",
		"--quality", fmt.Sprintf("%d", h.profile.QualityLevel),
		"-s", strings.Join(intToSlice(int(h.profile.SubtitleTracks)), ","),
	)

	if h.profile.AudioBitRate != 0 {
		h.args = append(h.args, "-B", fmt.Sprintf("%d", h.profile.AudioBitRate))
	}

	if len(h.profile.VideoTune) != 0 {
		h.args = append(h.args, "--encoder-tune", h.profile.VideoTune)
	}

	h.args = append(h.args, "-x", fmt.Sprintf("threads=%d", int(h.profile.Threads)))

}

// parseAudioTracks
// due to limitations in handbrake we must determine (based on the source) what kind of
// audio manipulations we're going to make
// This call sets our audio based handbrake arguments
func (h *Command) parseAudioTracks() {

	args := []string{
		"--audio-fallback", "aac",
		"--audio-copy-mask", "aac,ac3,dtshd,dts,mp3",
	}

	defer func() {
		h.args = append(h.args, args...)
	}()

	// Audio container format:
	mapping := h.profile.AudioCodecMap()
	tracks := h.title.AudioList

	// if there is only 1 codec
	if len(mapping) == 1 {

		// when that 1 codec is "copy"
		if _, ok := mapping[audio.CodecCopy]; ok {
			// copy all tracks
			args = append(args, "--aencoder", audio.CodecCopy.String())
			args = append(args, "--all-audio")
			return
		}

		// TODO remove duplicate tracks when there is just a single codec provided (?)

		var first audio.Codec
		for _, codec := range mapping {
			first = codec
			break
		}

		args = append(args, "--aencoder", first.String())
		args = append(args, "--all-audio")
		return
	}

	// handle foreign audio
	// apply simple converts only
	if h.title.MultiLanguageAudio() {

		var first audio.Codec
		for _, codec := range mapping {
			first = codec
			break
		}

		// convert based on the first mapping only
		args = append(args, "--aencoder", first.String())
		args = append(args, "--all-audio")
		return
	}

	conversion := make([]string, 0)
	highest := 0

	for idx, track := range tracks {
		if tracks[highest].AudioCodec().QualityRank() < track.AudioCodec().QualityRank() {
			highest = idx
			continue
		}

		if tracks[highest].Codec == track.Codec && tracks[highest].ChannelCount < track.ChannelCount {
			highest = idx
			continue
		}
	}

	best := tracks[highest] // the best available track

	// use the highest quality track for splitting
	for _, to := range mapping {

		// skip any codecs which would be of "higher" rank
		if best.AudioCodec().QualityRank() < to.QualityRank() {
			continue
		}

		// TODO is this needed?
		/*if from == audio.CodecCopy {
			conversion = append(conversion, to.CopyString())
			continue
		}*/

		conversion = append(conversion, to.String())
	}

	// if we couldnt successfully map to anything just convert to aac
	if len(conversion) == 0 {
		args = append(args, "--aencoder", audio.CodecAAC.String())
		args = append(args, "--all-audio")
		return
	}

	// below code used if we want to go a complete mapping route
	// we need to match the container with the track
	/*for _, track := range h.title.AudioList {
		if !track.IsKnown() {
			conversion = append(conversion, audio.CodecAAC.String())
			continue
		}

		// when we have a mapping for the audio
		if to, ok := mapping[track.AudioCodec()]; ok {
			conversion = append(conversion, to.String())
			continue
		}

		// when we do not have a mapping, we will just copy the track
		conversion = append(conversion, audio.CodecCopy.String())
	}*/

	args = append(
		args,
		"-a", repeatInt(highest+1, len(conversion)),
		"--aencoder", strings.Join(conversion, ","),
	)
}
