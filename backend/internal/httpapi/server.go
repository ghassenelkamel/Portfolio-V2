package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ghassenelkamel/portfolio-v2/internal/config"
	"github.com/ghassenelkamel/portfolio-v2/internal/middleware"
	"github.com/ghassenelkamel/portfolio-v2/internal/static"
)

func NewServer(cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/api/", NewAPI(cfg))
	mux.Handle("/", static.Handler(cfg.WebDir, slog.Default()))

	h := middleware.RequestLog(middleware.Gzip(middleware.Timeout(mux, 30*time.Second)))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}
}
