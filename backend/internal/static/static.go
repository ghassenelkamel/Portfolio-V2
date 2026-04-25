package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"log/slog"
)

func Handler(webDir string, logger *slog.Logger) http.Handler {
	// If no web dir exists (local backend-only), serve a minimal landing.
	if st, err := os.Stat(webDir); err != nil || !st.IsDir() {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("portfolio backend running (no static web dir). Try /api/healthz\n"))
		})
	}

	fs := http.FileServer(http.Dir(webDir))
	fallback := filepath.Join(webDir, "200.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve file
		rr := &responseRecorder{ResponseWriter: w, status: 200}
		fs.ServeHTTP(rr, r)

		// SPA fallback for routes (SvelteKit adapter-static with fallback 200.html)
		if rr.status == 404 {
			http.ServeFile(w, r, fallback)
			return
		}
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
