package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ghassenelkamel/portfolio-v2/internal/config"
	"github.com/ghassenelkamel/portfolio-v2/internal/contact"
	"github.com/ghassenelkamel/portfolio-v2/internal/content"
	"github.com/ghassenelkamel/portfolio-v2/internal/github"
	"github.com/ghassenelkamel/portfolio-v2/internal/profile"
)

type API struct {
	cfg        config.Config
	gh         *github.Client
	profileSvc *profile.Service
	contentSvc *content.Service
	contactSvc *contact.Service
}

func NewAPI(cfg config.Config) http.Handler {
	a := &API{
		cfg:        cfg,
		gh:         github.NewClient(cfg.GitHubToken),
		profileSvc: profile.NewService(),
		contentSvc: content.NewService(cfg),
		contactSvc: contact.NewService(cfg.ContactForwardURL),
	}
	return http.HandlerFunc(a.serveHTTP)
}

func (a *API) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")

	path := strings.TrimPrefix(r.URL.Path, "/api")

	switch {
	case r.Method == http.MethodGet && path == "/healthz":
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return

	case r.Method == http.MethodGet && path == "/github-repos":
		a.githubRepos(w, r)
		return

	case r.Method == http.MethodGet && path == "/profile-summary":
		a.profileSummary(w, r)
		return

	case r.Method == http.MethodGet && path == "/content/about":
		a.contentAbout(w, r)
		return

	case r.Method == http.MethodGet && path == "/content/skills":
		a.contentSkills(w, r)
		return

	case r.Method == http.MethodGet && path == "/content/experience":
		a.contentExperience(w, r)
		return

	case r.Method == http.MethodGet && path == "/content/contact":
		a.contentContact(w, r)
		return

	case r.Method == http.MethodPost && path == "/contact":
		a.contact(w, r)
		return

	default:
		http.NotFound(w, r)
		return
	}
}

func (a *API) githubRepos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=300")

	user := r.URL.Query().Get("user")
	if user == "" {
		user = a.cfg.GitHubUser
	}

	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	repos, rate, err := a.gh.ListRepos(r.Context(), user, limit)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"repos": []any{}, "error": err.Error(), "rate": rate})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"repos": repos, "error": nil, "rate": rate})
}

func (a *API) profileSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=300")
	fallback := profile.DefaultSummary()

	defer func() {
		if rec := recover(); rec != nil {
			writeJSON(w, http.StatusOK, map[string]any{"ok": false, "summary": fallback})
		}
	}()

	sum, ok := a.profileSvc.FetchSummary(r.Context())
	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "summary": fallback})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "summary": sum})
}

func (a *API) contentAbout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=120")
	data, err := a.contentSvc.GetAbout(r.Context())
	resp := map[string]any{"ok": err == nil, "data": data}
	if err != nil {
		resp["error"] = err.Error()
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *API) contentSkills(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=120")
	data, err := a.contentSvc.GetSkills(r.Context())
	resp := map[string]any{"ok": err == nil, "data": data}
	if err != nil {
		resp["error"] = err.Error()
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *API) contentExperience(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=120")
	data, err := a.contentSvc.GetExperience(r.Context())
	resp := map[string]any{"ok": err == nil, "data": data}
	if err != nil {
		resp["error"] = err.Error()
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *API) contentContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=120")
	data, err := a.contentSvc.GetContact(r.Context())
	resp := map[string]any{"ok": err == nil, "data": data}
	if err != nil {
		resp["error"] = err.Error()
	}
	writeJSON(w, http.StatusOK, resp)
}

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (a *API) contact(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, a.cfg.MaxBodyBytes)
	defer r.Body.Close()

	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "Invalid JSON."})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Message = strings.TrimSpace(req.Message)

	if req.Message == "" || len(req.Message) < 10 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "Message is too short."})
		return
	}
	if len(req.Message) > 4000 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "Message is too long."})
		return
	}
	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "Email is required."})
		return
	}

	if err := a.contactSvc.Handle(r.Context(), req.Name, req.Email, req.Subject, req.Message); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{"ok": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
