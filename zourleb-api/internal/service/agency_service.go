package service

import (
	"strings"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/phone"
	"github.com/zourleb/zourleb-api/pkg/response"
	"github.com/zourleb/zourleb-api/pkg/slug"
)

// AgencyService implements the agency portal: profile, tours, departures,
// prices, images, members, and the public application flow.
type AgencyService struct {
	agencies *repository.AgencyRepository
	tours    *repository.TourRepository
	bookings *repository.BookingRepository
	users    *repository.UserRepository
	defLang  string
}

func NewAgencyService(a *repository.AgencyRepository, t *repository.TourRepository, b *repository.BookingRepository, u *repository.UserRepository, defLang string) *AgencyService {
	return &AgencyService{agencies: a, tours: t, bookings: b, users: u, defLang: defLang}
}

// Apply registers a new agency (pending approval) and makes the caller its
// owner (agency-scoped role).
func (s *AgencyService) Apply(userID uint, req models.ApplyAgencyRequest) (*models.Agency, error) {
	ph, _ := phone.Normalize(req.Phone)
	a := &models.Agency{
		Slug:           slug.MakeUnique(req.Name),
		Name:           strings.TrimSpace(req.Name),
		Phone:          ph,
		Email:          strings.ToLower(req.Email),
		RegionID:       req.RegionID,
		Status:         models.AgencyStatusPending,
		CommissionRate: 10,
		CreatedBy:      userID,
	}
	if err := s.agencies.Create(a); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}

	ownerRoleID, err := s.users.RoleIDByKey(models.RoleAgencyOwner)
	if err == nil {
		_ = s.users.AssignRole(userID, ownerRoleID, &a.ID)
		_ = s.agencies.AddMember(&models.AgencyMember{AgencyID: a.ID, UserID: userID, RoleID: ownerRoleID})
	}
	return a, nil
}

// Profile returns the caller's agency.
func (s *AgencyService) Profile(agencyID uint) (*models.Agency, error) {
	a, err := s.agencies.FindByID(agencyID)
	if err != nil {
		return nil, response.ErrAgencyNotFound
	}
	return a, nil
}

// UpdateProfile applies a partial profile update with optional translations.
func (s *AgencyService) UpdateProfile(agencyID uint, req models.UpdateAgencyRequest) (*models.Agency, error) {
	a, err := s.agencies.FindByID(agencyID)
	if err != nil {
		return nil, response.ErrAgencyNotFound
	}
	if req.Name != nil {
		a.Name = *req.Name
	}
	if req.Phone != nil {
		if p, err := phone.Normalize(*req.Phone); err == nil {
			a.Phone = p
		}
	}
	if req.Email != nil {
		a.Email = strings.ToLower(*req.Email)
	}
	if req.Website != nil {
		a.Website = *req.Website
	}
	if req.RegionID != nil {
		a.RegionID = req.RegionID
	}
	if req.Logo != nil {
		a.Logo = *req.Logo
	}
	if req.Cover != nil {
		a.Cover = *req.Cover
	}
	if err := s.agencies.Save(a); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	for _, t := range req.Translations {
		_ = s.agencies.UpsertTranslation(&models.AgencyTranslation{
			AgencyID: agencyID, Locale: t.Locale, Description: t.Description, About: t.About,
		})
	}
	return s.agencies.FindByID(agencyID)
}

