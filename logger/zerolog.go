package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func GetLogger() (zerolog.Logger, error) {
	level := zerolog.DebugLevel

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(level)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger(), nil
}
