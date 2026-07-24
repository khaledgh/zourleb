// Package payment abstracts payment providers behind a single interface so
// local (Whish/OMT/Areeba), international (Stripe), and manual/offline flows
// are interchangeable.
package payment

import (
	"context"
	"errors"
)

// Status reflects a charge's lifecycle.
type Status string

const (
	StatusPending  Status = "pending"
	StatusPaid     Status = "paid"
	StatusFailed   Status = "failed"
	StatusRefunded Status = "refunded"
)

// ChargeRequest is a provider-agnostic charge.
type ChargeRequest struct {
	Amount         float64
	Currency       string
	Description    string
	IdempotencyKey string
	Reference      string // payable reference (booking code, boost id, …)
	ReturnURL      string
}

// ChargeResult is what a provider returns after initiating a charge.
type ChargeResult struct {
	ProviderRef string
	Status      Status
	RedirectURL string // for hosted-checkout providers
}

// Provider is implemented by each payment adapter.
type Provider interface {
	Name() string
	Charge(ctx context.Context, req ChargeRequest) (*ChargeResult, error)
	// Confirm is used by manual/offline providers (admin marks paid) and by
	// webhook handlers to settle a previously pending charge.
	Confirm(ctx context.Context, providerRef string) (*ChargeResult, error)
}

// Registry resolves providers by name.
type Registry struct {
	providers map[string]Provider
	def       string
}

func NewRegistry(def string, providers ...Provider) *Registry {
	m := make(map[string]Provider, len(providers))
	for _, p := range providers {
		m[p.Name()] = p
	}
	return &Registry{providers: m, def: def}
}

var ErrUnknownProvider = errors.New("unknown payment provider")

// Get returns a provider by name, or the default when name is empty.
func (r *Registry) Get(name string) (Provider, error) {
	if name == "" {
		name = r.def
	}
	p, ok := r.providers[name]
	if !ok {
		return nil, ErrUnknownProvider
	}
	return p, nil
}

// Names lists registered provider names.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.providers))
	for n := range r.providers {
		out = append(out, n)
	}
	return out
}
