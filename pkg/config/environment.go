package config

import (
	_ "github.com/joho/godotenv/autoload"
)

const (
	EnvLogLevel = "APP_LOG_LEVEL"

	EnvHandbrake = "APP_HANDBRAKE"

	EnvDBUsername = "DB_USERNAME"
	EnvDBPassword = "DB_PASSWORD"
	EnvDBHostname = "DB_HOSTNAME"
	EnvDBPort     = "DB_PORT"
	EnvDBName     = "DB_NAME"
	EnvDBDebug    = "DB_DEBUG" // enables debug query logging

	EnvEncoderConcurrency = "ENCODER_CONCURRENCY"
	EnvEncoderStage       = "ENCODER_STAGE"
	EnvEncoderReportLogs  = "ENCODER_REPORT_LOGGING"
)
