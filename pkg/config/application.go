package config

import "runtime"

const (
	AppBindTo = ":80"
)

var (
	Separator = "/"
)

func init() {
	if runtime.GOOS == "windows" {
		Separator = "\\"
	}
}
