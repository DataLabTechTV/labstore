package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
)

func NewIAMServerDescriptor(host string, port uint16) *ServerDescriptor {
	slog.Info(
		"iam server",
		"host", host,
		"port", port,
	)

	router := http.NewServeMux()
	loadIAMRoutes(router)

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
		Name:   "IAM",
		Server: &server,
	}

	return serverDescriptor
}

func loadIAMRoutes(router *http.ServeMux) {
	slog.Debug("loading iam routes")

	router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		action := iam.IAMOp(r.URL.Query().Get("Action"))

		switch action {
		case iam.OpCreateUser:
			iam.CreateUserHandler(w, r)
		// case iam.OpCreateAccessKey:
		// 	CreateAccessKeyHandler(w, r)
		// case iam.OpCreateGroup:
		// 	CreateGroupHandler(w, r)
		// case iam.OpAttachUserPolicy:
		// 	AttachUserPolicy(w, r)
		// case iam.OpAttachGroupPolicy:
		// 	AttachGroupPolicy(w, r)
		default:
			core.HandleError(w, iam.ErrNotImplemented(action))
		}
	})
}
