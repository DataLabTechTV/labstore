package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/lmittmann/tint"
)

const DefaultLogLevel = slog.LevelInfo

var (
	DefaultLogWriter  = os.Stderr
	DefaultTimeFormat = time.StampMilli
)

var (
	Level     slog.LevelVar
	AppLogger *slog.Logger
)

type Option func(*slog.Logger)

func Init(opts ...Option) {
	InitWithWriter(DefaultLogWriter, opts...)
}

func InitWithWriter(w io.Writer, opts ...Option) {
	Level.Set(DefaultLogLevel)

	noColor := true
	if file, ok := w.(*os.File); ok {
		noColor = !isatty.IsTerminal(file.Fd())
	}

	AppLogger = slog.New(
		tint.NewHandler(
			w,
			&tint.Options{
				Level:      &Level,
				TimeFormat: DefaultTimeFormat,
				NoColor:    noColor,
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

func Swap(w io.Writer, opts ...Option) func() {
	previous := slog.Default()

	revert := func() {
		AppLogger = previous
		slog.SetDefault(AppLogger)
		slog.Debug("reverted to default logger")
	}

	InitWithWriter(w, opts...)

	return revert
}

func NewDailyWriter(dir, prefix string) (io.Writer, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s-%s.log", prefix, time.Now().Format("2006-01-02"))
	path := filepath.Join(dir, filename)

	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
}
