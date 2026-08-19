// Package repository persists audit logs.
package repository

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// AuditRepository writes audit rows to the database.
type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// List returns a paginated list of audit logs ordered by most recent first.
func (r *AuditRepository) List(p pagination.Params) ([]models.AuditLog, int64, error) {
	q := r.db.Model(&models.AuditLog{}).Order("created_at desc")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.AuditLog
	err := q.Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *AuditRepository) Create(actorID *uint, action, targetType string, targetID uint, payload map[string]any, ip string) error {
	var data datatypes.JSON
	if payload != nil {
		if b, err := json.Marshal(payload); err == nil {
			data = b
		}
	}
	return r.db.Create(&models.AuditLog{
		ActorID:    actorID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Payload:    data,
		IP:         ip,
	}).Error
}
