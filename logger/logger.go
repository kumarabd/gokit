package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Handler is a type alias for zerolog.Logger
type Handler = zerolog.Logger

// New instantiates bucky logger instance
func New(appname string, opts Options) (*Handler, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.With().Str("app", appname).Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return &logger, nil
}
