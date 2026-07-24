// Package oauth verifies third-party identity tokens (Google sign-in).
package oauth

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

// GoogleProfile is the subset of verified claims we consume.
type GoogleProfile struct {
	Sub           string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

// GoogleVerifier validates Google ID tokens against the expected audience.
type GoogleVerifier struct {
	clientID string
}

func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{clientID: clientID}
}

// Configured reports whether a Google client ID was provided.
func (g *GoogleVerifier) Configured() bool { return g.clientID != "" }

// Verify validates the token signature, expiry, and audience, returning the
// profile. It performs a network fetch of Google's public keys (cached by the
// idtoken library).
func (g *GoogleVerifier) Verify(ctx context.Context, rawIDToken string) (*GoogleProfile, error) {
	if g.clientID == "" {
		return nil, errors.New("google client id not configured")
	}
	payload, err := idtoken.Validate(ctx, rawIDToken, g.clientID)
	if err != nil {
		return nil, err
	}

	profile := &GoogleProfile{Sub: payload.Subject}
	if v, ok := payload.Claims["email"].(string); ok {
		profile.Email = v
	}
	if v, ok := payload.Claims["email_verified"].(bool); ok {
		profile.EmailVerified = v
	}
	if v, ok := payload.Claims["name"].(string); ok {
		profile.Name = v
	}
	if v, ok := payload.Claims["picture"].(string); ok {
		profile.Picture = v
	}
	if profile.Email == "" {
		return nil, errors.New("google token missing email")
	}
	return profile, nil
}
