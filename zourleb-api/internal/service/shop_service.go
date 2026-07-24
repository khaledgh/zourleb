package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/i18n"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
	"github.com/zourleb/zourleb-api/pkg/slug"
)

// ShopService implements the toggleable shop: agency products, public catalog
// and orders. Route-level gating is enforced by the RequireFeature middleware.
type ShopService struct {
	repo    *repository.ShopRepository
	defLang string
}

func NewShopService(repo *repository.ShopRepository, defLang string) *ShopService {
	return &ShopService{repo: repo, defLang: defLang}
}

// --- Agency ---

func (s *ShopService) CreateProduct(agencyID uint, req models.SaveProductRequest) (*models.Product, error) {
	name := ""
	if len(req.Translations) > 0 {
		name = req.Translations[0].Name
	}
	status := req.Status
	if status == "" {
		status = models.ProductActive
	}
	p := &models.Product{
		AgencyID: agencyID, Slug: slug.MakeUnique(name), CategoryID: req.CategoryID,
		Price: req.Price, Currency: req.Currency, Stock: req.Stock, Status: status,
	}
	if err := s.repo.CreateProduct(p); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	for _, t := range req.Translations {
		_ = s.repo.UpsertProductTranslation(&models.ProductTranslation{
			ProductID: p.ID, Locale: t.Locale, Name: t.Name, Description: t.Description,
		})
	}
	return s.repo.FindProductForAgency(p.ID, agencyID)
}

func (s *ShopService) UpdateProduct(agencyID, id uint, req models.SaveProductRequest) (*models.Product, error) {
	p, err := s.repo.FindProductForAgency(id, agencyID)
	if err != nil {
		return nil, response.ErrNotFound.WithMessage("Product not found.")
	}
	p.CategoryID = req.CategoryID
	p.Price = req.Price
	p.Currency = req.Currency
	p.Stock = req.Stock
	if req.Status != "" {
		p.Status = req.Status
	}
	if err := s.repo.SaveProduct(p); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	for _, t := range req.Translations {
		_ = s.repo.UpsertProductTranslation(&models.ProductTranslation{
			ProductID: p.ID, Locale: t.Locale, Name: t.Name, Description: t.Description,
		})
	}
	return s.repo.FindProductForAgency(p.ID, agencyID)
}

func (s *ShopService) AgencyProducts(agencyID uint, p pagination.Params) ([]models.Product, pagination.Meta, error) {
	rows, total, err := s.repo.ListForAgency(agencyID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

// --- Public ---

func (s *ShopService) ListProducts(p pagination.Params, locale string) ([]models.ProductCard, pagination.Meta, error) {
	rows, total, err := s.repo.ListActive(p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	out := make([]models.ProductCard, 0, len(rows))
	for _, pr := range rows {
		out = append(out, s.toCard(pr, locale))
	}
	return out, pagination.NewMeta(p, total), nil
}

func (s *ShopService) ProductDetail(slug, locale string) (*models.ProductDetail, error) {
	p, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, response.ErrNotFound.WithMessage("Product not found.")
	}
	tr, _ := i18n.Resolve(p.Translations, locale, s.defLang)
	detail := &models.ProductDetail{
		ProductCard: s.toCard(*p, locale),
		Description: tr.Description,
	}
	for _, img := range p.Images {
		detail.Images = append(detail.Images, models.TourImageDTO{
			URL: img.URL, SortOrder: img.SortOrder, IsCover: img.IsCover,
		})
	}
	return detail, nil
}

// --- Orders ---

// CreateOrder builds an order from cart items, pricing from current product
// data and decrementing stock atomically.
func (s *ShopService) CreateOrder(userID uint, req models.CreateOrderRequest) (*models.OrderResponse, error) {
	ids := make([]uint, 0, len(req.Items))
	qty := map[uint]int{}
	for _, it := range req.Items {
		ids = append(ids, it.ProductID)
		qty[it.ProductID] += it.Quantity
	}
	products, err := s.repo.ProductsByIDs(ids)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	if len(products) != len(qty) {
		return nil, response.ErrBadRequest.WithMessage("One or more products are unavailable.")
	}

	var subtotal float64
	currency := ""
	order := &models.Order{
		Code: newBookingCode(), UserID: userID, Status: models.OrderPending,
	}
	for _, p := range products {
		if currency == "" {
			currency = p.Currency
		}
		line := p.Price * float64(qty[p.ID])
		subtotal += line
		order.Items = append(order.Items, models.OrderItem{
			ProductID: p.ID, Quantity: qty[p.ID], UnitPrice: p.Price, Currency: p.Currency,
		})
	}
	order.Subtotal = subtotal
	order.Currency = currency

	if err := s.repo.CreateOrder(order, qty); err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, response.ErrConflict.WithMessage("Insufficient stock for one or more items.")
		}
		return nil, response.ErrInternal.Wrap(err)
	}
	resp := models.NewOrderResponse(order)
	return &resp, nil
}

func (s *ShopService) ListOrders(userID uint, p pagination.Params) ([]models.Order, pagination.Meta, error) {
	rows, total, err := s.repo.ListOrders(userID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

func (s *ShopService) toCard(p models.Product, locale string) models.ProductCard {
	tr, _ := i18n.Resolve(p.Translations, locale, s.defLang)
	cover := ""
	for _, img := range p.Images {
		if img.IsCover {
			cover = img.URL
			break
		}
	}
	return models.ProductCard{
		ID: p.ID, Slug: p.Slug, Name: tr.Name, Price: p.Price,
		Currency: p.Currency, Stock: p.Stock, Cover: cover, AgencyID: p.AgencyID,
	}
}
