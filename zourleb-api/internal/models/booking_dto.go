package models

import "time"

// CreateBookingRequest is the tourist booking payload. The contact phone must
// already be verified (via OTP) when settings require it.
type CreateBookingRequest struct {
	TourID       uint               `json:"tour_id" validate:"required"`
	DepartureID  *uint              `json:"departure_id"`
	ContactPhone string             `json:"contact_phone" validate:"required,lb_phone"`
	Notes        string             `json:"notes" validate:"omitempty,max=1000"`
	Travelers    []TravelerRequest  `json:"travelers" validate:"required,min=1,dive"`
}

type TravelerRequest struct {
	FullName     string `json:"full_name" validate:"required,min=2,max=160"`
	TravelerType string `json:"traveler_type" validate:"required,oneof=adult child infant student"`
	Phone        string `json:"phone" validate:"omitempty,lb_phone"`
}

// BookingResponse is the API representation of a booking.
type BookingResponse struct {
	ID            uint      `json:"id"`
	Code          string    `json:"code"`
	TourID        uint      `json:"tour_id"`
	TourTitle     string    `json:"tour_title"`
	DepartureID   *uint     `json:"departure_id"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	TravelersCount int      `json:"travelers_count"`
	Subtotal      float64   `json:"subtotal"`
	Currency      string    `json:"currency"`
	ContactPhone  string    `json:"contact_phone"`
	CreatedAt     time.Time `json:"created_at"`
	Travelers     []TravelerResponse `json:"travelers,omitempty"`
}

type TravelerResponse struct {
	FullName     string `json:"full_name"`
	TravelerType string `json:"traveler_type"`
	Phone        string `json:"phone"`
}

// NewBookingResponse maps a Booking entity to its API shape.
func NewBookingResponse(b *Booking, tourTitle string) BookingResponse {
	resp := BookingResponse{
		ID: b.ID, Code: b.Code, TourID: b.TourID, TourTitle: tourTitle,
		DepartureID: b.DepartureID, Status: b.Status, PaymentStatus: b.PaymentStatus,
		TravelersCount: b.TravelersCount, Subtotal: b.Subtotal, Currency: b.Currency,
		ContactPhone: b.ContactPhone, CreatedAt: b.CreatedAt,
	}
	for _, t := range b.Travelers {
		resp.Travelers = append(resp.Travelers, TravelerResponse{
			FullName: t.FullName, TravelerType: t.TravelerType, Phone: t.Phone,
		})
	}
	return resp
}
