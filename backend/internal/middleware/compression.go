package middleware

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/klauspost/compress/zstd"
)

func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("compression middleware")

		var reader io.Reader = r.Body

		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			slog.Debug("decompressing gzip request")

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "invalid gzip", http.StatusBadRequest)
				return
			}
			defer helper.CloseWithErr(gz, &err)

			reader = gz

		case "zstd":
			slog.Debug("decompressing zstd request")

			zr, err := zstd.NewReader(r.Body)
			if err != nil {
				http.Error(w, "invalid zstd", http.StatusBadRequest)
				return
			}
			defer zr.Close()

			reader = zr
		}

		r.Body = io.NopCloser(reader)
		next.ServeHTTP(w, r)
	})
}
