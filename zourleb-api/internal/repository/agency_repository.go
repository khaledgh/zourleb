package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// AgencyRepository handles agency profile reads/writes and applications.
type AgencyRepository struct {
	db *gorm.DB
}

func NewAgencyRepository(db *gorm.DB) *AgencyRepository {
	return &AgencyRepository{db: db}
}

func (r *AgencyRepository) Create(a *models.Agency) error { return r.db.Create(a).Error }

func (r *AgencyRepository) FindByID(id uint) (*models.Agency, error) {
	var a models.Agency
	if err := r.db.Preload("Translations").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AgencyRepository) Save(a *models.Agency) error { return r.db.Save(a).Error }

func (r *AgencyRepository) UpsertTranslation(t *models.AgencyTranslation) error {
	return r.db.Where("agency_id = ? AND locale = ?", t.AgencyID, t.Locale).
		Assign(*t).FirstOrCreate(t).Error
}

// AddMember links a user to an agency with a role.
func (r *AgencyRepository) AddMember(m *models.AgencyMember) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "agency_id"}, {Name: "user_id"}},
		DoNothing: true,
	}).Create(m).Error
}

func (r *AgencyRepository) Members(agencyID uint) ([]models.AgencyMember, error) {
	var rows []models.AgencyMember
	err := r.db.Where("agency_id = ?", agencyID).Find(&rows).Error
	return rows, err
}

// --- Admin views ---

func (r *AgencyRepository) ListByStatus(status string, p pagination.Params) ([]models.Agency, int64, error) {
	q := r.db.Model(&models.Agency{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Agency
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *AgencyRepository) UpdateStatus(id uint, status string, verified *bool) error {
	updates := map[string]any{"status": status}
	if verified != nil {
		updates["verified"] = *verified
	}
	return r.db.Model(&models.Agency{}).Where("id = ?", id).Updates(updates).Error
}
