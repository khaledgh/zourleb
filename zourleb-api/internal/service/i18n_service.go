package service

import (
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// I18nService serves languages and UI translation bundles.
type I18nService struct {
	repo        *repository.I18nRepository
	defaultLang string
}

func NewI18nService(repo *repository.I18nRepository, defaultLang string) *I18nService {
	return &I18nService{repo: repo, defaultLang: defaultLang}
}

// Languages returns active languages for clients.
func (s *I18nService) Languages() ([]models.Language, error) {
	langs, err := s.repo.ActiveLanguages()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return langs, nil
}

// Bundle is the UI translation payload for one locale.
type Bundle struct {
	Locale       string                       `json:"locale"`
	IsRTL        bool                         `json:"is_rtl"`
	Translations map[string]map[string]string `json:"translations"`
}

// Bundle returns the UI strings for a locale, falling back to the default
// locale when the requested one is not active.
func (s *I18nService) Bundle(locale string) (*Bundle, error) {
	active, err := s.repo.IsActiveLocale(locale)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	if !active {
		locale = s.defaultLang
	}

	langs, err := s.repo.ActiveLanguages()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	isRTL := false
	for _, l := range langs {
		if l.Code == locale {
			isRTL = l.IsRTL
			break
		}
	}

	tr, err := s.repo.TranslationsForLocale(locale)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return &Bundle{Locale: locale, IsRTL: isRTL, Translations: tr}, nil
}

// ResolveLocale returns a valid active locale for the requested code, falling
// back to the default. Used by the Locale middleware.
func (s *I18nService) ResolveLocale(requested string) string {
	if requested == "" {
		return s.defaultLang
	}
	active, err := s.repo.IsActiveLocale(requested)
	if err != nil || !active {
		return s.defaultLang
	}
	return requested
}
