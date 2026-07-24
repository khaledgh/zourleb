package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// EngagementRepository handles favorites and reviews.
type EngagementRepository struct {
	db *gorm.DB
}

func NewEngagementRepository(db *gorm.DB) *EngagementRepository {
	return &EngagementRepository{db: db}
}

// --- Favorites ---

func (r *EngagementRepository) AddFavorite(userID, tourID uint) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&models.Favorite{UserID: userID, TourID: tourID}).Error
}

func (r *EngagementRepository) RemoveFavorite(userID, tourID uint) error {
	return r.db.Where("user_id = ? AND tour_id = ?", userID, tourID).
		Delete(&models.Favorite{}).Error
}

func (r *EngagementRepository) FavoriteTourIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.Favorite{}).Where("user_id = ?", userID).
		Pluck("tour_id", &ids).Error
	return ids, err
}

// --- Reviews ---

func (r *EngagementRepository) CreateReview(rv *models.Review) error {
	return r.db.Create(rv).Error
}

// HasCompletedBooking verifies the user actually traveled (verified review).
func (r *EngagementRepository) HasCompletedBooking(userID, tourID uint) (bool, error) {
	var n int64
	err := r.db.Model(&models.Booking{}).
		Where("user_id = ? AND tour_id = ? AND status IN ?", userID, tourID,
			[]string{models.BookingConfirmed, models.BookingCompleted}).
		Count(&n).Error
	return n > 0, err
}

// ApprovedReviewsForTour lists approved reviews, paginated.
func (r *EngagementRepository) ApprovedReviewsForTour(tourID uint, p pagination.Params) ([]models.Review, int64, error) {
	q := r.db.Model(&models.Review{}).
		Where("tour_id = ? AND status = ?", tourID, models.ReviewApproved)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Review
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

// RecalcAgencyRating recomputes an agency's rating from approved reviews across
// its tours. Called after moderation.
func (r *EngagementRepository) RecalcAgencyRating(agencyID uint) error {
	var stat struct {
		Avg   float64
		Count int
	}
	r.db.Model(&models.Review{}).
		Select("AVG(rating) as avg, COUNT(*) as count").
		Where("status = ? AND tour_id IN (?)", models.ReviewApproved,
			r.db.Model(&models.Tour{}).Select("id").Where("agency_id = ?", agencyID)).
		Scan(&stat)
	return r.db.Model(&models.Agency{}).Where("id = ?", agencyID).
		Updates(map[string]any{"rating_avg": stat.Avg, "rating_count": stat.Count}).Error
}

func (r *EngagementRepository) PendingReviews(p pagination.Params) ([]models.Review, int64, error) {
	q := r.db.Model(&models.Review{}).Where("status = ?", models.ReviewPending)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Review
	err := q.Order("created_at asc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *EngagementRepository) ModerateReview(id uint, status string) (*models.Review, error) {
	var rv models.Review
	if err := r.db.First(&rv, id).Error; err != nil {
		return nil, err
	}
	rv.Status = status
	if err := r.db.Save(&rv).Error; err != nil {
		return nil, err
	}
	return &rv, nil
}
