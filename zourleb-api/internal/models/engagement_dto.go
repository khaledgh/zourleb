package models

// CreateReviewRequest is the verified-booking review payload.
type CreateReviewRequest struct {
	TourID  uint   `json:"tour_id" validate:"required"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"omitempty,max=2000"`
}

// ReviewResponse is the API representation of a review.
type ReviewResponse struct {
	ID      uint   `json:"id"`
	TourID  uint   `json:"tour_id"`
	UserID  uint   `json:"user_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Status  string `json:"status"`
}

func NewReviewResponse(r *Review) ReviewResponse {
	return ReviewResponse{
		ID: r.ID, TourID: r.TourID, UserID: r.UserID,
		Rating: r.Rating, Comment: r.Comment, Status: r.Status,
	}
}
