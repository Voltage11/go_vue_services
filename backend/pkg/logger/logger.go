package logger

import (
	"os"

	"github.com/rs/zerolog"
)

//type LoggerType string
//
//const (
//	LoggerDebug LoggerType = "logger-debug"
//	LoggerInfo LoggerType = "logger-info"
//	LoggerError LoggerType = "logger-error"
//	LoggerFatal LoggerType = "logger-fatal"
//	LoggerPanic LoggerType = "logger-panic"
//)

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
