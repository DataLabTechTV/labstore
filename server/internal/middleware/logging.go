package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Tracks request time and logs request after all handlers
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		slog.Debug("logging middleware")

		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		}

		if r.URL.RawQuery != "" {
			attrs = append(attrs, slog.String("query", r.URL.RawQuery))
		}

		if len(r.Form) > 0 {
			attrs = append(attrs, slog.String("form", r.Form.Encode()))
		}

		attrs = append(attrs, slog.Int("status", rw.status))

		if rw.size > 0 {
			attrs = append(attrs, slog.Int("size", rw.size))
		}

		attrs = append(attrs, slog.Float64("duration_ms", float64(time.Since(start).Nanoseconds())/1e6))

		slog.Info("request", attrs...)
	})
}
