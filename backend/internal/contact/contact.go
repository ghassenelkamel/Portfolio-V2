package contact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"crypto/tls"
)

type SMTPConfig struct {
	Host       string
	Port       int
	User       string
	Pass       string
	From       string
	To         string
	RequireTLS bool
}

type Service struct {
	forwardURL string
	smtp       SMTPConfig
}

func NewService(forwardURL string, smtpCfg SMTPConfig) *Service {
	return &Service{forwardURL: strings.TrimSpace(forwardURL), smtp: smtpCfg}
}

func (s *Service) Handle(ctx context.Context, name, email, subject, message string) error {
	if s.smtpAnyConfigured() && !s.smtpConfigured() {
		return fmt.Errorf("smtp config incomplete: SMTP_HOST, SMTP_FROM, and SMTP_TO are required")
	}

	if s.smtpConfigured() {
		return s.sendSMTP(ctx, name, email, subject, message)
	}

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

func (s *Service) smtpConfigured() bool {
	return strings.TrimSpace(s.smtp.Host) != "" && strings.TrimSpace(s.smtp.From) != "" && strings.TrimSpace(s.smtp.To) != ""
}

func (s *Service) smtpAnyConfigured() bool {
	return strings.TrimSpace(s.smtp.Host) != "" ||
		strings.TrimSpace(s.smtp.From) != "" ||
		strings.TrimSpace(s.smtp.To) != "" ||
		strings.TrimSpace(s.smtp.User) != "" ||
		s.smtp.Pass != ""
}

func (s *Service) sendSMTP(ctx context.Context, name, email, subject, message string) error {
	if s.smtp.Port <= 0 {
		s.smtp.Port = 587
	}

	addr := net.JoinHostPort(strings.TrimSpace(s.smtp.Host), strconv.Itoa(s.smtp.Port))
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial failed: %w", err)
	}
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	client, err := smtp.NewClient(conn, s.smtp.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp client init failed: %w", err)
	}
	defer func() {
		_ = client.Close()
	}()

	startedTLS := false
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: s.smtp.Host, MinVersion: tls.VersionTLS12}
		if err := client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("smtp starttls failed: %w", err)
		}
		startedTLS = true
	} else if s.smtp.RequireTLS {
		return fmt.Errorf("smtp server does not advertise STARTTLS")
	}

	if strings.TrimSpace(s.smtp.User) != "" || s.smtp.Pass != "" {
		if !startedTLS && s.smtp.RequireTLS {
			return fmt.Errorf("refusing smtp auth without tls")
		}
		auth := smtp.PlainAuth("", s.smtp.User, s.smtp.Pass, s.smtp.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
	}

	from := strings.TrimSpace(s.smtp.From)
	to := strings.TrimSpace(s.smtp.To)
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}

	msg := buildMessage(from, to, name, email, subject, message)
	if _, err := w.Write([]byte(msg)); err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp finalize message failed: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("smtp quit failed: %w", err)
	}

	return nil
}

func buildMessage(from, to, name, email, subject, message string) string {
	n := sanitizeHeader(name)
	if n == "" {
		n = "Website Visitor"
	}
	reply := sanitizeHeader(email)
	sub := sanitizeHeader(subject)
	if sub == "" {
		sub = "New portfolio contact message"
	}
	body := strings.TrimSpace(message)

	lines := []string{
		"From: Portfolio Contact <" + sanitizeHeader(from) + ">",
		"To: <" + sanitizeHeader(to) + ">",
		"Subject: " + sub,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Date: " + time.Now().UTC().Format(time.RFC1123Z),
	}
	if reply != "" {
		lines = append(lines, "Reply-To: <"+reply+">")
	}

	headers := strings.Join(lines, "\r\n")
	bodyBlock := "Name: " + n + "\n" +
		"Email: " + reply + "\n\n" +
		"Message:\n" + body + "\n"

	return headers + "\r\n\r\n" + bodyBlock
}

func sanitizeHeader(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "\r", " ")
	v = strings.ReplaceAll(v, "\n", " ")
	return v
}
