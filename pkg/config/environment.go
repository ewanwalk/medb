package config

import (
	"github.com/joho/godotenv"
)

const (
	EnvLogLevel = "APP_LOG_LEVEL"

	EnvDBUsername = "DB_USERNAME"
	EnvDBPassword = "DB_PASSWORD"
	EnvDBHostname = "DB_HOSTNAME"
	EnvDBPort     = "DB_PORT"
	EnvDBName     = "DB_NAME"
)

func init() {
	_ = godotenv.Load()
}
