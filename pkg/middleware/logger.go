package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingResponseWriter wraps http.ResponseWriter to capture status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := &LoggingResponseWriter{w, http.StatusOK}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			logger.WithFields(map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     lrw.statusCode,
				"duration":   duration,
				"remote_ip":  r.RemoteAddr,
				"user_agent": r.UserAgent(),
				"referrer":   r.Referer(),
			}).Info("request completed")
		})
	}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
