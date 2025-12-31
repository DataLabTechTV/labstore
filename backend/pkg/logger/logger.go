package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/lmittmann/tint"
)

const defaultLogLevel = slog.LevelInfo

var defaultLogOutput = os.Stderr
var defaultTimeFormat = time.StampMilli

var Level slog.LevelVar
var AppLogger *slog.Logger

type Option func(*slog.Logger)

func Init(opts ...Option) {
	Level.Set(defaultLogLevel)

	AppLogger = slog.New(
		tint.NewHandler(
			defaultLogOutput,
			&tint.Options{
				Level:      &Level,
				TimeFormat: defaultTimeFormat,
				NoColor:    !isatty.IsTerminal(defaultLogOutput.Fd()),
			},
		),
	)

	for _, opt := range opts {
		opt(AppLogger)
	}

	slog.SetDefault(AppLogger)
}

func WithDebugFlag(debug bool) Option {
	if debug {
		return WithLevel(slog.LevelDebug)
	}

	return WithLevel(defaultLogLevel)
}

func WithLevel(level slog.Level) Option {
	return func(logger *slog.Logger) {
		Level.Set(level)

		if level == slog.LevelDebug {
			logger.Debug("Debug mode: on")
		}
	}
}
