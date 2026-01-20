package middleware

import (
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/iam"
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
			errs.Handle(w, errs.S3AccessDenied())
			return
		}

		next.ServeHTTP(w, r)
	})
}
