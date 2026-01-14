package router

import (
	"fmt"
	"log/slog"
	"net/http"

	iamimpl "github.com/IllumiKnowLabs/labstore/backend/internal/iam"
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
			iamimpl.CreateUserHandler(w, r)
		case iam.OpCreateAccessKey:
			iamimpl.CreateAccessKeyHandler(w, r)

		// --- User: Read ---
		case iam.OpGetUser:
			iamimpl.GetUserHandler(w, r)
		case iam.OpListAccessKeys:
			iamimpl.ListAccessKeysHandler(w, r)
		case iam.OpListAttachedUserPolicies:
			iamimpl.ListAttachedUserPoliciesHandler(w, r)

		// --- User: Delete ---
		case iam.OpDeleteUser:
			iamimpl.DeleteUserHandler(w, r)
		case iam.OpDeleteAccessKey:
			iamimpl.DeleteAccessKeyHandler(w, r)

		// --- Groups: Create ---
		case iam.OpCreateGroup:
			iamimpl.CreateGroupHandler(w, r)
		case iam.OpAddUserToGroup:
			iamimpl.AddUserToGroupHandler(w, r)

		// --- Groups: Read ---
		case iam.OpGetGroup:
			iamimpl.GetGroupHandler(w, r)
		case iam.OpListAttachedGroupPolicies:
			iamimpl.ListAttachedGroupPoliciesHandler(w, r)

		// --- Groups: Delete ---
		case iam.OpDeleteGroup:
			iamimpl.DeleteGroupHandler(w, r)
		case iam.OpRemoveUserFromGroup:
			iamimpl.RemoveUserFromGroupHandler(w, r)

		// --- Policies: Create ---
		case iam.OpCreatePolicy:
			iamimpl.CreatePolicyHandler(w, r)
		case iam.OpAttachUserPolicy:
			iamimpl.AttachUserPolicyHandler(w, r)
		case iam.OpAttachGroupPolicy:
			iamimpl.AttachGroupPolicyHandler(w, r)

		// --- Policies: Read ---
		case iam.OpGetPolicy:
			iamimpl.GetPolicyHandler(w, r)

		// --- Policies: Delete ---
		case iam.OpDeletePolicy:
			iamimpl.DeletePolicyHandler(w, r)
		case iam.OpDetachUserPolicy:
			iamimpl.DetachUserPolicyHandler(w, r)
		case iam.OpDetachGroupPolicy:
			iamimpl.DetachGroupPolicyHandler(w, r)

		default:
			errs.Handle(w, errs.IAMNotImplemented(action))
		}
	})
}
