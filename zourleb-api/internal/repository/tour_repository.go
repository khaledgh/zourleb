package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// TourRepository handles tour writes (agency portal) and pricing lookups.
type TourRepository struct {
	db *gorm.DB
}

func NewTourRepository(db *gorm.DB) *TourRepository {
	return &TourRepository{db: db}
}

func (r *TourRepository) FindByID(id uint) (*models.Tour, error) {
	var t models.Tour
	if err := r.db.Preload("Translations").Preload("Prices").First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TourRepository) FindDeparture(id uint) (*models.TourDeparture, error) {
	var d models.TourDeparture
	if err := r.db.First(&d, id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

// PricesForTour returns prices for a tour (optionally for a specific departure).
func (r *TourRepository) PricesForTour(tourID uint) ([]models.TourPrice, error) {
	var rows []models.TourPrice
	err := r.db.Where("tour_id = ?", tourID).Find(&rows).Error
	return rows, err
}

func (r *TourRepository) TitleFor(tourID uint, locale, def string) string {
	var tr models.TourTranslation
	if err := r.db.Where("tour_id = ? AND locale = ?", tourID, locale).First(&tr).Error; err == nil {
		return tr.Title
	}
	if err := r.db.Where("tour_id = ? AND locale = ?", tourID, def).First(&tr).Error; err == nil {
		return tr.Title
	}
	return ""
}

// --- Agency portal writes ---

func (r *TourRepository) Create(t *models.Tour) error { return r.db.Create(t).Error }
func (r *TourRepository) Save(t *models.Tour) error   { return r.db.Save(t).Error }

func (r *TourRepository) FindByIDForAgency(id, agencyID uint) (*models.Tour, error) {
	var t models.Tour
	err := r.db.Preload("Translations").Preload("Images").Preload("Departures").Preload("Prices").
		Where("id = ? AND agency_id = ?", id, agencyID).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TourRepository) ListForAgency(agencyID uint, p pagination.Params) ([]models.Tour, int64, error) {
	q := r.db.Model(&models.Tour{}).Where("agency_id = ?", agencyID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Tour
	err := q.Preload("Translations").Preload("Images", "is_cover = ?", true).
		Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *TourRepository) UpsertTranslation(tr *models.TourTranslation) error {
	return r.db.Where("tour_id = ? AND locale = ?", tr.TourID, tr.Locale).
		Assign(*tr).FirstOrCreate(tr).Error
}

func (r *TourRepository) AddImage(img *models.TourImage) error  { return r.db.Create(img).Error }
func (r *TourRepository) AddDeparture(d *models.TourDeparture) error { return r.db.Create(d).Error }
func (r *TourRepository) AddPrice(pr *models.TourPrice) error   { return r.db.Create(pr).Error }

func (r *TourRepository) DeleteImage(id, agencyID uint) error {
	return r.db.Where("id = ? AND tour_id IN (?)", id,
		r.db.Model(&models.Tour{}).Select("id").Where("agency_id = ?", agencyID)).
		Delete(&models.TourImage{}).Error
}

// UpdatePriceFrom recomputes a tour's price_from from its cheapest adult price.
func (r *TourRepository) UpdatePriceFrom(tourID uint) error {
	var min struct{ Amount float64 }
	r.db.Model(&models.TourPrice{}).
		Select("MIN(amount) as amount").
		Where("tour_id = ? AND traveler_type = ?", tourID, "adult").
		Scan(&min)
	return r.db.Model(&models.Tour{}).Where("id = ?", tourID).
		Update("price_from", min.Amount).Error
}
