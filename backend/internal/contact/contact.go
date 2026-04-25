package contact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	forwardURL string
}

func NewService(forwardURL string) *Service { return &Service{forwardURL: forwardURL} }

func (s *Service) Handle(ctx context.Context, name, email, subject, message string) error {
	payload := map[string]string{
		"name":    name,
		"email":   email,
		"subject": subject,
		"message": message,
	}

	// If no forwarding configured, accept (useful for local/dev).
	if s.forwardURL == "" {
		return nil
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.forwardURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("forward failed: %s", resp.Status)
	}
	return nil
}
