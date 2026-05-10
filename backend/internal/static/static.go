package static

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"log/slog"
)

func Handler(webDir string, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	// If no web dir exists (local backend-only), serve a minimal landing.
	if st, err := os.Stat(webDir); err != nil || !st.IsDir() {
		logger.Info("static web dir not found; serving backend-only message", "web_dir", webDir)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("cache-control", "no-store")
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("portfolio backend running (no static web dir). Try /api/healthz\n"))
		})
	}

	fs := http.FileServer(http.Dir(webDir))
	fallback := pickFallbackFile(webDir)
	if fallback != "" {
		logger.Info("static SPA fallback enabled", "fallback", fallback)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			setStaticCacheHeaders(w, r.URL.Path, false)
			fs.ServeHTTP(w, r)
			return
		}

		if !isSPAPath(r.URL.Path) {
			setStaticCacheHeaders(w, r.URL.Path, false)
			fs.ServeHTTP(w, r)
			return
		}

		if pathLooksLikeAsset(r.URL.Path) {
			setStaticCacheHeaders(w, r.URL.Path, false)
			fs.ServeHTTP(w, r)
			return
		}

		if fallback != "" {
			setStaticCacheHeaders(w, r.URL.Path, true)
			http.ServeFile(w, r, fallback)
			return
		}

		setStaticCacheHeaders(w, r.URL.Path, false)
		fs.ServeHTTP(w, r)
	})
}

func pickFallbackFile(webDir string) string {
	for _, name := range []string{"index.html", "200.html"} {
		full := filepath.Join(webDir, name)
		if st, err := os.Stat(full); err == nil && !st.IsDir() {
			return full
		}
	}
	return ""
}

func isSPAPath(reqPath string) bool {
	return reqPath != "" && reqPath != "/"
}

func pathLooksLikeAsset(reqPath string) bool {
	clean := path.Clean("/" + reqPath)
	if strings.HasPrefix(clean, "/assets/") {
		return true
	}
	if ext := path.Ext(clean); ext != "" {
		return true
	}
	if strings.Contains(clean, ".well-known") {
		return true
	}
	return false
}

func setStaticCacheHeaders(w http.ResponseWriter, reqPath string, fallback bool) {
	if fallback {
		w.Header().Set("cache-control", "no-cache")
		return
	}

	clean := path.Clean("/" + reqPath)
	ext := strings.ToLower(path.Ext(clean))

	if clean == "/" || ext == ".html" {
		w.Header().Set("cache-control", "no-cache")
		return
	}

	if strings.HasPrefix(clean, "/assets/") || strings.Contains(clean, "/immutable/") {
		w.Header().Set("cache-control", "public, max-age=31536000, immutable")
		return
	}

	switch ext {
	case ".js", ".mjs", ".css", ".svg", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".woff", ".woff2":
		w.Header().Set("cache-control", "public, max-age=604800, stale-while-revalidate=86400")
	default:
		w.Header().Set("cache-control", "public, max-age=300, stale-while-revalidate=300")
	}
}
