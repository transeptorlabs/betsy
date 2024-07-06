package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetLogger(logLevel string) (zerolog.Logger, error) {
	level, err := getLevel(logLevel)
	if err != nil {
		return log.Logger, err
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(level)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger(), nil
}

func getLevel(logLevel string) (zerolog.Level, error) {
	// Defaults to info loevel
	level := zerolog.InfoLevel
	switch logLevel {
	case "TRACE":
		level = zerolog.TraceLevel
	case "DEBUG":
		level = zerolog.DebugLevel
	case "INFO":
		level = zerolog.InfoLevel
	case "WARN":
		level = zerolog.WarnLevel
	case "ERROR":
		level = zerolog.ErrorLevel
	case "OFF":
		level = zerolog.Disabled
	default:
		return level, fmt.Errorf("Invalid log level: %s", logLevel)
	}
	return level, nil
}
