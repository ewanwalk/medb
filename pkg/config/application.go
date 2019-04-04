package config

import "runtime"

var (
	Separator = "/"
)

func init() {
	if runtime.GOOS == "windows" {
		Separator = "\\"
	}
}
