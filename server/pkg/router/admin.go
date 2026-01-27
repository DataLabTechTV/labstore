package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/IllumiKnowLabs/labstore/server/internal/middleware"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
)

var healthStatus []*atomic.Bool

func NewAdminServerDescriptor(host string, port uint16, monitor []*ServerDescriptor) *ServerDescriptor {
	slog.Info("admin api server", "host", host, "port", port, "monitor", len(monitor))

	healthStatus = make([]*atomic.Bool, len(monitor))
	for i := range monitor {
		healthStatus[i] = &monitor[i].Healthy
	}

	router := http.NewServeMux()
	loadAdminRoutes(router)

	addr := fmt.Sprintf("%s:%d", host, port)

	mw := middleware.Stack(
		middleware.LoggingMiddleware,
		middleware.CompressionMiddleware,
		middleware.NormalizationMiddleware,
	)

	server := http.Server{
		Addr:    addr,
		Handler: mw(router),
	}

	serverDescriptor := &ServerDescriptor{
		Name:   "Admin API",
		Server: &server,
	}

	return serverDescriptor
}

func loadAdminRoutes(router *http.ServeMux) {
	slog.Debug("loading admin routes")
	router.HandleFunc("GET /health", http.HandlerFunc(healthCheckHandler))
}

// healthCheckHandler: GET /health
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	for _, hs := range healthStatus {
		if healthy := hs.Load(); !healthy {
			helper.Must(w.Write([]byte("LabStore is not running")))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	helper.Must(w.Write([]byte("LabStore is running")))
	w.WriteHeader(http.StatusOK)
}
