package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// ShopHandler exposes the shop module (gated by the shop.enabled flag).
type ShopHandler struct {
	shop *service.ShopService
}

func NewShopHandler(shop *service.ShopService) *ShopHandler {
	return &ShopHandler{shop: shop}
}

// --- Public ---

func (h *ShopHandler) ListProducts(c echo.Context) error {
	rows, meta, err := h.shop.ListProducts(pagination.FromQuery(c), middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *ShopHandler) GetProduct(c echo.Context) error {
	detail, err := h.shop.ProductDetail(c.Param("slug"), middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, detail)
}

// --- Orders (tourist) ---

func (h *ShopHandler) CreateOrder(c echo.Context) error {
	var req models.CreateOrderRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.shop.CreateOrder(middleware.UserID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, res)
}

func (h *ShopHandler) ListOrders(c echo.Context) error {
	rows, meta, err := h.shop.ListOrders(middleware.UserID(c), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

// --- Agency products ---

func (h *ShopHandler) AgencyProducts(c echo.Context) error {
	rows, meta, err := h.shop.AgencyProducts(middleware.AgencyID(c), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *ShopHandler) CreateProduct(c echo.Context) error {
	var req models.SaveProductRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	p, err := h.shop.CreateProduct(middleware.AgencyID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, p)
}

func (h *ShopHandler) UpdateProduct(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.SaveProductRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	p, err := h.shop.UpdateProduct(middleware.AgencyID(c), id, req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, p)
}
