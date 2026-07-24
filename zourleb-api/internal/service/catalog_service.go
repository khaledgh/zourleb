package service

import (
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/i18n"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// CatalogService serves the public catalog, localizing content per request.
type CatalogService struct {
	repo        *repository.CatalogRepository
	settings    *SettingsService
	defaultLang string
}

func NewCatalogService(repo *repository.CatalogRepository, settings *SettingsService, defaultLang string) *CatalogService {
	return &CatalogService{repo: repo, settings: settings, defaultLang: defaultLang}
}

// ListTours returns localized tour cards plus pagination meta.
func (s *CatalogService) ListTours(f models.TourFilter, p pagination.Params, locale string) ([]models.TourCard, pagination.Meta, error) {
	tours, total, err := s.repo.ListTours(f, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}

	agencyIDs := make([]uint, 0, len(tours))
	for _, t := range tours {
		agencyIDs = append(agencyIDs, t.AgencyID)
	}
	agencies, _ := s.repo.AgenciesByIDs(agencyIDs)

	cards := make([]models.TourCard, 0, len(tours))
	for _, t := range tours {
		cards = append(cards, s.toCard(t, agencies[t.AgencyID], locale))
	}
	return cards, pagination.NewMeta(p, total), nil
}

// TourDetail returns the full localized tour, splitting current vs past gallery.
func (s *CatalogService) TourDetail(slug, locale string) (*models.TourDetail, error) {
	t, err := s.repo.FindTourBySlug(slug)
	if err != nil {
		return nil, response.ErrTourNotFound
	}
	agencies, _ := s.repo.AgenciesByIDs([]uint{t.AgencyID})
	card := s.toCard(*t, agencies[t.AgencyID], locale)

	tr, _ := i18n.Resolve(t.Translations, locale, s.defaultLang)

	detail := &models.TourDetail{
		TourCard:    card,
		Description: tr.Description,
		Itinerary:   tr.Itinerary,
		Included:    tr.Included,
		Excluded:    tr.Excluded,
		MinAge:      t.MinAge,
		MaxCapacity: t.MaxCapacity,
	}
	for _, img := range t.Images {
		dto := models.TourImageDTO{URL: img.URL, SortOrder: img.SortOrder, IsCover: img.IsCover}
		if img.Source == models.ImageSourcePast {
			detail.PastGallery = append(detail.PastGallery, dto)
		} else {
			detail.Images = append(detail.Images, dto)
		}
	}
	for _, d := range t.Departures {
		detail.Departures = append(detail.Departures, models.TourDepartureDTO{
			ID: d.ID, StartDate: d.StartDate, EndDate: d.EndDate,
			Capacity: d.Capacity, SeatsLeft: d.SeatsLeft(), Status: d.Status,
		})
	}
	for _, pr := range t.Prices {
		detail.Prices = append(detail.Prices, models.TourPriceDTO{
			DepartureID: pr.DepartureID, TravelerType: pr.TravelerType,
			Currency: pr.Currency, Amount: pr.Amount,
		})
	}
	return detail, nil
}

// Categories returns localized categories.
func (s *CatalogService) Categories(locale string) ([]models.CategoryDTO, error) {
	rows, err := s.repo.Categories()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	out := make([]models.CategoryDTO, 0, len(rows))
	for _, c := range rows {
		tr, _ := i18n.Resolve(c.Translations, locale, s.defaultLang)
		out = append(out, models.CategoryDTO{
			ID: c.ID, Slug: c.Slug, Icon: c.Icon, Name: tr.Name, SortOrder: c.SortOrder,
		})
	}
	return out, nil
}

// Regions returns localized regions.
func (s *CatalogService) Regions(locale string) ([]models.RegionDTO, error) {
	rows, err := s.repo.Regions()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	out := make([]models.RegionDTO, 0, len(rows))
	for _, rg := range rows {
		tr, _ := i18n.Resolve(rg.Translations, locale, s.defaultLang)
		out = append(out, models.RegionDTO{
			ID: rg.ID, Slug: rg.Slug, Name: tr.Name, ParentID: rg.ParentID,
		})
	}
	return out, nil
}

// AgencyDetail returns a localized public agency profile.
func (s *CatalogService) AgencyDetail(slug, locale string) (*models.AgencyDetail, error) {
	a, err := s.repo.FindAgencyBySlug(slug)
	if err != nil {
		return nil, response.ErrAgencyNotFound
	}
	tr, _ := i18n.Resolve(a.Translations, locale, s.defaultLang)
	return &models.AgencyDetail{
		ID: a.ID, Slug: a.Slug, Name: a.Name, Logo: a.Logo, Cover: a.Cover,
		Description: tr.Description, About: tr.About, Phone: a.Phone,
		Website: a.Website, Verified: a.Verified,
		RatingAvg: a.RatingAvg, RatingCount: a.RatingCount,
	}, nil
}

// Home aggregates the home-screen payload.
func (s *CatalogService) Home(locale string) (*models.HomePayload, error) {
	banners, _ := s.repo.ActiveBanners()
	featured, _ := s.repo.FeaturedTours(10)
	cats, _ := s.Categories(locale)
	regions, _ := s.Regions(locale)

	agencyIDs := make([]uint, 0, len(featured))
	for _, t := range featured {
		agencyIDs = append(agencyIDs, t.AgencyID)
	}
	agencies, _ := s.repo.AgenciesByIDs(agencyIDs)

	payload := &models.HomePayload{
		Categories: cats,
		Regions:    regions,
		Modules:    s.settings.ModuleFlags(),
	}
	for _, b := range banners {
		payload.Banners = append(payload.Banners, models.BannerDTO{
			ID: b.ID, Image: b.Image, Title: b.Title, Link: b.Link, TourID: b.TourID,
		})
	}
	for _, t := range featured {
		payload.Featured = append(payload.Featured, s.toCard(t, agencies[t.AgencyID], locale))
	}
	return payload, nil
}

// toCard maps a tour + agency to a localized card.
func (s *CatalogService) toCard(t models.Tour, a models.Agency, locale string) models.TourCard {
	tr, _ := i18n.Resolve(t.Translations, locale, s.defaultLang)
	cover := ""
	for _, img := range t.Images {
		if img.IsCover {
			cover = img.URL
			break
		}
	}
	return models.TourCard{
		ID: t.ID, Slug: t.Slug, Title: tr.Title, Summary: tr.Summary,
		Cover: cover, Type: t.Type, DurationDays: t.DurationDays,
		Difficulty: t.Difficulty, PriceFrom: t.PriceFrom, Currency: t.BaseCurrency,
		Featured: t.Featured, RegionID: t.RegionID, CategoryID: t.CategoryID,
		Agency: models.AgencyMini{
			ID: a.ID, Slug: a.Slug, Name: a.Name, Logo: a.Logo, Verified: a.Verified,
		},
		RatingAvg: a.RatingAvg,
	}
}
