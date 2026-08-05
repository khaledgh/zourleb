package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// AdminRepository covers platform-management entities: languages, banners,
// users, categories and regions.
type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// --- Languages ---

func (r *AdminRepository) AllLanguages() ([]models.Language, error) {
	var rows []models.Language
	err := r.db.Order("sort_order asc").Find(&rows).Error
	return rows, err
}

func (r *AdminRepository) CreateLanguage(l *models.Language) error { return r.db.Create(l).Error }

func (r *AdminRepository) SaveLanguage(l *models.Language) error { return r.db.Save(l).Error }

func (r *AdminRepository) FindLanguage(id uint) (*models.Language, error) {
	var l models.Language
	if err := r.db.First(&l, id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

// --- Banners ---

func (r *AdminRepository) AllBanners() ([]models.Banner, error) {
	var rows []models.Banner
	err := r.db.Order("sort_order asc").Find(&rows).Error
	return rows, err
}

func (r *AdminRepository) CreateBanner(b *models.Banner) error { return r.db.Create(b).Error }
func (r *AdminRepository) SaveBanner(b *models.Banner) error   { return r.db.Save(b).Error }
func (r *AdminRepository) DeleteBanner(id uint) error {
	return r.db.Delete(&models.Banner{}, id).Error
}

// --- Users ---

func (r *AdminRepository) ListUsers(p pagination.Params, search string) ([]models.User, int64, error) {
	q := r.db.Model(&models.User{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR email LIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.User
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *AdminRepository) SetUserStatus(id uint, status string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}

// --- Categories & Regions ---

func (r *AdminRepository) CreateCategory(c *models.Category) error { return r.db.Create(c).Error }
func (r *AdminRepository) UpsertCategoryTranslation(t *models.CategoryTranslation) error {
	return r.db.Where("category_id = ? AND locale = ?", t.CategoryID, t.Locale).
		Assign(*t).FirstOrCreate(t).Error
}

func (r *AdminRepository) CreateRegion(rg *models.Region) error { return r.db.Create(rg).Error }
func (r *AdminRepository) UpsertRegionTranslation(t *models.RegionTranslation) error {
	return r.db.Where("region_id = ? AND locale = ?", t.RegionID, t.Locale).
		Assign(*t).FirstOrCreate(t).Error
}

// AgencyIDForTour returns the owning agency id for a tour.
func (r *AdminRepository) AgencyIDForTour(tourID uint) (uint, error) {
	var t models.Tour
	if err := r.db.Select("agency_id").First(&t, tourID).Error; err != nil {
		return 0, err
	}
	return t.AgencyID, nil
}

// --- Counts for analytics ---

func (r *AdminRepository) Count(model interface{}, where ...interface{}) (int64, error) {
	var n int64
	q := r.db.Model(model)
	if len(where) == 2 {
		q = q.Where(where[0], where[1])
	}
	err := q.Count(&n).Error
	return n, err
}

func (r *AdminRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *AdminRepository) ClearUserRoles(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error
}

func (r *AdminRepository) DeleteAgency(id uint) error {
	return r.db.Delete(&models.Agency{}, id).Error
}

func (r *AdminRepository) ListTours(p pagination.Params, search, status string) ([]models.Tour, int64, error) {
	q := r.db.Model(&models.Tour{}).Preload("Translations")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("id IN (?) OR slug LIKE ?",
			r.db.Model(&models.TourTranslation{}).
				Select("tour_id").
				Where("title LIKE ? OR summary LIKE ?", like, like),
			like,
		)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Tour
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}
