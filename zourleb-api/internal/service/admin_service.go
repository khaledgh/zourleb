package service

import (
	"context"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// AdminService implements the super-admin surface: approvals, languages,
// translations, settings, banners, categories/regions, review moderation,
// payment confirmation and analytics.
type AdminService struct {
	admin    *repository.AdminRepository
	agencies *repository.AgencyRepository
	i18n     *repository.I18nRepository
	engage   *repository.EngagementRepository
	boosts   *repository.BoostRepository
	payments *repository.PaymentRepository
	settings *SettingsService
	boostSvc *BoostService
}

func NewAdminService(
	admin *repository.AdminRepository,
	agencies *repository.AgencyRepository,
	i18n *repository.I18nRepository,
	engage *repository.EngagementRepository,
	boosts *repository.BoostRepository,
	payments *repository.PaymentRepository,
	settings *SettingsService,
	boostSvc *BoostService,
) *AdminService {
	return &AdminService{
		admin: admin, agencies: agencies, i18n: i18n, engage: engage,
		boosts: boosts, payments: payments, settings: settings, boostSvc: boostSvc,
	}
}

// --- Agencies ---

func (s *AdminService) ListAgencies(status string, p pagination.Params) ([]models.Agency, pagination.Meta, error) {
	rows, total, err := s.agencies.ListByStatus(status, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

func (s *AdminService) ApproveAgency(id uint, approve bool) error {
	status := models.AgencyStatusApproved
	verified := approve
	if !approve {
		status = models.AgencyStatusSuspended
	}
	if err := s.agencies.UpdateStatus(id, status, &verified); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}

// --- Languages ---

func (s *AdminService) Languages() ([]models.Language, error) {
	rows, err := s.admin.AllLanguages()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

func (s *AdminService) CreateLanguage(req models.SaveLanguageRequest) (*models.Language, error) {
	l := &models.Language{
		Code: req.Code, Name: req.Name, NativeName: req.NativeName,
		IsRTL: req.IsRTL, IsActive: req.IsActive, SortOrder: req.SortOrder,
	}
	if err := s.admin.CreateLanguage(l); err != nil {
		return nil, response.ErrConflict.WithMessage("Language code already exists.")
	}
	return l, nil
}

func (s *AdminService) UpdateLanguage(id uint, req models.SaveLanguageRequest) (*models.Language, error) {
	l, err := s.admin.FindLanguage(id)
	if err != nil {
		return nil, response.ErrNotFound
	}
	l.Name, l.NativeName = req.Name, req.NativeName
	l.IsRTL, l.IsActive, l.SortOrder = req.IsRTL, req.IsActive, req.SortOrder
	if err := s.admin.SaveLanguage(l); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return l, nil
}

// --- Translations ---

func (s *AdminService) SaveTranslation(req models.SaveTranslationRequest) error {
	return wrapInternal(s.i18n.UpsertTranslation(&models.Translation{
		Locale: req.Locale, Namespace: req.Namespace, Key: req.Key, Value: req.Value,
	}))
}

func (s *AdminService) BulkTranslations(req models.BulkTranslationsRequest) error {
	for _, it := range req.Items {
		if err := s.i18n.UpsertTranslation(&models.Translation{
			Locale: it.Locale, Namespace: it.Namespace, Key: it.Key, Value: it.Value,
		}); err != nil {
			return response.ErrInternal.Wrap(err)
		}
	}
	return nil
}

// --- Settings ---

func (s *AdminService) Settings() ([]models.Setting, error) {
	return s.settings.All()
}

func (s *AdminService) SaveSetting(req models.SaveSettingRequest) error {
	return s.settings.Upsert(req.Key, req.Value, req.Type)
}

// --- Banners ---

func (s *AdminService) Banners() ([]models.Banner, error) {
	rows, err := s.admin.AllBanners()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

func (s *AdminService) CreateBanner(req models.SaveBannerRequest) (*models.Banner, error) {
	typ := req.Type
	if typ == "" {
		typ = models.BannerManual
	}
	b := &models.Banner{
		Type: typ, TourID: req.TourID, Image: req.Image, Title: req.Title,
		Link: req.Link, SortOrder: req.SortOrder, Active: req.Active,
	}
	if err := s.admin.CreateBanner(b); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return b, nil
}

func (s *AdminService) DeleteBanner(id uint) error {
	return wrapInternal(s.admin.DeleteBanner(id))
}

// --- Categories & Regions ---

func (s *AdminService) CreateCategory(req models.SaveCategoryRequest) (*models.Category, error) {
	c := &models.Category{Slug: req.Slug, Icon: req.Icon, SortOrder: req.SortOrder}
	if err := s.admin.CreateCategory(c); err != nil {
		return nil, response.ErrConflict.WithMessage("Category slug already exists.")
	}
	for _, t := range req.Translations {
		_ = s.admin.UpsertCategoryTranslation(&models.CategoryTranslation{
			CategoryID: c.ID, Locale: t.Locale, Name: t.Name,
		})
	}
	return c, nil
}

func (s *AdminService) CreateRegion(req models.SaveRegionRequest) (*models.Region, error) {
	rg := &models.Region{Slug: req.Slug, ParentID: req.ParentID}
	if err := s.admin.CreateRegion(rg); err != nil {
		return nil, response.ErrConflict.WithMessage("Region slug already exists.")
	}
	for _, t := range req.Translations {
		_ = s.admin.UpsertRegionTranslation(&models.RegionTranslation{
			RegionID: rg.ID, Locale: t.Locale, Name: t.Name,
		})
	}
	return rg, nil
}

// --- Reviews ---

func (s *AdminService) PendingReviews(p pagination.Params) ([]models.Review, pagination.Meta, error) {
	rows, total, err := s.engage.PendingReviews(p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

func (s *AdminService) ModerateReview(id uint, status string) error {
	rv, err := s.engage.ModerateReview(id, status)
	if err != nil {
		return response.ErrNotFound
	}
	// recompute the agency rating for the reviewed tour's agency (best-effort)
	if agencyID, err := s.admin.AgencyIDForTour(rv.TourID); err == nil {
		_ = s.engage.RecalcAgencyRating(agencyID)
	}
	return nil
}

// --- Users ---

func (s *AdminService) ListUsers(p pagination.Params, search string) ([]models.User, pagination.Meta, error) {
	rows, total, err := s.admin.ListUsers(p, search)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

func (s *AdminService) SetUserStatus(id uint, status string) error {
	return wrapInternal(s.admin.SetUserStatus(id, status))
}

// --- Boost approvals / payment confirmation ---

func (s *AdminService) PendingBoosts() ([]models.Boost, error) {
	rows, err := s.boosts.PendingForAdmin()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

// ConfirmPayment settles a pending (manual/offline) payment and, when it backs a
// boost, activates that boost.
func (s *AdminService) ConfirmPayment(ctx context.Context, paymentID uint) error {
	pay, err := s.payments.FindByID(paymentID)
	if err != nil {
		return response.ErrNotFound
	}
	if pay.Status == models.PayStatusPaid {
		return nil
	}
	if err := s.payments.UpdateStatus(pay.ID, models.PayStatusPaid, pay.ProviderRef); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	if pay.PayableType == models.PayableBoost {
		if err := s.boostSvc.Activate(pay.PayableID, pay.ID); err != nil {
			return err
		}
	}
	return nil
}

// --- Analytics ---

func (s *AdminService) Analytics() (*models.AnalyticsSummary, error) {
	users, _ := s.admin.Count(&models.User{})
	agencies, _ := s.admin.Count(&models.Agency{})
	tours, _ := s.admin.Count(&models.Tour{})
	bookings, _ := s.admin.Count(&models.Booking{})
	activeBoosts, _ := s.admin.Count(&models.Boost{}, "status = ?", models.BoostActive)
	revenue, _ := s.payments.RevenueSummary()
	return &models.AnalyticsSummary{
		Users: users, Agencies: agencies, Tours: tours, Bookings: bookings,
		ActiveBoosts: activeBoosts, Revenue: revenue,
	}, nil
}
