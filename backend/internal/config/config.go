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
		Addr:             getenv("ADDR", ":8080"),
		GitHubUser:        getenv("GITHUB_USER", "ghassenelkamel"),
		GitHubToken:       os.Getenv("GITHUB_TOKEN"),
		ContactForwardURL: os.Getenv("CONTACT_FORWARD_URL"),
		MaxBodyBytes:      maxBody,
		WebDir:            getenv("WEB_DIR", "./web"),
	}
}

func getenv(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}
