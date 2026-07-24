package service

import (
	"context"
	"io"

	"github.com/zourleb/zourleb-api/pkg/response"
	"github.com/zourleb/zourleb-api/pkg/storage"
)

// UploadService stores uploaded media and returns its public URL.
type UploadService struct {
	store storage.Storage
}

func NewUploadService(store storage.Storage) *UploadService {
	return &UploadService{store: store}
}

// allowed image content types
var allowedImage = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

const maxUploadBytes = 8 << 20 // 8 MiB

// SaveImage validates and stores an uploaded image, returning its URL.
func (s *UploadService) SaveImage(ctx context.Context, prefix, filename, contentType string, size int64, r io.Reader) (string, error) {
	if !allowedImage[contentType] {
		return "", response.ErrBadRequest.WithMessage("Only JPEG, PNG, WEBP or GIF images are allowed.")
	}
	if size > maxUploadBytes {
		return "", response.ErrBadRequest.WithMessage("Image exceeds the 8 MB limit.")
	}
	key := storage.Key(prefix, filename)
	url, err := s.store.Save(ctx, key, r, contentType)
	if err != nil {
		return "", response.ErrInternal.Wrap(err)
	}
	return url, nil
}
