package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// loggingResponseWriter wraps http.ResponseWriter to capture the status code for logging.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader captures the status code and calls the underlying WriteHeader method.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware is an HTTP middleware that logs each request with structured fields.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)

		fields := logrus.Fields{
			"component":   "http",
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     lrw.statusCode,
			"duration_ms": duration.Milliseconds(),
			"client_ip":  clientIP,
		}
		if r.URL.RawQuery != "" {
			fields["query"] = r.URL.RawQuery
		}

		logrus.WithFields(fields).Info("request completed")
	})
}
