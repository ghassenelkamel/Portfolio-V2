package content

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type AboutDocument struct {
	SchemaVersion int      `json:"schemaVersion"`
	UpdatedAt     string   `json:"updatedAt"`
	Headline      string   `json:"headline"`
	Subtitle      string   `json:"subtitle"`
	Summary       []string `json:"summary"`
	Keywords      []string `json:"keywords"`
	Specialties   []string `json:"specialties"`
	Currently     []string `json:"currently"`
	Goals         []string `json:"goals"`
	Email         string   `json:"email"`
	Site          string   `json:"site"`
	GitHub        string   `json:"github"`
}

type SkillItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Order int    `json:"order"`
}

type SkillsDocument struct {
	SchemaVersion int         `json:"schemaVersion"`
	UpdatedAt     string      `json:"updatedAt"`
	Eyebrow       string      `json:"eyebrow"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Items         []SkillItem `json:"items"`
}

type ExperienceItem struct {
	ID        string   `json:"id"`
	Role      string   `json:"role"`
	Org       string   `json:"org"`
	OrgURL    string   `json:"orgUrl,omitempty"`
	Dates     string   `json:"dates"`
	Location  string   `json:"location"`
	Summary   []string `json:"summary"`
	Bullets   []string `json:"bullets"`
	Image     string   `json:"image,omitempty"`
	ProofURL  string   `json:"proofUrl,omitempty"`
	Featured  bool     `json:"featured,omitempty"`
	SortOrder int      `json:"sortOrder"`
}

type ExperienceDocument struct {
	SchemaVersion int              `json:"schemaVersion"`
	UpdatedAt     string           `json:"updatedAt"`
	Eyebrow       string           `json:"eyebrow"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Items         []ExperienceItem `json:"items"`
}

