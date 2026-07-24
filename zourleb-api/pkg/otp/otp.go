// Package otp generates numeric codes and dispatches them over a channel.
package otp

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
)

// Sender delivers an OTP code to a phone number. Implementations: WhatsApp, SMS.
type Sender interface {
	Send(ctx context.Context, phone, code string) error
	Channel() string
}

// Generate returns a cryptographically random numeric code of the given length.
func Generate(length int) (string, error) {
	if length < 4 {
		length = 4
	}
	var b strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b.WriteByte(byte('0' + n.Int64()))
	}
	return b.String(), nil
}

// Dispatcher routes a code to the sender matching the requested channel.
type Dispatcher struct {
	senders map[string]Sender
}

func NewDispatcher(senders ...Sender) *Dispatcher {
	m := make(map[string]Sender, len(senders))
	for _, s := range senders {
		if s != nil {
			m[s.Channel()] = s
		}
	}
	return &Dispatcher{senders: m}
}

// Send delivers via the named channel, returning whether a sender existed.
func (d *Dispatcher) Send(ctx context.Context, channel, phone, code string) (bool, error) {
	s, ok := d.senders[channel]
	if !ok {
		return false, nil
	}
	return true, s.Send(ctx, phone, code)
}
