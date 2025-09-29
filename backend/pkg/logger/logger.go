package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type AppLogger struct {
	Logger *zerolog.Logger
}

func New() *AppLogger {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &AppLogger{
		Logger: &logger,
	}
}

func (l *AppLogger) SetLogLevel(level int) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}
