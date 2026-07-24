package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
)

// I18nRepository serves languages, UI translations, and settings.
type I18nRepository struct {
	db *gorm.DB
}

func NewI18nRepository(db *gorm.DB) *I18nRepository {
	return &I18nRepository{db: db}
}

// ActiveLanguages returns enabled languages ordered for display.
func (r *I18nRepository) ActiveLanguages() ([]models.Language, error) {
	var langs []models.Language
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&langs).Error
	return langs, err
}

// AllLanguages returns every language (admin view).
func (r *I18nRepository) AllLanguages() ([]models.Language, error) {
	var langs []models.Language
	err := r.db.Order("sort_order asc").Find(&langs).Error
	return langs, err
}

// DefaultLanguageCode returns the configured default locale code.
func (r *I18nRepository) DefaultLanguageCode() (string, error) {
	var lang models.Language
	if err := r.db.Where("is_default = ?", true).First(&lang).Error; err != nil {
		return "", err
	}
	return lang.Code, nil
}

// IsActiveLocale reports whether a locale code is active.
func (r *I18nRepository) IsActiveLocale(code string) (bool, error) {
	var n int64
	err := r.db.Model(&models.Language{}).
		Where("code = ? AND is_active = ?", code, true).Count(&n).Error
	return n > 0, err
}

// TranslationsForLocale returns all UI strings for a locale, namespaced.
// Shape: { namespace: { key: value } }.
func (r *I18nRepository) TranslationsForLocale(locale string) (map[string]map[string]string, error) {
	var rows []models.Translation
	if err := r.db.Where("locale = ?", locale).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := map[string]map[string]string{}
	for _, t := range rows {
		if out[t.Namespace] == nil {
			out[t.Namespace] = map[string]string{}
		}
		out[t.Namespace][t.Key] = t.Value
	}
	return out, nil
}

// UpsertTranslation creates or updates one UI string.
func (r *I18nRepository) UpsertTranslation(t *models.Translation) error {
	return r.db.Where("locale = ? AND namespace = ? AND `key` = ?", t.Locale, t.Namespace, t.Key).
		Assign(models.Translation{Value: t.Value}).
		FirstOrCreate(t).Error
}

// --- Settings ---

func (r *I18nRepository) AllSettings() ([]models.Setting, error) {
	var rows []models.Setting
	err := r.db.Find(&rows).Error
	return rows, err
}

func (r *I18nRepository) GetSetting(key string) (*models.Setting, error) {
	var s models.Setting
	if err := r.db.Where("`key` = ?", key).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *I18nRepository) UpsertSetting(s *models.Setting) error {
	return r.db.Save(s).Error
}

var _ = gorm.ErrRecordNotFound
