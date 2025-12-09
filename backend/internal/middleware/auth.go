package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/auth"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
)

const accessKeyCtx ContextKey = "accessKey"

func ErrSignatureDoesNotMatch() *core.S3Error {
	return &core.S3Error{
		Code:       "SignatureDoesNotMatch",
		Message:    "The request signature we calculate does not match the signature you provided.",
		StatusCode: http.StatusForbidden,
	}
}

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
			core.HandleError(w, ErrSignatureDoesNotMatch())
			return
		}

		if res.IsStreaming {
			r.Body = auth.NewSigV4ChunkedReader(r, res)
		}

		ctx := context.WithValue(r.Context(), accessKeyCtx, res.Credential.AccessKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
