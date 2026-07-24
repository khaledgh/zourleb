package service

import (
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// EngagementService handles favorites and reviews.
type EngagementService struct {
	repo     *repository.EngagementRepository
	settings *SettingsService
}

func NewEngagementService(repo *repository.EngagementRepository, settings *SettingsService) *EngagementService {
	return &EngagementService{repo: repo, settings: settings}
}

func (s *EngagementService) AddFavorite(userID, tourID uint) error {
	if err := s.repo.AddFavorite(userID, tourID); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}

func (s *EngagementService) RemoveFavorite(userID, tourID uint) error {
	if err := s.repo.RemoveFavorite(userID, tourID); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}

func (s *EngagementService) Favorites(userID uint) ([]uint, error) {
	ids, err := s.repo.FavoriteTourIDs(userID)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return ids, nil
}

// CreateReview enforces the reviews flag and verified-booking rule, storing the
// review as pending for moderation.
func (s *EngagementService) CreateReview(userID uint, req models.CreateReviewRequest) (*models.ReviewResponse, error) {
	if !s.settings.Bool(models.SettingReviewsEnabled, true) {
		return nil, response.ErrFeatureOff
	}
	ok, err := s.repo.HasCompletedBooking(userID, req.TourID)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	if !ok {
		return nil, response.ErrForbidden.WithMessage("Only travelers who booked this tour can review it.")
	}
	rv := &models.Review{
		UserID: userID, TourID: req.TourID,
		Rating: req.Rating, Comment: req.Comment, Status: models.ReviewPending,
	}
	if err := s.repo.CreateReview(rv); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	resp := models.NewReviewResponse(rv)
	return &resp, nil
}

// TourReviews lists approved reviews for a tour.
func (s *EngagementService) TourReviews(tourID uint, p pagination.Params) ([]models.ReviewResponse, pagination.Meta, error) {
	rows, total, err := s.repo.ApprovedReviewsForTour(tourID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	out := make([]models.ReviewResponse, 0, len(rows))
	for i := range rows {
		out = append(out, models.NewReviewResponse(&rows[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}
