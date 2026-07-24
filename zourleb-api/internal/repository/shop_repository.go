package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// ShopRepository persists products and orders.
type ShopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

// --- Products (agency) ---

func (r *ShopRepository) CreateProduct(p *models.Product) error { return r.db.Create(p).Error }
func (r *ShopRepository) SaveProduct(p *models.Product) error   { return r.db.Save(p).Error }

func (r *ShopRepository) UpsertProductTranslation(t *models.ProductTranslation) error {
	return r.db.Where("product_id = ? AND locale = ?", t.ProductID, t.Locale).
		Assign(*t).FirstOrCreate(t).Error
}

func (r *ShopRepository) FindProductForAgency(id, agencyID uint) (*models.Product, error) {
	var p models.Product
	err := r.db.Preload("Translations").Preload("Images").
		Where("id = ? AND agency_id = ?", id, agencyID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ShopRepository) ListForAgency(agencyID uint, p pagination.Params) ([]models.Product, int64, error) {
	q := r.db.Model(&models.Product{}).Where("agency_id = ?", agencyID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Product
	err := q.Preload("Translations").Preload("Images", "is_cover = ?", true).
		Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

// --- Products (public) ---

func (r *ShopRepository) ListActive(p pagination.Params) ([]models.Product, int64, error) {
	q := r.db.Model(&models.Product{}).Where("status = ?", models.ProductActive)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Product
	err := q.Preload("Translations").Preload("Images", "is_cover = ?", true).
		Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *ShopRepository) FindBySlug(slug string) (*models.Product, error) {
	var p models.Product
	err := r.db.Preload("Translations").Preload("Images").
		Where("slug = ? AND status = ?", slug, models.ProductActive).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ShopRepository) ProductsByIDs(ids []uint) ([]models.Product, error) {
	var rows []models.Product
	err := r.db.Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

// --- Orders ---

// CreateOrder persists an order and decrements stock atomically.
func (r *ShopRepository) CreateOrder(o *models.Order, stockDeltas map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for productID, qty := range stockDeltas {
			res := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", productID, qty).
				UpdateColumn("stock", gorm.Expr("stock - ?", qty))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return gorm.ErrInvalidData // insufficient stock
			}
		}
		return tx.Create(o).Error
	})
}

func (r *ShopRepository) ListOrders(userID uint, p pagination.Params) ([]models.Order, int64, error) {
	q := r.db.Model(&models.Order{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Order
	err := q.Preload("Items").Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}
