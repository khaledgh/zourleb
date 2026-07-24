package models

// --- Product authoring (agency) ---

type SaveProductRequest struct {
	CategoryID   *uint                     `json:"category_id"`
	Price        float64                   `json:"price" validate:"required,min=0"`
	Currency     string                    `json:"currency" validate:"required,oneof=USD LBP"`
	Stock        int                       `json:"stock" validate:"min=0"`
	Status       string                    `json:"status" validate:"omitempty,oneof=active hidden out_of_stock"`
	Translations []ProductTranslationInput `json:"translations" validate:"required,min=1,dive"`
}

type ProductTranslationInput struct {
	Locale      string `json:"locale" validate:"required,max=8"`
	Name        string `json:"name" validate:"required,max=200"`
	Description string `json:"description" validate:"omitempty"`
}

// --- Public product representation (localized) ---

type ProductCard struct {
	ID       uint    `json:"id"`
	Slug     string  `json:"slug"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Stock    int     `json:"stock"`
	Cover    string  `json:"cover"`
	AgencyID uint    `json:"agency_id"`
}

type ProductDetail struct {
	ProductCard
	Description string         `json:"description"`
	Images     []TourImageDTO `json:"images"`
}

// --- Orders ---

type CreateOrderRequest struct {
	Items []OrderItemInput `json:"items" validate:"required,min=1,dive"`
}

type OrderItemInput struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type OrderResponse struct {
	ID       uint    `json:"id"`
	Code     string  `json:"code"`
	Status   string  `json:"status"`
	Subtotal float64 `json:"subtotal"`
	Currency string  `json:"currency"`
}

func NewOrderResponse(o *Order) OrderResponse {
	return OrderResponse{
		ID: o.ID, Code: o.Code, Status: o.Status,
		Subtotal: o.Subtotal, Currency: o.Currency,
	}
}
