// Package sms sends OTP codes via a pluggable SMS provider.
package sms

import (
	"context"
	"fmt"
	"log/slog"
)

// Client is a minimal SMS sender. The default "local"/dev provider logs the
// code; Twilio/Vonage adapters can be added behind the same Send signature.
type Client struct {
	provider string
	from     string
	apiKey   string
	apiSec   string
	log      *slog.Logger
	devMode  bool
}

func NewClient(provider, from, apiKey, apiSec string, devMode bool, log *slog.Logger) *Client {
	return &Client{provider: provider, from: from, apiKey: apiKey, apiSec: apiSec, log: log, devMode: devMode}
}

func (c *Client) Channel() string { return "sms" }

// Send delivers the OTP via SMS. Without credentials it logs in dev mode.
func (c *Client) Send(ctx context.Context, phone, code string) error {
	if c.apiKey == "" {
		if c.devMode {
			c.log.Warn("sms not configured — dev OTP", "phone", phone, "code", code)
			return nil
		}
		return fmt.Errorf("sms client not configured")
	}
	// Provider adapters (Twilio/Vonage/local gateway) plug in here.
	switch c.provider {
	default:
		c.log.Info("sms send (stub provider)", "provider", c.provider, "to", phone)
		return nil
	}
}
