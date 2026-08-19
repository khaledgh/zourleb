// Package service provides audit logging for sensitive actions.
package service

import (
	"context"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// AuditService logs accountable actions without failing the caller.
type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// Log records an action best-effort; errors are swallowed to avoid failing ops.
func (s *AuditService) Log(ctx context.Context, actorID *uint, action, targetType string, targetID uint, payload map[string]any, ip string) {
	_ = s.repo.Create(actorID, action, targetType, targetID, payload, ip)
}

// List returns audit logs for review in the admin panel.
func (s *AuditService) List(p pagination.Params) ([]models.AuditLog, int64, error) {
	return s.repo.List(p)
}
