package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ghassenelkamel/portfolio-v2/internal/config"
	"github.com/ghassenelkamel/portfolio-v2/internal/httpapi"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: parseLogLevel(os.Getenv("LOG_LEVEL"))}))
	slog.SetDefault(logger)

	logger.Info(
		"starting portfolio api",
		"addr", cfg.Addr,
		"web_dir", cfg.WebDir,
		"content_repo", cfg.ContentRepoOwner+"/"+cfg.ContentRepoName,
		"content_ref", cfg.ContentRepoRef,
		"smtp_host", cfg.SMTPHost,
		"smtp_to", cfg.SMTPTo,
		"smtp_tls", cfg.SMTPRequireTLS,
		"forward_set", cfg.ContactForwardURL != "",
	)

	srv := httpapi.NewServer(cfg)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func parseLogLevel(v string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
