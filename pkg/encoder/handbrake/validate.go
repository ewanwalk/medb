package handbrake

import (
	"encoder-backend/pkg/config"
	"errors"
	"os"
	"os/exec"
)

var (
	ErrBinaryNotFound = errors.New("handbrake.validate: binary not found")
)

func (h *Command) validate() error {

	binary := "HandBrakeCLI"

	if env := os.Getenv(config.EnvHandbrake); len(env) != 0 {
		binary = env
	}

	cmd := exec.Command("/bin/sh", "-c", "command -v", binary)
	if err := cmd.Run(); err != nil {
		return ErrBinaryNotFound
	}

	h.binary = binary

	return nil
}
