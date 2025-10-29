package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const defaultLogLevel = logrus.InfoLevel

var (
	defaultLogOutput    = os.Stdout
	defaultLogFormatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
)

var Log *logrus.Logger

type Option func(*logrus.Logger)

func Init(opts ...Option) {
	Log = logrus.New()

	Log.SetOutput(defaultLogOutput)
	Log.SetLevel(defaultLogLevel)
	Log.SetFormatter(defaultLogFormatter)

	for _, opt := range opts {
		opt(Log)
	}
}

func WithDebugFlag(debug bool) Option {
	if debug {
		return WithLevel(logrus.DebugLevel)
	}

	return WithLevel(defaultLogLevel)
}

func WithLevel(level logrus.Level) Option {
	return func(logger *logrus.Logger) {
		if level == logrus.DebugLevel {
			Log.Debug("Debug mode: on")
		}

		logger.SetLevel(level)
	}
}

func WithFormatter(formatter logrus.Formatter) Option {
	return func(logger *logrus.Logger) {
		logger.SetFormatter(formatter)
	}
}

func WithOutput(output *os.File) Option {
	return func(logger *logrus.Logger) {
		logger.SetOutput(output)
	}
}
