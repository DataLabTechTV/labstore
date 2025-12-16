package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
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
		action := r.URL.Query().Get("Action")

		// TODO: support for missing operations
		switch iam.IAMOp(action) {
		// --- User ---
		case iam.OpCreateUser:
			iam.CreateUserHandler(w, r)
		case iam.OpCreateAccessKey:
			iam.CreateAccessKeyHandler(w, r)
		case iam.OpAttachUserPolicy:
			iam.AttachUserPolicyHandler(w, r)

		// --- Groups ---
		case iam.OpCreateGroup:
			iam.CreateGroupHandler(w, r)
		case iam.OpAddUserToGroup:
			iam.AddUserToGroupHandler(w, r)
		case iam.OpAttachGroupPolicy:
			iam.AttachGroupPolicyHandler(w, r)

		// --- Policies ---
		case iam.OpCreatePolicy:
			iam.CreatePolicyHandler(w, r)

		default:
			errs.Handle(w, errs.IAMNotImplemented(action))
		}
	})
}
