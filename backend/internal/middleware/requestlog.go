package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type responseLogger struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseLogger) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseLogger) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func RequestLog(next http.Handler) http.Handler {
	logger := slog.Default()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseLogger{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		status := rw.status
		if status == 0 {
			status = http.StatusOK
		}
		durationMs := time.Since(start).Milliseconds()

		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"bytes", rw.bytes,
			"duration_ms", durationMs,
			"remote_ip", clientIP(r),
		}

		if ua := strings.TrimSpace(r.UserAgent()); ua != "" {
			attrs = append(attrs, "ua", ua)
		}
		if reqID := strings.TrimSpace(r.Header.Get("X-Request-ID")); reqID != "" {
			attrs = append(attrs, "request_id", reqID)
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("http_request", attrs...)
		case status >= http.StatusBadRequest:
			logger.Warn("http_request", attrs...)
		case durationMs >= 2000:
			logger.Warn("http_request_slow", attrs...)
		case r.URL.Path == "/api/healthz":
			logger.Debug("http_request", attrs...)
		default:
			logger.Info("http_request", attrs...)
		}
	})
}

func clientIP(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); v != "" {
		if first, _, ok := strings.Cut(v, ","); ok {
			return strings.TrimSpace(first)
		}
		return v
	}

	if v := strings.TrimSpace(r.Header.Get("X-Real-IP")); v != "" {
		return v
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
