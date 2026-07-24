package service

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/onesignal"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// NotificationService persists in-app notifications and dispatches push.
type NotificationService struct {
	repo *repository.NotificationRepository
	push *onesignal.Client
}

func NewNotificationService(repo *repository.NotificationRepository, push *onesignal.Client) *NotificationService {
	return &NotificationService{repo: repo, push: push}
}

// Notify persists a notification row and best-effort sends a push. Push errors
// never fail the caller (notifications are non-critical, sent async-ish).
func (s *NotificationService) Notify(ctx context.Context, userID uint, typ, title, body string, data map[string]any) {
	var payload datatypes.JSON
	if data != nil {
		if b, err := json.Marshal(data); err == nil {
			payload = b
		}
	}
	_ = s.repo.Create(&models.Notification{
		UserID: userID, Type: typ, Title: title, Body: body, Data: payload,
	})

	if ids, err := s.repo.PlayerIDsForUser(userID); err == nil && len(ids) > 0 {
		_ = s.push.Send(ctx, onesignal.Notification{
			PlayerIDs: ids, Title: title, Body: body, Data: data,
		})
	}
}

// List returns the user's notifications.
func (s *NotificationService) List(userID uint, p pagination.Params) ([]models.Notification, pagination.Meta, error) {
	rows, total, err := s.repo.ListForUser(userID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	return rows, pagination.NewMeta(p, total), nil
}

// MarkRead flags a notification as read.
func (s *NotificationService) MarkRead(userID, id uint) error {
	if err := s.repo.MarkRead(userID, id); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}
