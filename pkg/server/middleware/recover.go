package middleware

import (
	"net/http"
	"runtime/debug"

	"Cloud/pkg/logger"
)

// Recover catch panic, log and return HTTP 500
func Recover(logger logger.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.WithError(err.(error)).
						With("stack", string(stack)).
						Error("panic recover")

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("internal server error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
