package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/auth"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
)

type ContextKey string

const accessKeyCtx ContextKey = "accessKey"

func GetRequestAccessKey(r *http.Request) string {
	if accessKey := r.Context().Value(accessKeyCtx); accessKey != nil {
		return accessKey.(string)
	}

	return ""
}

// Must come before middleware that changes the request (e.g., NormalizeMiddleware)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("auth middleware")

		res, err := auth.VerifySigV4(r)
		if err != nil {
			slog.Error("sigv4", "err", err)
			errs.Handle(w, errs.S3SignatureDoesNotMatch())
			return
		}

		if res.IsStreaming {
			r.Body = auth.NewSigV4ChunkedReader(r, res)
		}

		ctx := context.WithValue(r.Context(), accessKeyCtx, res.Credential.AccessKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
