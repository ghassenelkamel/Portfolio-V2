package profile

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

func fetchReadmeRaw(ctx context.Context) (string, bool) {
	urls := []string{
		"https://raw.githubusercontent.com/ghassenelkamel/ghassenelkamel/main/README.md",
		"https://raw.githubusercontent.com/ghassenelkamel/ghassenelkamel/master/README.md",
	}

	client := &http.Client{Timeout: 15 * time.Second}

	for _, u := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "portfolio-v2-profile-fetcher")
		req.Header.Set("Accept", "text/plain")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			md := strings.TrimSpace(string(b))
			if md != "" {
				return md, true
			}
		}
	}
	return "", false
}