// Members lists the agency's members.
func (s *AgencyService) Members(agencyID uint) ([]models.AgencyMember, error) {
	rows, err := s.agencies.Members(agencyID)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

// --- Tours ---

// CreateTour drafts a new tour with its translations.
func (s *AgencyService) CreateTour(agencyID, userID uint, req models.SaveTourRequest) (*models.Tour, error) {
	title := ""
	if len(req.Translations) > 0 {
		title = req.Translations[0].Title
	}
	cur := req.BaseCurrency
	if cur == "" {
		cur = "USD"
	}
	t := &models.Tour{
		AgencyID:     agencyID,
		Slug:         slug.MakeUnique(title),
		CategoryID:   req.CategoryID,
		RegionID:     req.RegionID,
		Type:         req.Type,
		DurationDays: req.DurationDays,
		Difficulty:   req.Difficulty,
		MinAge:       req.MinAge,
		MaxCapacity:  req.MaxCapacity,
		BaseCurrency: cur,
		Status:       models.TourStatusDraft,
		IsShowcase:   req.IsShowcase,
		StartsFromDate: req.StartsFromDate,
		CreatedBy:    userID,
	}
	if err := s.tours.Create(t); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	for _, tr := range req.Translations {
		_ = s.tours.UpsertTranslation(&models.TourTranslation{
			TourID: t.ID, Locale: tr.Locale, Title: tr.Title, Summary: tr.Summary,
			Description: tr.Description, Itinerary: tr.Itinerary,
			Included: tr.Included, Excluded: tr.Excluded,
		})
	}
	return s.tours.FindByIDForAgency(t.ID, agencyID)
}

// UpdateTour edits an existing tour the agency owns.
func (s *AgencyService) UpdateTour(agencyID, tourID uint, req models.SaveTourRequest) (*models.Tour, error) {
	t, err := s.tours.FindByIDForAgency(tourID, agencyID)
	if err != nil {
		return nil, response.ErrTourNotFound
	}
	t.CategoryID = req.CategoryID
	t.RegionID = req.RegionID
	t.Type = req.Type
	t.DurationDays = req.DurationDays
	t.Difficulty = req.Difficulty
	t.MinAge = req.MinAge
	t.MaxCapacity = req.MaxCapacity
	if req.BaseCurrency != "" {
		t.BaseCurrency = req.BaseCurrency
	}
	t.IsShowcase = req.IsShowcase
	t.StartsFromDate = req.StartsFromDate
	if err := s.tours.Save(t); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	for _, tr := range req.Translations {
		_ = s.tours.UpsertTranslation(&models.TourTranslation{
			TourID: t.ID, Locale: tr.Locale, Title: tr.Title, Summary: tr.Summary,
			Description: tr.Description, Itinerary: tr.Itinerary,
			Included: tr.Included, Excluded: tr.Excluded,
		})
	}
	return s.tours.FindByIDForAgency(t.ID, agencyID)
}

// ListTours returns the agency's own tours.
func (s *AgencyService) ListTours(agencyID uint, p pagination.Params) ([]models.Tour, pagination.Meta, error) {
	rows, total, err := s.tours.ListForAgency(agencyID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

// GetTour returns one of the agency's tours.
func (s *AgencyService) GetTour(agencyID, tourID uint) (*models.Tour, error) {
	t, err := s.tours.FindByIDForAgency(tourID, agencyID)
	if err != nil {
		return nil, response.ErrTourNotFound
	}
	return t, nil
}

// Publish toggles a tour between draft and published.
func (s *AgencyService) Publish(agencyID, tourID uint, publish bool) error {
	t, err := s.tours.FindByIDForAgency(tourID, agencyID)
	if err != nil {
		return response.ErrTourNotFound
	}
	if publish {
		t.Status = models.TourStatusPublished
	} else {
		t.Status = models.TourStatusDraft
	}
	if err := s.tours.Save(t); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}

// AddDeparture adds a departure to a tour.
func (s *AgencyService) AddDeparture(agencyID, tourID uint, req models.AddDepartureRequest) error {
	if _, err := s.tours.FindByIDForAgency(tourID, agencyID); err != nil {
		return response.ErrTourNotFound
	}
	return wrapInternal(s.tours.AddDeparture(&models.TourDeparture{
		TourID: tourID, StartDate: req.StartDate, EndDate: req.EndDate,
		Capacity: req.Capacity, Status: models.DepartureOpen,
	}))
}

// AddPrice adds a price row and refreshes the tour's price_from.
func (s *AgencyService) AddPrice(agencyID, tourID uint, req models.AddPriceRequest) error {
	if _, err := s.tours.FindByIDForAgency(tourID, agencyID); err != nil {
		return response.ErrTourNotFound
	}
	if err := s.tours.AddPrice(&models.TourPrice{
		TourID: tourID, DepartureID: req.DepartureID, TravelerType: req.TravelerType,
		Currency: req.Currency, Amount: req.Amount,
	}); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	_ = s.tours.UpdatePriceFrom(tourID)
	return nil
}

// AddImage attaches an image (current or past-gallery) to a tour.
func (s *AgencyService) AddImage(agencyID, tourID uint, req models.AddImageRequest) error {
	if _, err := s.tours.FindByIDForAgency(tourID, agencyID); err != nil {
		return response.ErrTourNotFound
	}
	src := req.Source
	if src == "" {
		src = models.ImageSourceCurrent
	}
	return wrapInternal(s.tours.AddImage(&models.TourImage{
		TourID: tourID, URL: req.URL, IsCover: req.IsCover, Source: src,
	}))
}

// DeleteImage removes a tour image owned by the agency.
func (s *AgencyService) DeleteImage(agencyID, imageID uint) error {
	return wrapInternal(s.tours.DeleteImage(imageID, agencyID))
}

// Bookings lists bookings across the agency's tours.
func (s *AgencyService) Bookings(agencyID uint, p pagination.Params) ([]models.Booking, pagination.Meta, error) {
	rows, total, err := s.bookings.ListForAgency(agencyID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

func wrapInternal(err error) error {
	if err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}
