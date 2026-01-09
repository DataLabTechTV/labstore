package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/lmittmann/tint"
)

const DefaultLogLevel = slog.LevelInfo

var (
	DefaultLogOutput  = os.Stderr
	DefaultTimeFormat = time.StampMilli
)

var (
	Level     slog.LevelVar
	AppLogger *slog.Logger
)

type Option func(*slog.Logger)

func Temporary(output io.Writer, opts ...Option) func() {
	previous := slog.Default()

	revert := func() {
		AppLogger = previous
		slog.SetDefault(AppLogger)
	}

	InitWithOutput(output, opts...)

	return revert
}

func Init(opts ...Option) {
	InitWithOutput(DefaultLogOutput, opts...)
}

func InitWithOutput(output io.Writer, opts ...Option) {
	Level.Set(DefaultLogLevel)

	AppLogger = slog.New(
		tint.NewHandler(
			output,
			&tint.Options{
				Level:      &Level,
				TimeFormat: DefaultTimeFormat,
				NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
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

	return WithLevel(DefaultLogLevel)
}

func WithLevel(level slog.Level) Option {
	return func(logger *slog.Logger) {
		Level.Set(level)

		if level == slog.LevelDebug {
			logger.Debug("Debug mode: on")
		}
	}
}
