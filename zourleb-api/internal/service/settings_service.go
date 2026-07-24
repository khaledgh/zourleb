package service

import (
	"strconv"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// SettingsService reads and writes platform settings / feature flags.
type SettingsService struct {
	repo *repository.I18nRepository
}

func NewSettingsService(repo *repository.I18nRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

// Bool reads a boolean setting, defaulting when missing/unparsable.
func (s *SettingsService) Bool(key string, def bool) bool {
	st, err := s.repo.GetSetting(key)
	if err != nil {
		return def
	}
	b, err := strconv.ParseBool(st.Value)
	if err != nil {
		return def
	}
	return b
}

// String reads a string setting with a default.
func (s *SettingsService) String(key, def string) string {
	st, err := s.repo.GetSetting(key)
	if err != nil {
		return def
	}
	return st.Value
}

// ModuleFlags reports which optional modules are enabled (for clients).
func (s *SettingsService) ModuleFlags() models.ModuleFlags {
	return models.ModuleFlags{
		Shop:    s.Bool(models.SettingShopEnabled, false),
		Reviews: s.Bool(models.SettingReviewsEnabled, true),
		Boost:   s.Bool(models.SettingBoostEnabled, true),
	}
}

// All returns every setting (admin view).
func (s *SettingsService) All() ([]models.Setting, error) {
	rows, err := s.repo.AllSettings()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

// Upsert creates/updates a setting (admin).
func (s *SettingsService) Upsert(key, value, typ string) error {
	if typ == "" {
		typ = models.SettingString
	}
	if err := s.repo.UpsertSetting(&models.Setting{Key: key, Value: value, Type: typ}); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}
