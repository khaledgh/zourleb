package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
)

// UserRepository handles persistence for users, roles, and refresh tokens.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ErrNotFound is returned when a lookup misses; services map it to a domain error.
var ErrNotFound = gorm.ErrRecordNotFound

func (r *UserRepository) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByGoogleID(gid string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("google_id = ?", gid).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) TouchLogin(id uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Update("last_login_at", now).Error
}

// AssignRole attaches a (possibly agency-scoped) role to a user, ignoring dups.
func (r *UserRepository) AssignRole(userID, roleID uint, agencyID *uint) error {
	ur := models.UserRole{UserID: userID, RoleID: roleID, AgencyID: agencyID}
	err := r.db.Where("user_id = ? AND role_id = ? AND agency_id <=> ?", userID, roleID, agencyID).
		FirstOrCreate(&ur).Error
	return err
}

// RoleKeys returns the distinct role keys held by a user (any scope).
func (r *UserRepository) RoleKeys(userID uint) ([]string, error) {
	var keys []string
	err := r.db.Model(&models.UserRole{}).
		Select("DISTINCT roles.key").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("roles.key", &keys).Error
	return keys, err
}

// Permissions returns the distinct permission keys granted to a user through
// all their roles, along with the agency scopes those roles apply to.
func (r *UserRepository) Permissions(userID uint) (map[string]struct{}, error) {
	var rows []models.Permission
	err := r.db.Distinct("permissions.*").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ?", userID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{}, len(rows))
	for _, p := range rows {
		set[p.Key] = struct{}{}
	}
	return set, nil
}

// AgencyIDsForUser returns agency IDs the user has any role in.
func (r *UserRepository) AgencyIDsForUser(userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND agency_id IS NOT NULL", userID).
		Pluck("agency_id", &ids).Error
	return ids, err
}

func (r *UserRepository) RoleIDByKey(key string) (uint, error) {
	var role models.Role
	if err := r.db.Where("`key` = ?", key).First(&role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

// --- Refresh tokens ---

func (r *UserRepository) CreateRefreshToken(t *models.RefreshToken) error {
	return r.db.Create(t).Error
}

func (r *UserRepository) FindRefreshByHash(hash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	if err := r.db.Where("token_hash = ?", hash).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *UserRepository) RevokeRefreshToken(id uint) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).Where("id = ?", id).
		Update("revoked_at", now).Error
}

func (r *UserRepository) RevokeAllRefreshTokens(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// IsNotFound reports whether err is a missing-record error.
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
