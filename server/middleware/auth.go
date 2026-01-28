package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/auth"
	"github.com/IllumiKnowLabs/labstore/server/errs"
)

type contextKey string

const accessKeyCtx contextKey = "accessKey"

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

		sigV4Ctx, err := auth.VerifySigV4(r)
		if err != nil {
			slog.Error("sigv4", "err", err)
			errs.Handle(w, errs.S3SignatureDoesNotMatch())
			return
		}

		if sigV4Ctx.IsStreaming {
			r.Body = auth.NewSigV4ChunkedDecoder(sigV4Ctx, r.Body)
		}

		ctx := context.WithValue(r.Context(), accessKeyCtx, sigV4Ctx.Credential.AccessKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
