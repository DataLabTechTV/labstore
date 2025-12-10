package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
)

func NewAdminServerDescriptor(host string, port uint16) *ServerDescriptor {
	slog.Info(
		"admin api server",
		"host", host,
		"port", port,
	)

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

// healthCheckHandler: GET /labstore/health
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
