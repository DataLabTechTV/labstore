package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/iam"
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
		if err := r.ParseForm(); err != nil {
			slog.Error("load iam router: could not parse form")
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		action := r.Form.Get("Action")

		switch iam.IAMOp(action) {
		// --- User: Create ---
		case iam.OpCreateUser:
			iam.CreateUserHandler(w, r)
		case iam.OpCreateAccessKey:
			iam.CreateAccessKeyHandler(w, r)

		// --- User: Read ---
		case iam.OpGetUser:
			iam.GetUserHandler(w, r)
		case iam.OpListAccessKeys:
			iam.ListAccessKeysHandler(w, r)
		case iam.OpListAttachedUserPolicies:
			iam.ListAttachedUserPoliciesHandler(w, r)

		// --- User: Delete ---
		case iam.OpDeleteUser:
			iam.DeleteUserHandler(w, r)
		case iam.OpDeleteAccessKey:
			iam.DeleteAccessKeyHandler(w, r)

		// --- Groups: Create ---
		case iam.OpCreateGroup:
			iam.CreateGroupHandler(w, r)
		case iam.OpAddUserToGroup:
			iam.AddUserToGroupHandler(w, r)

		// --- Groups: Read ---
		case iam.OpGetGroup:
			iam.GetGroupHandler(w, r)
		case iam.OpListAttachedGroupPolicies:
			iam.ListAttachedGroupPoliciesHandler(w, r)

		// --- Groups: Delete ---
		case iam.OpDeleteGroup:
			iam.DeleteGroupHandler(w, r)
		case iam.OpRemoveUserFromGroup:
			iam.RemoveUserFromGroupHandler(w, r)

		// --- Policies: Create ---
		case iam.OpCreatePolicy:
			iam.CreatePolicyHandler(w, r)
		case iam.OpAttachUserPolicy:
			iam.AttachUserPolicyHandler(w, r)
		case iam.OpAttachGroupPolicy:
			iam.AttachGroupPolicyHandler(w, r)

		// --- Policies: Read ---
		case iam.OpGetPolicy:
			iam.GetPolicyHandler(w, r)

		// --- Policies: Delete ---
		case iam.OpDeletePolicy:
			iam.DeletePolicyHandler(w, r)
		case iam.OpDetachUserPolicy:
			iam.DetachUserPolicyHandler(w, r)
		case iam.OpDetachGroupPolicy:
			iam.DetachGroupPolicyHandler(w, r)

		default:
			errs.Handle(w, errs.IAMNotImplemented(action))
		}
	})
}
