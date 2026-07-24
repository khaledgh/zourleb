package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zourleb/zourleb-api/internal/models"
)

// DeviceRepository persists OneSignal device tokens.
type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Upsert inserts or updates a device token by its OneSignal player id.
func (r *DeviceRepository) Upsert(dt *models.DeviceToken) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "one_signal_player"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "platform", "updated_at"}),
	}).Create(dt).Error
}

// ByUser returns all device tokens for a user.
func (r *DeviceRepository) ByUser(userID uint) ([]models.DeviceToken, error) {
	var rows []models.DeviceToken
	err := r.db.Where("user_id = ?", userID).Find(&rows).Error
	return rows, err
}
