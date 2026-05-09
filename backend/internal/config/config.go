package config

import (
	"bufio"
	"os"
	"path/filepath"
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
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPass          string
	SMTPFrom          string
	SMTPTo            string
	SMTPRequireTLS    bool

	MaxBodyBytes int64
	WebDir       string
}

func Load() Config {
	loadDotEnvIfPresent(".env")
	loadDotEnvIfPresent(filepath.Join("..", ".env"))

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
		SMTPHost:          strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:          getenvInt("SMTP_PORT", 587),
		SMTPUser:          strings.TrimSpace(os.Getenv("SMTP_USER")),
		SMTPPass:          os.Getenv("SMTP_PASS"),
		SMTPFrom:          strings.TrimSpace(os.Getenv("SMTP_FROM")),
		SMTPTo:            strings.TrimSpace(os.Getenv("SMTP_TO")),
		SMTPRequireTLS:    getenvBool("SMTP_REQUIRE_TLS", true),
		MaxBodyBytes:      maxBody,
		WebDir:            getenv("WEB_DIR", "./web"),
	}
}

func loadDotEnvIfPresent(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		_ = os.Setenv(key, val)
	}
}

func getenvBool(k string, def bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(k)))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
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
