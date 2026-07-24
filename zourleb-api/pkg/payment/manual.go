package payment

import "context"

// ManualProvider implements offline payment (bank transfer / cash on arrival).
// Charges start pending and are settled when an admin confirms receipt — very
// practical for the Lebanese market.
type ManualProvider struct{}

func NewManualProvider() *ManualProvider { return &ManualProvider{} }

func (p *ManualProvider) Name() string { return "manual" }

// Charge records a pending offline payment; no redirect.
func (p *ManualProvider) Charge(_ context.Context, req ChargeRequest) (*ChargeResult, error) {
	return &ChargeResult{
		ProviderRef: "manual_" + req.IdempotencyKey,
		Status:      StatusPending,
	}, nil
}

// Confirm marks the offline payment as paid (admin action).
func (p *ManualProvider) Confirm(_ context.Context, providerRef string) (*ChargeResult, error) {
	return &ChargeResult{ProviderRef: providerRef, Status: StatusPaid}, nil
}
