// Package onesignal sends push notifications via the OneSignal REST API.
package onesignal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Client posts notifications to OneSignal. When unconfigured it no-ops (logs),
// so the rest of the system runs without push credentials in dev.
type Client struct {
	appID  string
	apiKey string
	http   *http.Client
	log    *slog.Logger
}

func NewClient(appID, apiKey string, log *slog.Logger) *Client {
	return &Client{appID: appID, apiKey: apiKey, http: &http.Client{Timeout: 10 * time.Second}, log: log}
}

func (c *Client) Configured() bool { return c.appID != "" && c.apiKey != "" }

// Notification is a localized push payload targeted at player ids.
type Notification struct {
	PlayerIDs []string
	Title     string
	Body      string
	Data      map[string]any
}

// Send delivers a push to the given player ids.
func (c *Client) Send(ctx context.Context, n Notification) error {
	if !c.Configured() {
		c.log.Debug("onesignal not configured — skipping push", "title", n.Title, "targets", len(n.PlayerIDs))
		return nil
	}
	if len(n.PlayerIDs) == 0 {
		return nil
	}
	payload := map[string]any{
		"app_id":             c.appID,
		"include_player_ids": n.PlayerIDs,
		"headings":           map[string]string{"en": n.Title},
		"contents":           map[string]string{"en": n.Body},
		"data":               n.Data,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://onesignal.com/api/v1/notifications", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("onesignal send failed: status %d", resp.StatusCode)
	}
	return nil
}
