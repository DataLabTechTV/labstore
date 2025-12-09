package middleware

import (
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
)

func WithIAM(action iam.Action, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("with iam", "action", action)
		if action == "" {
			return
		}

		bucket := r.PathValue("bucket")
		key := r.PathValue("key")
		accessKey := GetRequestAccessKey(r)

		if !iam.CheckPolicy(accessKey, bucket, key, iam.Action(action)) {
			// !FIXME: AWS compliant error handling?
			core.HandleError(w, ErrAccessDenied())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ErrAccessDenied() *core.S3Error {
	return &core.S3Error{
		Code:       "AccessDenied",
		Message:    "AccessDenied",
		StatusCode: 403,
	}
}
