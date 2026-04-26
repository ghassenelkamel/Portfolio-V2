package profile

import (
	"context"
	"regexp"
	"strings"
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

const fallbackAbout = `I’m an IoT and Software Systems Engineer passionate about building secure, scalable, and efficient infrastructures for connected devices and platforms.

My work sits at the intersection of embedded systems, backend development, and network security, transforming real-world hardware into reliable digital ecosystems.

I specialize in:
- Designing IoT architectures with secure data flow from device to cloud
- Developing back-end services in Go for automation, monitoring, and data processing
- Integrating secure communication layers — VPN, MQTTS, and TLS`

const fallbackWorking = `- Building Go-based IoT backbones and secure APIs for data exchange
- Setting up VPN and MQTTS infrastructures for secure remote communication
- Designing modular dashboards for real-time monitoring and control`

const fallbackGoals = `- Contribute to open-source infrastructure and IoT projects
- Advance toward architect-level expertise in backend and platform engineering
- Keep improving security, scalability, and maintainability in production systems`

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

func DefaultSummary() Summary {
	return Summary{
		Headline:  "Ghassen El Kamel",
		Subtitle:  "IoT & Systems Engineer · Backend · Infrastructure",
		About:     fallbackAbout,
		Working:   fallbackWorking,
		Goals:     fallbackGoals,
		Email:     "Ghassenelkamel@live.fr",
		Portfolio: "https://ghassenelkamel.fr",
	}
}

func (s *Service) FetchSummary(ctx context.Context) (Summary, bool) {
	base := DefaultSummary()

	md, ok := fetchReadmeRaw(ctx)
	if !ok {
		return Summary{}, false
	}

	about := extractSectionAny(md, []string{
		"💡 About Me",
		"About Me",
		"About",
	})
	working := extractSectionAny(md, []string{
		"🧩 What I’m Working On",
		"What I’m Working On",
		"What I'm Working On",
		"Working On",
	})
	goals := extractSectionAny(md, []string{
		"🚀 Vision & Goals",
		"Vision & Goals",
		"Vision and Goals",
		"Goals",
	})

	if about == "" {
		about = base.About
	}
	if working == "" {
		working = base.Working
	}
	if goals == "" {
		goals = base.Goals
	}

	email := pickFirstEmail(md)
	if email == "" {
		email = base.Email
	}
	portfolio := pickPortfolio(md)
	if portfolio == "" {
		portfolio = base.Portfolio
	}

	return Summary{
		Headline:  base.Headline,
		Subtitle:  base.Subtitle,
		About:     about,
		Working:   working,
		Goals:     goals,
		Email:     email,
		Portfolio: portfolio,
	}, true
}

// extractSection finds a "### <title>" heading and returns its body until the next heading (## or ###) or a line with "---".
func extractSection(md, title string) string {
	return extractSectionAny(md, []string{title})
}

func extractSectionAny(md string, titles []string) string {
	lines := strings.Split(md, "\n")
	titleSet := make(map[string]struct{}, len(titles))
	for _, t := range titles {
		n := normalizeHeading(t)
		if n != "" {
			titleSet[n] = struct{}{}
		}
	}

	// Find the heading line
	start := -1
	for i, l := range lines {
		t := strings.TrimSpace(l)
		heading := ""
		switch {
		case strings.HasPrefix(t, "### "):
			heading = strings.TrimSpace(strings.TrimPrefix(t, "### "))
		case strings.HasPrefix(t, "## "):
			heading = strings.TrimSpace(strings.TrimPrefix(t, "## "))
		default:
			continue
		}

		if _, ok := titleSet[normalizeHeading(heading)]; ok {
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

func normalizeHeading(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.NewReplacer("’", "'", "`", "", "_", "", "*", "").Replace(s)
	s = nonAlnum.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
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
