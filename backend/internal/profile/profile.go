package profile

import (
	"context"
	"regexp"
	"strings"
)

type Summary struct {
	Headline  string `json:"headline"`
	Subtitle  string `json:"subtitle"`
	About     string `json:"about"`
	Working   string `json:"working,omitempty"`
	Goals     string `json:"goals,omitempty"`
	Email     string `json:"email"`
	Portfolio string `json:"portfolio"`
}

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) FetchSummary(ctx context.Context) (Summary, bool) {
	md, ok := fetchReadmeRaw(ctx)
	if !ok {
		return Summary{}, false
	}

	about := extractSection(md, "💡 About Me")
	working := extractSection(md, "🧩 What I’m Working On")
	goals := extractSection(md, "🚀 Vision & Goals")

	email := pickFirstEmail(md)
	if email == "" {
		email = "Ghassenelkamel@live.fr"
	}
	portfolio := pickPortfolio(md)
	if portfolio == "" {
		portfolio = "https://ghassenelkamel.fr"
	}

	return Summary{
		Headline:  "Ghassen El Kamel",
		Subtitle:  "IoT & Systems Engineer · Backend · Infrastructure",
		About:     about,
		Working:   working,
		Goals:     goals,
		Email:     email,
		Portfolio: portfolio,
	}, true
}

// extractSection finds a "### <title>" heading and returns its body until the next heading (## or ###) or a line with "---".
func extractSection(md, title string) string {
	lines := strings.Split(md, "\n")

	// Find the heading line
	start := -1
	target := "### " + title
	for i, l := range lines {
		if strings.TrimSpace(l) == target {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}

	// Capture until next section boundary
	end := len(lines)
	for i := start; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])

		if strings.HasPrefix(t, "### ") || strings.HasPrefix(t, "## ") || t == "---" {
			end = i
			break
		}
	}

	out := strings.Join(lines[start:end], "\n")
	return strings.TrimSpace(out)
}

func pickFirstEmail(md string) string {
	re := regexp.MustCompile(`(?i)[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}`)
	return strings.TrimSpace(re.FindString(md))
}

func pickPortfolio(md string) string {
	re := regexp.MustCompile(`(?i)Portfolio:\s*\[.*?\]\((https?://[^\)]+)\)`)
	m := re.FindStringSubmatch(md)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}
