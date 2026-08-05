package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// AdminHandler exposes the super-admin endpoints.
type AdminHandler struct {
	admin *service.AdminService
}

func NewAdminHandler(admin *service.AdminService) *AdminHandler {
	return &AdminHandler{admin: admin}
}

// --- Agencies ---

func (h *AdminHandler) ListAgencies(c echo.Context) error {
	rows, meta, err := h.admin.ListAgencies(c.QueryParam("status"), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *AdminHandler) ApproveAgency(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.ApproveAgencyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.ApproveAgency(id, req.Approve); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"approved": req.Approve})
}

// --- Languages ---

func (h *AdminHandler) Languages(c echo.Context) error {
	rows, err := h.admin.Languages()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

func (h *AdminHandler) CreateLanguage(c echo.Context) error {
	var req models.SaveLanguageRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	l, err := h.admin.CreateLanguage(req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, l)
}

func (h *AdminHandler) UpdateLanguage(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.SaveLanguageRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	l, err := h.admin.UpdateLanguage(id, req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, l)
}

// --- Translations ---

func (h *AdminHandler) SaveTranslation(c echo.Context) error {
	var req models.SaveTranslationRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.SaveTranslation(req); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"saved": true})
}

func (h *AdminHandler) BulkTranslations(c echo.Context) error {
	var req models.BulkTranslationsRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.BulkTranslations(req); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]int{"imported": len(req.Items)})
}

// --- Settings ---

func (h *AdminHandler) Settings(c echo.Context) error {
	rows, err := h.admin.Settings()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

func (h *AdminHandler) SaveSetting(c echo.Context) error {
	var req models.SaveSettingRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.SaveSetting(req); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"saved": true})
}

// --- Banners ---

func (h *AdminHandler) Banners(c echo.Context) error {
	rows, err := h.admin.Banners()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

func (h *AdminHandler) CreateBanner(c echo.Context) error {
	var req models.SaveBannerRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	b, err := h.admin.CreateBanner(req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, b)
}

func (h *AdminHandler) DeleteBanner(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.admin.DeleteBanner(id); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"deleted": true})
}

// --- Categories & Regions ---

func (h *AdminHandler) CreateCategory(c echo.Context) error {
	var req models.SaveCategoryRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	cat, err := h.admin.CreateCategory(req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, cat)
}

func (h *AdminHandler) CreateRegion(c echo.Context) error {
	var req models.SaveRegionRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	rg, err := h.admin.CreateRegion(req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, rg)
}

// --- Reviews ---

func (h *AdminHandler) PendingReviews(c echo.Context) error {
	rows, meta, err := h.admin.PendingReviews(pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *AdminHandler) ModerateReview(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.ModerateReviewRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.ModerateReview(id, req.Status); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]string{"status": req.Status})
}

// --- Users ---

func (h *AdminHandler) ListUsers(c echo.Context) error {
	rows, meta, err := h.admin.ListUsers(pagination.FromQuery(c), c.QueryParam("q"))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *AdminHandler) SetUserStatus(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.SetUserStatusRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.SetUserStatus(id, req.Status); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]string{"status": req.Status})
}

func (h *AdminHandler) CreateUser(c echo.Context) error {
	var req models.AdminCreateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	u, err := h.admin.CreateUser(req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, u)
}

func (h *AdminHandler) UpdateUser(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.AdminUpdateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	u, err := h.admin.UpdateUser(id, req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, u)
}

func (h *AdminHandler) DeleteUser(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.admin.DeleteUser(id); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"deleted": true})
}

func (h *AdminHandler) CreateAgency(c echo.Context) error {
	var req models.AdminCreateAgencyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	a, err := h.admin.CreateAgency(req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, a)
}

func (h *AdminHandler) UpdateAgency(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.AdminUpdateAgencyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	a, err := h.admin.UpdateAgency(id, req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, a)
}

func (h *AdminHandler) DeleteAgency(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.admin.DeleteAgency(id); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"deleted": true})
}

// --- Boosts & payments ---

func (h *AdminHandler) PendingBoosts(c echo.Context) error {
	rows, err := h.admin.PendingBoosts()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

func (h *AdminHandler) ListTours(c echo.Context) error {
	rows, meta, err := h.admin.ListTours(pagination.FromQuery(c), c.QueryParam("q"), c.QueryParam("status"))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *AdminHandler) ConfirmPayment(c echo.Context) error {
	var req models.ConfirmPaymentRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.admin.ConfirmPayment(c.Request().Context(), req.PaymentID); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"confirmed": true})
}

// --- Analytics ---

func (h *AdminHandler) Analytics(c echo.Context) error {
	summary, err := h.admin.Analytics()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, summary)
}
