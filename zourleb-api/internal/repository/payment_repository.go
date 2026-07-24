package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
)

// PaymentRepository persists the polymorphic payments ledger.
type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(p *models.Payment) error { return r.db.Create(p).Error }

func (r *PaymentRepository) FindByID(id uint) (*models.Payment, error) {
	var p models.Payment
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByIdempotencyKey supports safe retries on money paths.
func (r *PaymentRepository) FindByIdempotencyKey(key string) (*models.Payment, error) {
	var p models.Payment
	if err := r.db.Where("idempotency_key = ?", key).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) UpdateStatus(id uint, status, providerRef string) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "provider_ref": providerRef}).Error
}

// RevenueSummary aggregates paid revenue grouped by payable type (admin).
func (r *PaymentRepository) RevenueSummary() (map[string]float64, error) {
	type row struct {
		PayableType string
		Total       float64
	}
	var rows []row
	err := r.db.Model(&models.Payment{}).
		Select("payable_type, SUM(amount) as total").
		Where("status = ?", models.PayStatusPaid).
		Group("payable_type").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := map[string]float64{}
	for _, x := range rows {
		out[x.PayableType] = x.Total
	}
	return out, nil
}
