package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// UploadHandler accepts multipart image uploads and returns stored URLs.
type UploadHandler struct {
	upload *service.UploadService
}

func NewUploadHandler(upload *service.UploadService) *UploadHandler {
	return &UploadHandler{upload: upload}
}

// Image handles POST /uploads/image (multipart form field "file").
func (h *UploadHandler) Image(c echo.Context) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest.WithMessage("Missing 'file' field."))
	}
	src, err := fh.Open()
	if err != nil {
		return response.Fail(c, response.ErrInternal.Wrap(err))
	}
	defer src.Close()

	prefix := c.QueryParam("prefix")
	if prefix == "" {
		prefix = "uploads"
	}
	url, err := h.upload.SaveImage(
		c.Request().Context(), prefix, fh.Filename,
		fh.Header.Get("Content-Type"), fh.Size, src,
	)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, map[string]string{"url": url})
}