type ContactDocument struct {
	SchemaVersion int    `json:"schemaVersion"`
	UpdatedAt     string `json:"updatedAt"`
	Eyebrow       string `json:"eyebrow"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Email         string `json:"email"`
	LinkedIn      string `json:"linkedin"`
	GitHub        string `json:"github"`
	Portfolio     string `json:"portfolio"`
}

func DefaultAbout() AboutDocument {
	return AboutDocument{
		SchemaVersion: 1,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Headline:      "Ghassen El Kamel",
		Subtitle:      "IoT & Systems Engineer · Backend · Infrastructure",
		Summary: []string{
			"I’m an IoT and Software Systems Engineer passionate about building secure, scalable, and efficient infrastructures for connected devices and platforms.",
			"My work sits at the intersection of embedded systems, backend development, and network security, transforming real-world hardware into reliable digital ecosystems.",
		},
		Keywords: []string{
			"IoT architecture",
			"Go backend",
			"Embedded Linux",
			"MQTT / TLS",
			"VPN networking",
			"Automation APIs",
		},
		Specialties: []string{
			"Designing IoT architectures with secure data flow from device to cloud",
			"Developing back-end services in Go for automation, monitoring, and data processing",
			"Integrating secure communication layers — VPN, MQTTS, and TLS",
			"Optimizing workflows and systems to improve uptime and efficiency",
			"Building dashboards and tools to empower clients with actionable insights",
		},
		Currently: []string{
			"Building Go-based IoT backbones and secure APIs for data exchange",
			"Setting up VPN and MQTTS infrastructures for secure remote communication",
			"Designing modular dashboards for real-time monitoring and control",
		},
		Goals: []string{
			"Contribute to open-source infrastructure and IoT projects",
			"Advance toward architect-level expertise in backend and platform engineering",
			"Keep improving security, scalability, and maintainability in production systems",
		},
		Email:  "Ghassenelkamel@live.fr",
		Site:   "https://ghassenelkamel.fr",
		GitHub: "https://github.com/ghassenelkamel",
	}
}

func DefaultSkills() SkillsDocument {
	return SkillsDocument{
		SchemaVersion: 1,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Eyebrow:       "Skills",
		Title:         "Capability Monitor",
		Description:   "Track live proficiency signals across core tools and technologies.",
		Items: []SkillItem{
			{ID: "go", Name: "Go (Backend)", Level: 88, Order: 1},
			{ID: "mqtt-vpn-tls", Name: "MQTT / VPN / TLS", Level: 86, Order: 2},
			{ID: "linux", Name: "Linux / Unix", Level: 86, Order: 3},
			{ID: "postgresql", Name: "PostgreSQL", Level: 82, Order: 4},
			{ID: "docker", Name: "Docker", Level: 78, Order: 5},
			{ID: "svelte", Name: "JavaScript / Svelte", Level: 78, Order: 6},
			{ID: "c-cpp", Name: "C / C++", Level: 72, Order: 7},
			{ID: "python", Name: "Python", Level: 70, Order: 8},
		},
	}
}

func DefaultExperience() ExperienceDocument {
	return ExperienceDocument{
		SchemaVersion: 1,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Eyebrow:       "Experience",
		Title:         "Professional Journey",
		Description:   "Select an experience to explore scope, impact, and tools used.",
		Items: []ExperienceItem{
			{
				ID:       "spade-integrity-2024",
				Role:     "IoT Engineer",
				Org:      "Spade Integrity",
				OrgURL:   "https://spade-integrity.com/",
				Dates:    "May 2024 - Present",
				Location: "Paris, France · On-site",
				Summary: []string{
					"Contribute to client platform development with focus on reliability, security, and operational continuity.",
					"Work on secure data recovery flows, NAS integration, and synchronization systems for alarms and critical data.",
					"Support end-to-end IoT/backend delivery by improving data pipelines, monitoring quality, and secure service integrations.",
				},
				Bullets:   []string{"Go", "IoT platform", "Secure data recovery", "NAS integration", "Alarm synchronization", "Monitoring", "Infrastructure"},
				Image:     "assets/spade-integrity.svg",
				SortOrder: 10,
			},
			{
				ID:       "bookbeo-2023",
				Role:     "Full Stack Mobile Developer (Internship)",
				Org:      "bookBeo",
				OrgURL:   "https://www.bookbeo.com/fr",
				Dates:    "Mar 2023 - Sep 2023 · 7 mos",
				Location: "Rennes, Brittany, France · On-site",
				Summary: []string{
					"Developed mobile product features and backend-connected workflows during internship delivery cycles.",
					"Collaborated on release quality and implementation speed with structured team practices and version control discipline.",
				},
				Bullets:   []string{"GitLab", "Git", "React Native", "API integration", "Product delivery", "Performance", "Team collaboration"},
				Image:     "/assets/p7.jpg",
				SortOrder: 20,
			},
			{
				ID:        "st-2022",
				Role:      "Security Project Team Member (Study Project)",
				Org:       "STMicroelectronics",
				OrgURL:    "https://www.st.com",
				Dates:     "Sep 2022 – Mar 2023",
				Location:  "Le Mans, France · On-site",
				Summary:   []string{"Worked on PSA Crypto API implementation under embedded constraints and strict security requirements.", "Delivered C code with structured Git workflows and Linux-based validation tooling."},
				Bullets:   []string{"C", "PSA Crypto", "Embedded security", "Linux tooling"},
				Image:     "/assets/p6.jpg",
				SortOrder: 30,
			},
			{
				ID:        "falcon-2022",
				Role:      "Mobile Developer Intern",
				Org:       "Falcon Inspection & Services",
				OrgURL:    "https://www.falconinspec.com/",
				Dates:     "Jun 2022 – Aug 2022",
				Location:  "Tunis, Tunisia · On-site",
				Summary:   []string{"Built a mobile workflow to record and validate multiple industrial tests in the field.", "Implemented per-test data models and rule validation to improve consistency of reports."},
				Bullets:   []string{"Flutter", "Node.js", "Validation flows", "Docker"},
				Image:     "/assets/p5.jpg",
				SortOrder: 40,
			},
			{
				ID:        "dall-2021",
				Role:      "Weather Monitoring & Data Storage (Internship)",
				Org:       "DigiArtLivingLab",
				OrgURL:    "https://www.dall4all.org/",
				Dates:     "Jul 2021 – Aug 2021",
				Location:  "Tunis, Tunisia · On-site",
				Summary:   []string{"Automated weather data collection and storage pipelines for continuous monitoring.", "Built UI modules and API integrations with JSON parsing and backend persistence."},
				Bullets:   []string{"Java", "JavaFX", "API integration", "Automation"},
				Image:     "/assets/p1.jpg",
				SortOrder: 50,
			},
			{
				ID:        "sagemcom-2020",
				Role:      "Electronic Components Validation (Internship)",
				Org:       "SAGEMCOM",
				OrgURL:    "https://sagemcom.com",
				Dates:     "Feb 2020 – May 2020",
				Location:  "Tunis, Tunisia · On-site",
				Summary:   []string{"Built LabView workflows for detection and qualification of electronic components.", "Trained a CNN model and deployed a lightweight inference variant on Raspberry Pi hardware."},
				Bullets:   []string{"LabView", "CNN", "TensorFlow Lite", "Raspberry Pi", "SolidWorks"},
				Image:     "/assets/p2.jpg",
				SortOrder: 60,
			},
			{
				ID:        "dall-2019",
				Role:      "IoT Station + Mobile Apps (Internship)",
				Org:       "DigiArtLivingLab",
				OrgURL:    "https://www.dall4all.org/",
				Dates:     "Jan 2019 – Feb 2019",
				Location:  "Tunis, Tunisia · On-site",
				Summary:   []string{"Developed Arduino-based sensing builds with Bluetooth communication.", "Created Android companion apps and live data sync under tight timeline constraints."},
				Bullets:   []string{"Arduino", "Bluetooth", "MIT App Inventor", "Firebase"},
				Image:     "/assets/p3.jpg",
				SortOrder: 70,
			},
			{
				ID:        "dall-2018",
				Role:      "Educational Game Development (Internship)",
				Org:       "DigiArtLivingLab",
				OrgURL:    "https://www.dall4all.org/",
				Dates:     "Jan 2018 – Feb 2018",
				Location:  "Nabeul, Tunisia · On-site",
				Summary:   []string{"Built a Unity educational game focused on directional learning for children.", "Coordinated assets, level pacing, and gameplay polish with a small collaborative team."},
				Bullets:   []string{"Unity", "Gameplay design", "Team collaboration", "C++ basics"},
				Image:     "/assets/p4.jpg",
				SortOrder: 80,
			},
		},
	}
}

func DefaultContact() ContactDocument {
	return ContactDocument{
		SchemaVersion: 1,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Eyebrow:       "Contact",
		Title:         "Professional Reach",
		Description:   "Select a contact channel to discuss scope, impact, and collaboration.",
		Email:         "Ghassenelkamel@live.fr",
		LinkedIn:      "https://linkedin.com/in/ghassenelkamel",
		GitHub:        "https://github.com/ghassenelkamel",
		Portfolio:     "https://ghassenelkamel.fr",
	}
}

func (d AboutDocument) Validate() error {
	if d.SchemaVersion != 1 {
		return fmt.Errorf("about schemaVersion must be 1")
	}
	if err := validateISODateTime(d.UpdatedAt); err != nil {
		return err
	}
	if strings.TrimSpace(d.Headline) == "" || strings.TrimSpace(d.Subtitle) == "" {
		return fmt.Errorf("about headline/subtitle required")
	}
	if len(d.Summary) == 0 || len(d.Specialties) == 0 {
		return fmt.Errorf("about summary and specialties are required")
	}
	if len(d.Keywords) == 0 {
		return fmt.Errorf("about keywords are required")
	}
	if !looksLikeEmail(d.Email) {
		return fmt.Errorf("about email invalid")
	}
	if !looksLikeURL(d.Site) || !looksLikeURL(d.GitHub) {
		return fmt.Errorf("about links must be valid urls")
	}
	return nil
}

func (d SkillsDocument) Validate() error {
	if d.SchemaVersion != 1 {
		return fmt.Errorf("skills schemaVersion must be 1")
	}
	if err := validateISODateTime(d.UpdatedAt); err != nil {
		return err
	}
	if strings.TrimSpace(d.Eyebrow) == "" || strings.TrimSpace(d.Title) == "" {
		return fmt.Errorf("skills title fields required")
	}
	if len(d.Items) == 0 {
		return fmt.Errorf("skills items required")
	}
	ids := map[string]struct{}{}
	for _, it := range d.Items {
		if strings.TrimSpace(it.ID) == "" || strings.TrimSpace(it.Name) == "" {
			return fmt.Errorf("skills item id/name required")
		}
		if _, dup := ids[it.ID]; dup {
			return fmt.Errorf("duplicate skill id: %s", it.ID)
		}
		ids[it.ID] = struct{}{}
		if it.Level < 0 || it.Level > 100 {
			return fmt.Errorf("skill level out of range: %s", it.ID)
		}
	}
	return nil
}

func (d ExperienceDocument) Validate() error {
	if d.SchemaVersion != 1 {
		return fmt.Errorf("experience schemaVersion must be 1")
	}
	if err := validateISODateTime(d.UpdatedAt); err != nil {
		return err
	}
	if strings.TrimSpace(d.Eyebrow) == "" || strings.TrimSpace(d.Title) == "" {
		return fmt.Errorf("experience title fields required")
	}
	if len(d.Items) == 0 {
		return fmt.Errorf("experience items required")
	}
	ids := map[string]struct{}{}
	for _, it := range d.Items {
		if strings.TrimSpace(it.ID) == "" {
			return fmt.Errorf("experience id required")
		}
		if _, dup := ids[it.ID]; dup {
			return fmt.Errorf("duplicate experience id: %s", it.ID)
		}
		ids[it.ID] = struct{}{}
		if strings.TrimSpace(it.Role) == "" || strings.TrimSpace(it.Org) == "" || strings.TrimSpace(it.Dates) == "" {
			return fmt.Errorf("experience role/org/dates required")
		}
		if len(it.Summary) == 0 {
			return fmt.Errorf("experience summary required: %s", it.ID)
		}
		if it.Image != "" {
			img := strings.TrimSpace(it.Image)
			if !(strings.HasPrefix(img, "assets/") || strings.HasPrefix(img, "/assets/") || looksLikeURL(img)) {
				return fmt.Errorf("experience image must be repo asset path or absolute url: %s", it.ID)
			}
		}
		if it.ProofURL != "" && !looksLikeURL(it.ProofURL) {
			return fmt.Errorf("experience proofUrl invalid: %s", it.ID)
		}
		if it.OrgURL != "" && !looksLikeURL(it.OrgURL) {
			return fmt.Errorf("experience orgUrl invalid: %s", it.ID)
		}
	}
	return nil
}

func (d ContactDocument) Validate() error {
	if d.SchemaVersion != 1 {
		return fmt.Errorf("contact schemaVersion must be 1")
	}
	if err := validateISODateTime(d.UpdatedAt); err != nil {
		return err
	}
	if strings.TrimSpace(d.Eyebrow) == "" || strings.TrimSpace(d.Title) == "" {
		return fmt.Errorf("contact title fields required")
	}
	if !looksLikeEmail(d.Email) {
		return fmt.Errorf("contact email invalid")
	}
	if !looksLikeURL(d.LinkedIn) || !looksLikeURL(d.GitHub) || !looksLikeURL(d.Portfolio) {
		return fmt.Errorf("contact links must be valid urls")
	}
	return nil
}

func validateISODateTime(v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("updatedAt is required")
	}
	if _, err := time.Parse(time.RFC3339, v); err != nil {
		return fmt.Errorf("updatedAt must be RFC3339")
	}
	return nil
}

func looksLikeURL(v string) bool {
	u, err := url.Parse(strings.TrimSpace(v))
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func looksLikeEmail(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	parts := strings.Split(v, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}
