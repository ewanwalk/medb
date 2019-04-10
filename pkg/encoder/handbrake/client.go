package handbrake

import (
	"context"
	"encoder-backend/pkg/config"
	"encoder-backend/pkg/models"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	binary  string
	profile models.QualityProfile
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

	staging := "./staging"
	if env := os.Getenv(config.EnvEncoderStage); len(env) != 0 {
		staging = env
	}

	info, err := os.Stat(staging)
	if err != nil {

		err := os.MkdirAll(staging, 0666)
		if err != nil {
			return nil, ErrStagingDirectory
		}
	}

	if info != nil && !info.IsDir() {
		return nil, ErrStagingDirectory
	}

	h.staged = filepath.Join(staging, filepath.Base(file))

	h.args = append(h.args, "-i", file, "-o", h.staged)

	/*log.WithFields(log.Fields{
		"binary":  h.binary,
		"staging": h.staged,
		"args":    strings.Join(h.args, " "),
	}).Debug("handbrake.get: exec")*/

	return exec.CommandContext(ctx, h.binary, h.args...), nil
}

func (h *Command) parse() {

	h.args = append(h.args,
		"--encoder-level", "4.1",
		"--format", h.profile.VideoContainer,
		"--encoder", h.profile.Codec,
		"--quality", fmt.Sprintf("%d", h.profile.QualityLevel),
		"--aencoder", h.profile.AudioContainer,
		"-B", fmt.Sprintf("%d", h.profile.AudioBitRate),
		"-a", strings.Join(intToSlice(int(h.profile.AudioTracks)), ","),
		"-s", strings.Join(intToSlice(int(h.profile.SubtitleTracks)), ","),
	)

	if len(h.profile.VideoTune) != 0 {
		h.args = append(h.args, "--encoder-tune", h.profile.Codec)
	}

	h.args = append(h.args, "-x", fmt.Sprintf("threads=%d", int(h.profile.Threads)))

}
