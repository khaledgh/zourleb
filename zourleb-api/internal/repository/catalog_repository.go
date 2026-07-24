package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// CatalogRepository reads the public catalog: tours, categories, regions.
type CatalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) *CatalogRepository {
	return &CatalogRepository{db: db}
}

// publishedScope restricts to published, non-showcase tours.
func (r *CatalogRepository) publishedScope(q *gorm.DB) *gorm.DB {
	return q.Where("status = ? AND is_showcase = ?", models.TourStatusPublished, false)
}

// ListTours returns a filtered, paginated page of tours with translations,
// cover image and agency preloaded.
func (r *CatalogRepository) ListTours(f models.TourFilter, p pagination.Params) ([]models.Tour, int64, error) {
	q := r.publishedScope(r.db.Model(&models.Tour{}))

	if f.RegionID != nil {
		q = q.Where("region_id = ?", *f.RegionID)
	}
	if f.CategoryID != nil {
		q = q.Where("category_id = ?", *f.CategoryID)
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Featured != nil {
		q = q.Where("featured = ?", *f.Featured)
	}
	if f.PriceMin != nil {
		q = q.Where("price_from >= ?", *f.PriceMin)
	}
	if f.PriceMax != nil {
		q = q.Where("price_from <= ?", *f.PriceMax)
	}
	if f.DateFrom != nil {
		q = q.Where("starts_from_date >= ?", *f.DateFrom)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("id IN (?)",
			r.db.Model(&models.TourTranslation{}).
				Select("tour_id").
				Where("title LIKE ? OR summary LIKE ?", like, like),
		)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch f.Sort {
	case "price_asc":
		q = q.Order("price_from asc")
	case "price_desc":
		q = q.Order("price_from desc")
	case "featured":
		q = q.Order("featured desc, created_at desc")
	default:
		q = q.Order("created_at desc")
	}

	var tours []models.Tour
	err := q.
		Preload("Translations").
		Preload("Images", "is_cover = ?", true).
		Offset(p.Offset()).Limit(p.Limit()).
		Find(&tours).Error
	return tours, total, err
}

// AgenciesByIDs loads compact agency rows for cards (avoids N+1).
func (r *CatalogRepository) AgenciesByIDs(ids []uint) (map[uint]models.Agency, error) {
	out := map[uint]models.Agency{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []models.Agency
	if err := r.db.Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, a := range rows {
		out[a.ID] = a
	}
	return out, nil
}

// FindTourBySlug loads a full tour for the detail view.
func (r *CatalogRepository) FindTourBySlug(slug string) (*models.Tour, error) {
	var t models.Tour
	err := r.db.
		Preload("Translations").
		Preload("Images").
		Preload("Departures", func(db *gorm.DB) *gorm.DB {
			return db.Where("status <> ?", models.DepartureCancelled).Order("start_date asc")
		}).
		Preload("Prices").
		Where("slug = ? AND status = ?", slug, models.TourStatusPublished).
		First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FeaturedTours returns up to limit featured published tours for the home page.
func (r *CatalogRepository) FeaturedTours(limit int) ([]models.Tour, error) {
	var tours []models.Tour
	err := r.publishedScope(r.db.Model(&models.Tour{})).
		Where("featured = ?", true).
		Preload("Translations").
		Preload("Images", "is_cover = ?", true).
		Order("created_at desc").Limit(limit).
		Find(&tours).Error
	return tours, err
}

func (r *CatalogRepository) Categories() ([]models.Category, error) {
	var rows []models.Category
	err := r.db.Preload("Translations").Order("sort_order asc").Find(&rows).Error
	return rows, err
}

func (r *CatalogRepository) Regions() ([]models.Region, error) {
	var rows []models.Region
	err := r.db.Preload("Translations").Find(&rows).Error
	return rows, err
}

func (r *CatalogRepository) FindAgencyBySlug(slug string) (*models.Agency, error) {
	var a models.Agency
	err := r.db.Preload("Translations").
		Where("slug = ? AND status = ?", slug, models.AgencyStatusApproved).
		First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// --- Home support: banners ---

func (r *CatalogRepository) ActiveBanners() ([]models.Banner, error) {
	var rows []models.Banner
	err := r.db.Where("active = ?", true).Order("sort_order asc").Find(&rows).Error
	return rows, err
}
