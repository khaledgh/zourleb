package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// AgencyHandler exposes the agency portal endpoints (RBAC + agency-scoped).
type AgencyHandler struct {
	agency *service.AgencyService
	upload *service.UploadService
}

func NewAgencyHandler(agency *service.AgencyService, upload *service.UploadService) *AgencyHandler {
	return &AgencyHandler{agency: agency, upload: upload}
}

// Apply handles POST /agency/apply — any authenticated user can apply.
func (h *AgencyHandler) Apply(c echo.Context) error {
	var req models.ApplyAgencyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	a, err := h.agency.Apply(middleware.UserID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, a)
}

// Profile handles GET /agency/profile.
func (h *AgencyHandler) Profile(c echo.Context) error {
	a, err := h.agency.Profile(middleware.AgencyID(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, a)
}

// UpdateProfile handles PATCH /agency/profile.
func (h *AgencyHandler) UpdateProfile(c echo.Context) error {
	var req models.UpdateAgencyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	a, err := h.agency.UpdateProfile(middleware.AgencyID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, a)
}

// Members handles GET /agency/members.
func (h *AgencyHandler) Members(c echo.Context) error {
	rows, err := h.agency.Members(middleware.AgencyID(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

// --- Tours ---

func (h *AgencyHandler) ListTours(c echo.Context) error {
	rows, meta, err := h.agency.ListTours(middleware.AgencyID(c), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

func (h *AgencyHandler) CreateTour(c echo.Context) error {
	var req models.SaveTourRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	t, err := h.agency.CreateTour(middleware.AgencyID(c), middleware.UserID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, t)
}

func (h *AgencyHandler) GetTour(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	t, err := h.agency.GetTour(middleware.AgencyID(c), id)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, t)
}

func (h *AgencyHandler) UpdateTour(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.SaveTourRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	t, err := h.agency.UpdateTour(middleware.AgencyID(c), id, req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, t)
}

func (h *AgencyHandler) Publish(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.PublishTourRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.agency.Publish(middleware.AgencyID(c), id, req.Publish); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"published": req.Publish})
}

func (h *AgencyHandler) AddDeparture(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.AddDepartureRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.agency.AddDeparture(middleware.AgencyID(c), id, req); err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, map[string]bool{"added": true})
}

func (h *AgencyHandler) AddPrice(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.AddPriceRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.agency.AddPrice(middleware.AgencyID(c), id, req); err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, map[string]bool{"added": true})
}

func (h *AgencyHandler) AddImage(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.AddImageRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.agency.AddImage(middleware.AgencyID(c), id, req); err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, map[string]bool{"added": true})
}

func (h *AgencyHandler) DeleteImage(c echo.Context) error {
	imageID, err := paramUint(c, "imageId")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.agency.DeleteImage(middleware.AgencyID(c), imageID); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"deleted": true})
}

func (h *AgencyHandler) Bookings(c echo.Context) error {
	rows, meta, err := h.agency.Bookings(middleware.AgencyID(c), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}
