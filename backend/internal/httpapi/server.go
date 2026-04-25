package httpapi

import (
	"net/http"
	"os"
	"time"

	"github.com/ghassenelkamel/portfolio-v2/internal/config"
	"github.com/ghassenelkamel/portfolio-v2/internal/middleware"
)

func NewServer(cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/api/", NewAPI(cfg))

	// In prod, Go will serve static files from WEB_DIR (copied into ./web in Dockerfile).
	// In dev (no ./web), this provides a minimal response.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(cfg.WebDir); err == nil {
			http.FileServer(http.Dir(cfg.WebDir)).ServeHTTP(w, r)
			return
		}
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("dev backend running. try /api/healthz\n"))
	})

	h := middleware.Gzip(middleware.Timeout(mux, 30*time.Second))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}
}
