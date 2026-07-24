package service

import (
	"strings"
	"time"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// AccountService handles the authenticated user's own profile and devices.
type AccountService struct {
	users   *repository.UserRepository
	devices *repository.DeviceRepository
}

func NewAccountService(users *repository.UserRepository, devices *repository.DeviceRepository) *AccountService {
	return &AccountService{users: users, devices: devices}
}

// Me returns the current user as an API response.
func (s *AccountService) Me(userID uint) (*models.UserResponse, error) {
	u, err := s.users.FindByID(userID)
	if err != nil {
		return nil, response.ErrNotFound
	}
	roles, _ := s.users.RoleKeys(u.ID)
	resp := models.NewUserResponse(u, roles)
	return &resp, nil
}

// Update applies a partial profile update.
func (s *AccountService) Update(userID uint, req models.UpdateMeRequest) (*models.UserResponse, error) {
	u, err := s.users.FindByID(userID)
	if err != nil {
		return nil, response.ErrNotFound
	}
	if req.Name != nil {
		u.Name = strings.TrimSpace(*req.Name)
	}
	if req.Locale != nil {
		u.Locale = *req.Locale
	}
	if req.Avatar != nil {
		u.Avatar = *req.Avatar
	}
	if err := s.users.Update(u); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	roles, _ := s.users.RoleKeys(u.ID)
	resp := models.NewUserResponse(u, roles)
	return &resp, nil
}

// RegisterDevice stores/updates a OneSignal player id for push delivery.
func (s *AccountService) RegisterDevice(userID uint, req models.RegisterDeviceRequest) error {
	dt := &models.DeviceToken{
		UserID:          userID,
		OneSignalPlayer: req.OneSignalPlayerID,
		Platform:        req.Platform,
	}
	if err := s.devices.Upsert(dt); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}

var _ = time.Now
