package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr string

	GitHubUser  string
	GitHubToken string

	ContentRepoOwner string
	ContentRepoName  string
	ContentRepoRef   string
	ContentRepoDir   string
	ContentCacheTTL  int

	ContactForwardURL string

	MaxBodyBytes int64
	WebDir       string
}

func Load() Config {
	maxBody := int64(256 * 1024)
	if v := os.Getenv("MAX_BODY_BYTES"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			maxBody = n
		}
	}

	return Config{
		Addr:              getenv("ADDR", ":8080"),
		GitHubUser:        getenv("GITHUB_USER", "ghassenelkamel"),
		GitHubToken:       os.Getenv("GITHUB_TOKEN"),
		ContentRepoOwner:  getenv("CONTENT_REPO_OWNER", "ghassenelkamel"),
		ContentRepoName:   getenv("CONTENT_REPO_NAME", "portfolio-content"),
		ContentRepoRef:    getenv("CONTENT_REPO_REF", "main"),
		ContentRepoDir:    strings.Trim(getenv("CONTENT_REPO_DIR", "content"), "/"),
		ContentCacheTTL:   getenvInt("CONTENT_CACHE_TTL", 300),
		ContactForwardURL: os.Getenv("CONTACT_FORWARD_URL"),
		MaxBodyBytes:      maxBody,
		WebDir:            getenv("WEB_DIR", "./web"),
	}
}

func getenvInt(k string, def int) int {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func getenv(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}
