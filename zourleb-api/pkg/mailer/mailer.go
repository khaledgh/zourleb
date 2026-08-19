// Package mailer abstracts sending transactional email.
package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Sender dispatches email messages.
type Sender interface {
	Send(ctx context.Context, to, subject, body string) error
	Configured() bool
}

// LogSender logs the message to stdout. Used in dev when SMTP is not configured.
type LogSender struct {
	log *slog.Logger
}

func NewLogSender(log *slog.Logger) *LogSender { return &LogSender{log: log} }

func (s *LogSender) Configured() bool { return true }
func (s *LogSender) Send(_ context.Context, to, subject, body string) error {
	s.log.Info("[email]", "to", to, "subject", subject, "body", body)
	return nil
}

// SMTPClient sends email through a configured SMTP relay.
type SMTPClient struct {
	host     string
	port     string
	from     string
	user     string
	password string
	dry      bool
}

func NewSMTPClient(host, port, from, user, password string) *SMTPClient {
	return &SMTPClient{
		host: host, port: port, from: from,
		user: user, password: password,
		dry: os.Getenv("SMTP_DRY_RUN") == "true",
	}
}

func (c *SMTPClient) Configured() bool {
	return c.host != "" && c.port != "" && c.from != ""
}

func (c *SMTPClient) Send(_ context.Context, to, subject, body string) error {
	if !c.Configured() {
		return fmt.Errorf("SMTP not configured")
	}
	if c.dry {
		fmt.Fprintf(os.Stderr, "[DRY-RUN] email to=%s subject=%s body=%s\n", to, subject, body)
		return nil
	}
	addr := c.host + ":" + c.port
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))
	auth := smtp.PlainAuth("", c.user, c.password, c.host)
	return smtp.SendMail(addr, auth, c.from, []string{to}, msg)
}

// Token generates a URL-safe random verification token.
func Token() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789abcdefghjkmnpqrstuvwxyz"
	b := make([]byte, 32)
	for i := range b {
		b[i] = alphabet[rand.IntN(len(alphabet))]
	}
	return string(b)
}

// LinkTemplate builds a plain-text verification message.
func LinkTemplate(baseURL, token string, ttl time.Duration) (string, string) {
	subject := "Verify your email"
	body := strings.Join([]string{
		"Hi,",
		"",
		fmt.Sprintf("Please verify your email by clicking the link below. It expires in %d minutes:", int(ttl.Minutes())),
		"",
		fmt.Sprintf("%s/verify-email?token=%s", strings.TrimRight(baseURL, "/"), token),
		"",
		"If you did not request this, you can safely ignore it.",
	}, "\n")
	return subject, body
}
