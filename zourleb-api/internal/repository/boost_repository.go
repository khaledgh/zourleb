package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
)

// BoostRepository persists boost packages and boosts.
type BoostRepository struct {
	db *gorm.DB
}

func NewBoostRepository(db *gorm.DB) *BoostRepository {
	return &BoostRepository{db: db}
}

// --- Packages ---

func (r *BoostRepository) Packages() ([]models.BoostPackage, error) {
	var rows []models.BoostPackage
	err := r.db.Order("priority desc").Find(&rows).Error
	return rows, err
}

func (r *BoostRepository) FindPackage(id uint) (*models.BoostPackage, error) {
	var p models.BoostPackage
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *BoostRepository) CreatePackage(p *models.BoostPackage) error { return r.db.Create(p).Error }
func (r *BoostRepository) SavePackage(p *models.BoostPackage) error   { return r.db.Save(p).Error }

// --- Boosts ---

func (r *BoostRepository) Create(b *models.Boost) error { return r.db.Create(b).Error }

func (r *BoostRepository) FindByID(id uint) (*models.Boost, error) {
	var b models.Boost
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BoostRepository) Save(b *models.Boost) error { return r.db.Save(b).Error }

func (r *BoostRepository) ListForAgency(agencyID uint) ([]models.Boost, error) {
	var rows []models.Boost
	err := r.db.Where("agency_id = ?", agencyID).Order("created_at desc").Find(&rows).Error
	return rows, err
}

// ActiveByPlacement returns active boosts for a placement, highest priority
// first, capped at the package's max_active_slots window.
func (r *BoostRepository) ActiveByPlacement(placement string, now time.Time, limit int) ([]models.Boost, error) {
	var rows []models.Boost
	err := r.db.
		Joins("JOIN boost_packages bp ON bp.id = boosts.package_id").
		Where("boosts.placement = ? AND boosts.status = ? AND boosts.starts_at <= ? AND boosts.ends_at >= ?",
			placement, models.BoostActive, now, now).
		Order("bp.priority DESC, boosts.starts_at DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

// IncrImpressions / IncrClicks bump counters for analytics.
func (r *BoostRepository) IncrImpressions(id uint) error {
	return r.db.Model(&models.Boost{}).Where("id = ?", id).
		UpdateColumn("impressions", gorm.Expr("impressions + 1")).Error
}

func (r *BoostRepository) IncrClicks(id uint) error {
	return r.db.Model(&models.Boost{}).Where("id = ?", id).
		UpdateColumn("clicks", gorm.Expr("clicks + 1")).Error
}

// ExpireDue marks active boosts past their end as expired; returns affected agency ids.
func (r *BoostRepository) ExpireDue(now time.Time) (int64, error) {
	res := r.db.Model(&models.Boost{}).
		Where("status = ? AND ends_at < ?", models.BoostActive, now).
		Update("status", models.BoostExpired)
	return res.RowsAffected, res.Error
}

// PendingForAdmin lists boosts awaiting approval/payment confirmation.
func (r *BoostRepository) PendingForAdmin() ([]models.Boost, error) {
	var rows []models.Boost
	err := r.db.Where("status = ?", models.BoostPendingPayment).
		Order("created_at asc").Find(&rows).Error
	return rows, err
}
