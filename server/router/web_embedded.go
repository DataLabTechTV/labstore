//go:build embedassets

package router

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/helper"
)

//go:embed assets/**
var assetsFS embed.FS

func NewWebServerDescriptor(host string, port uint16) (*ServerDescriptor, error) {
	slog.Info("web ui server", "host", host, "port", port)
	slog.Debug("using embedded web ui assets")

	addr := fmt.Sprintf("%s:%d", host, port)

	contentFS := helper.Must(fs.Sub(assetsFS, "assets"))
	httpFS := http.FS(contentFS)
	handler := http.FileServer(httpFS)

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	serverDescriptor := &ServerDescriptor{
		Name:   "Web UI",
		Server: &server,
	}

	return serverDescriptor, nil
}
