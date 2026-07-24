package models

// Product status values.
const (
	ProductActive     = "active"
	ProductHidden     = "hidden"
	ProductOutOfStock = "out_of_stock"
)

// Order status values.
const (
	OrderPending   = "pending"
	OrderPaid      = "paid"
	OrderFulfilled = "fulfilled"
	OrderCancelled = "cancelled"
)

// Product is a sellable item in the (toggleable) shop module.
type Product struct {
	SoftBase
	AgencyID   uint    `gorm:"index" json:"agency_id"`
	Slug       string  `gorm:"size:191;uniqueIndex" json:"slug"`
	CategoryID *uint   `gorm:"index" json:"category_id"`
	Price      float64 `gorm:"type:decimal(12,2)" json:"price"`
	Currency   string  `gorm:"size:3;default:USD" json:"currency"`
	Stock      int     `gorm:"default:0" json:"stock"`
	Status     string  `gorm:"size:16;default:active;index" json:"status"`

	Translations []ProductTranslation `gorm:"foreignKey:ProductID" json:"translations,omitempty"`
	Images       []ProductImage       `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

type ProductTranslation struct {
	Base
	ProductID   uint   `gorm:"index;uniqueIndex:uq_product_locale" json:"product_id"`
	Locale      string `gorm:"size:8;uniqueIndex:uq_product_locale" json:"locale"`
	Name        string `gorm:"size:200" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

type ProductImage struct {
	Base
	ProductID uint   `gorm:"index" json:"product_id"`
	URL       string `gorm:"size:512" json:"url"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsCover   bool   `gorm:"default:false" json:"is_cover"`
}

type Order struct {
	SoftBase
	Code     string  `gorm:"size:16;uniqueIndex" json:"code"`
	UserID   uint    `gorm:"index" json:"user_id"`
	Status   string  `gorm:"size:16;default:pending;index" json:"status"`
	Subtotal float64 `gorm:"type:decimal(14,2)" json:"subtotal"`
	Currency string  `gorm:"size:3" json:"currency"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	Base
	OrderID   uint    `gorm:"index" json:"order_id"`
	ProductID uint    `gorm:"index" json:"product_id"`
	Quantity  int     `gorm:"default:1" json:"quantity"`
	UnitPrice float64 `gorm:"type:decimal(12,2)" json:"unit_price"`
	Currency  string  `gorm:"size:3" json:"currency"`
}
