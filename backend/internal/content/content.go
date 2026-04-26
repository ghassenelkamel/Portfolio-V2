package content

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ghassenelkamel/portfolio-v2/internal/config"
)

type cacheEntry struct {
	raw     []byte
	expires time.Time
}

type Service struct {
	cfg    config.Config
	client *http.Client

	mu    sync.RWMutex
	cache map[string]cacheEntry
}

func NewService(cfg config.Config) *Service {
	return &Service{
		cfg:    cfg,
		client: &http.Client{Timeout: 12 * time.Second},
		cache:  map[string]cacheEntry{},
	}
}

func (s *Service) GetAbout(ctx context.Context) (AboutDocument, error) {
	fallback := DefaultAbout()
	raw, err := s.loadRaw(ctx, "about.json", func(b []byte) error {
		var d AboutDocument
		if err := decodeStrict(b, &d); err != nil {
			return err
		}
		return d.Validate()
	})
	if err != nil {
		return fallback, err
	}

	var out AboutDocument
	if err := decodeStrict(raw, &out); err != nil {
		return fallback, err
	}
	if err := out.Validate(); err != nil {
		return fallback, err
	}
	return out, nil
}

func (s *Service) GetSkills(ctx context.Context) (SkillsDocument, error) {
	fallback := DefaultSkills()
	raw, err := s.loadRaw(ctx, "skills.json", func(b []byte) error {
		var d SkillsDocument
		if err := decodeStrict(b, &d); err != nil {
			return err
		}
		return d.Validate()
	})
	if err != nil {
		return fallback, err
	}

	var out SkillsDocument
	if err := decodeStrict(raw, &out); err != nil {
		return fallback, err
	}
	if err := out.Validate(); err != nil {
		return fallback, err
	}
	return out, nil
}

func (s *Service) GetExperience(ctx context.Context) (ExperienceDocument, error) {
	fallback := DefaultExperience()
	raw, err := s.loadRaw(ctx, "experience.json", func(b []byte) error {
		var d ExperienceDocument
		if err := decodeStrict(b, &d); err != nil {
			return err
		}
		return d.Validate()
	})
	if err != nil {
		return fallback, err
	}

	var out ExperienceDocument
	if err := decodeStrict(raw, &out); err != nil {
		return fallback, err
	}
	if err := out.Validate(); err != nil {
		return fallback, err
	}
	return out, nil
}

func (s *Service) GetContact(ctx context.Context) (ContactDocument, error) {
	fallback := DefaultContact()
	raw, err := s.loadRaw(ctx, "contact.json", func(b []byte) error {
		var d ContactDocument
		if err := decodeStrict(b, &d); err != nil {
			return err
		}
		return d.Validate()
	})
	if err != nil {
		return fallback, err
	}

	var out ContactDocument
	if err := decodeStrict(raw, &out); err != nil {
		return fallback, err
	}
	if err := out.Validate(); err != nil {
		return fallback, err
	}
	return out, nil
}

func (s *Service) loadRaw(ctx context.Context, file string, validate func([]byte) error) ([]byte, error) {
	if raw, err := s.getCachedFresh(file); err == nil {
		if validate(raw) == nil {
			return raw, nil
		}
	}

	raw, ferr := s.fetchRaw(ctx, file)
	if ferr == nil {
		if validate(raw) == nil {
			s.putCache(file, raw)
			return raw, nil
		}
	}

	if raw, err := s.getCachedAny(file); err == nil {
		if validate(raw) == nil {
			return raw, nil
		}
	}

	if ferr != nil {
		return nil, ferr
	}
	return nil, errors.New("content validation failed")
}

func (s *Service) getCachedFresh(file string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.cache[file]
	if !ok || time.Now().After(entry.expires) || len(entry.raw) == 0 {
		return nil, errors.New("cache miss")
	}
	return append([]byte(nil), entry.raw...), nil
}

func (s *Service) getCachedAny(file string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.cache[file]
	if !ok || len(entry.raw) == 0 {
		return nil, errors.New("cache empty")
	}
	return append([]byte(nil), entry.raw...), nil
}

func (s *Service) putCache(file string, raw []byte) {
	ttl := s.cfg.ContentCacheTTL
	if ttl <= 0 {
		ttl = 300
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[file] = cacheEntry{
		raw:     append([]byte(nil), raw...),
		expires: time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

func (s *Service) fetchRaw(ctx context.Context, file string) ([]byte, error) {
	u := s.rawURL(file)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "portfolio-v2-content-fetcher")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("content fetch failed: %s", resp.Status)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return nil, errors.New("content file is empty")
	}

	return b, nil
}

func (s *Service) rawURL(file string) string {
	p := file
	if dir := strings.Trim(s.cfg.ContentRepoDir, "/"); dir != "" {
		p = path.Join(dir, file)
	}

	owner := strings.TrimSpace(s.cfg.ContentRepoOwner)
	if owner == "" {
		owner = "ghassenelkamel"
	}
	repo := strings.TrimSpace(s.cfg.ContentRepoName)
	if repo == "" {
		repo = "portfolio-content"
	}
	ref := strings.TrimSpace(s.cfg.ContentRepoRef)
	if ref == "" {
		ref = "main"
	}

	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, p)
}

func decodeStrict(raw []byte, out any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}
