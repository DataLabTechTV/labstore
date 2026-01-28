package profiler

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
)

const DefaultHost = "localhost"
const DefaultPort = 6060

type Profiler struct {
	Host string
	Port int
}

type Option func(*Profiler)

func NewProfiler(opts ...Option) *Profiler {
	profiler := &Profiler{
		Host: DefaultHost,
		Port: DefaultPort,
	}

	for _, opt := range opts {
		opt(profiler)
	}

	return profiler
}

func WithHost(host string) Option {
	return func(profiler *Profiler) {
		profiler.Host = host
	}
}

func WithPort(port int) Option {
	return func(profiler *Profiler) {
		profiler.Port = port
	}
}

func (profiler *Profiler) Start() {
	slog.Info("starting pprof profiler", "port", profiler.Port)

	go func() {
		addr := fmt.Sprintf("%s:%d", profiler.Host, profiler.Port)
		fmt.Printf("🕵️  Profiler listening on http://%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}
