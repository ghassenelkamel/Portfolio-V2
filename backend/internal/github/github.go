package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	token string

	mu      sync.Mutex
	cacheAt time.Time
	cacheK  string
	cacheV  []Repo
	cacheR  Rate
}

type Rate struct {
	Limit     string `json:"limit"`
	Remaining string `json:"remaining"`
	Reset     string `json:"reset"`
}

type Repo struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	HTMLURL     string     `json:"html_url"`
	Description *string    `json:"description"`
	Language    *string    `json:"language"`
	Fork        bool       `json:"fork"`
	Archived    bool       `json:"archived"`
	Stars       int        `json:"stargazers_count"`
	Forks       int        `json:"forks_count"`
	PushedAt    *time.Time `json:"pushed_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func NewClient(token string) *Client { return &Client{token: token} }

func (c *Client) ListRepos(ctx context.Context, user string, limit int) ([]Repo, Rate, error) {
	key := user + ":" + strconv.Itoa(limit)

	c.mu.Lock()
	if c.cacheK == key && time.Since(c.cacheAt) < 5*time.Minute && len(c.cacheV) > 0 {
		out := make([]Repo, len(c.cacheV))
		copy(out, c.cacheV)
		rate := c.cacheR
		c.mu.Unlock()
		return out, rate, nil
	}
	c.mu.Unlock()

	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=%d&sort=pushed", user, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, Rate{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "portfolio-v2")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, Rate{}, err
	}
	defer resp.Body.Close()

	rate := Rate{
		Limit:     resp.Header.Get("x-ratelimit-limit"),
		Remaining: resp.Header.Get("x-ratelimit-remaining"),
		Reset:     resp.Header.Get("x-ratelimit-reset"),
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32*1024))
		return nil, rate, fmt.Errorf("github api error: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, rate, err
	}

	filtered := repos[:0]
	for _, r := range repos {
		if r.Fork || r.Archived {
			continue
		}
		filtered = append(filtered, r)
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		ai := timeOrZero(filtered[i].PushedAt, filtered[i].UpdatedAt)
		aj := timeOrZero(filtered[j].PushedAt, filtered[j].UpdatedAt)
		return aj.After(ai)
	})

	c.mu.Lock()
	c.cacheAt = time.Now()
	c.cacheK = key
	c.cacheV = append([]Repo(nil), filtered...)
	c.cacheR = rate
	c.mu.Unlock()

	return filtered, rate, nil
}

func timeOrZero(a, b *time.Time) time.Time {
	if a != nil {
		return *a
	}
	if b != nil {
		return *b
	}
	return time.Time{}
}
