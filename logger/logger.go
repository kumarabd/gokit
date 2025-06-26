package logger

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

// Handler wraps zerolog.Logger to allow method definitions
type Handler struct {
	zerolog.Logger
}

// New instantiates bucky logger instance
func New(appname string, opts Options) (*Handler, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.With().Str("app", appname).Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return &Handler{logger}, nil
}

func (l *Handler) AsLogrLogger() logr.Logger {
	return zerologr.New(&l.Logger)
}
