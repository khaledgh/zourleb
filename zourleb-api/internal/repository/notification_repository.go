package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// NotificationRepository persists in-app notifications.
type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepository) ListForUser(userID uint, p pagination.Params) ([]models.Notification, int64, error) {
	q := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Notification
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

func (r *NotificationRepository) MarkRead(userID, id uint) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read_at", now).Error
}

// PlayerIDsForUser returns OneSignal player ids for push delivery.
func (r *NotificationRepository) PlayerIDsForUser(userID uint) ([]string, error) {
	var ids []string
	err := r.db.Model(&models.DeviceToken{}).Where("user_id = ?", userID).
		Pluck("one_signal_player", &ids).Error
	return ids, err
}
