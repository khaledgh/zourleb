// Package whatsapp sends OTP codes via the WhatsApp Cloud API (Meta).
package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Client sends WhatsApp template messages. When unconfigured it logs the code
// in development instead of failing (so the OTP flow is testable offline).
type Client struct {
	phoneID  string
	token    string
	template string
	http     *http.Client
	log      *slog.Logger
	devMode  bool
}

func NewClient(phoneID, token, template string, devMode bool, log *slog.Logger) *Client {
	return &Client{
		phoneID:  phoneID,
		token:    token,
		template: template,
		http:     &http.Client{Timeout: 10 * time.Second},
		log:      log,
		devMode:  devMode,
	}
}

func (c *Client) Channel() string { return "whatsapp" }

func (c *Client) configured() bool { return c.phoneID != "" && c.token != "" }

// Send delivers the OTP via a pre-approved template with the code as a body
// parameter. Falls back to logging in dev when not configured.
func (c *Client) Send(ctx context.Context, phone, code string) error {
	if !c.configured() {
		if c.devMode {
			c.log.Warn("whatsapp not configured — dev OTP", "phone", phone, "code", code)
			return nil
		}
		return fmt.Errorf("whatsapp client not configured")
	}

	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                strings.TrimPrefix(phone, "+"),
		"type":              "template",
		"template": map[string]any{
			"name":     c.template,
			"language": map[string]string{"code": "en"},
			"components": []any{
				map[string]any{
					"type": "body",
					"parameters": []any{
						map[string]string{"type": "text", "text": code},
					},
				},
			},
		},
	}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/messages", c.phoneID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("whatsapp send failed: status %d", resp.StatusCode)
	}
	return nil
}
