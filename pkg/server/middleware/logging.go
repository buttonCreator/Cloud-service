package middleware

import (
	"net/http"
	"time"

	"Cloud/pkg/logger"
)

// AccessLog middleware logs HTTP requests
func AccessLog(logger logger.ILogger, checkPathIsProbe func(path string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "/status" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r)

			statusCode := rw.status
			if statusCode >= 200 && statusCode < 300 && checkPathIsProbe(path) {
				return
			}

			logger.With(
				"duration", time.Since(start).Seconds(),
				"code", statusCode,
			).Info("HTTP request")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
