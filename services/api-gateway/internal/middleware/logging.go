package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			if log != nil {
				log.Info("request", zap.String("path", r.URL.Path), zap.Duration("duration", time.Since(start)))
			}
		})
	}
}
