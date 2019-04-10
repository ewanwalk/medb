package job

import (
	"encoder-backend/pkg/encoder/handbrake"
	"encoder-backend/pkg/models"
)

type option func(*options) error

type options struct {
	handbrake *handbrake.Command
}

// withProfile
// utilizing a quality profile to produce a command
func withProfile(profile models.QualityProfile) option {
	return func(o *options) error {

		command, err := handbrake.New(profile)
		if err != nil {
			return err
		}

		o.handbrake = command

		return nil

	}
}
