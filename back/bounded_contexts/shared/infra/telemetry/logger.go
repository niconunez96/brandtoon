package telemetry

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	LogError(msg string, err error, attrs any)
	LogInfo(msg string, attrs any)
}

type ZeroLogLogger struct {
	zerolog.Logger
}

func newZeroLogLogger() ZeroLogLogger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return ZeroLogLogger{
		zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

func (l ZeroLogLogger) LogError(msg string, err error, attrs any) {
	l.Error().Fields(attrs).Err(err).Msg(msg)
}

func (l ZeroLogLogger) LogInfo(msg string, attrs any) {
	l.Info().Fields(attrs).Msg(msg)
}

func newLogger() Logger {
	return newZeroLogLogger()
}
