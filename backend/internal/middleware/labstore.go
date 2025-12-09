package middleware

import (
	"net/http"
	"strings"
)

func LabStoreMiddleware(next http.Handler) http.Handler {
	labstoreRouter := http.NewServeMux()
	mw := Stack(NormalizationMiddleware)

	labstoreRouter.HandleFunc("GET /labstore/health", http.HandlerFunc(HealthCheckHandler))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/labstore/") {
			mw(labstoreRouter).ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HealthCheckHandler: GET /labstore/health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
