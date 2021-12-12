package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string) *zerolog.Logger {
	var lev zerolog.Level

	switch level {
	case "info":
		lev = zerolog.InfoLevel
	case "debug":
		lev = zerolog.DebugLevel
	default:
		lev = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stderr).With().Logger().Level(lev)

	return &logger
}
