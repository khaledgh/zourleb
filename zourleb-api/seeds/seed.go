// Package seeds populates baseline data: languages, RBAC, settings, super admin.
package seeds

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/hash"
)

// Run is idempotent: it upserts baseline rows so it can run on every boot.
func Run(db *gorm.DB, cfg *config.Config, log *slog.Logger) error {
	if err := seedLanguages(db); err != nil {
		return err
	}
	if err := seedSettings(db, cfg); err != nil {
		return err
	}
	permIDs, err := seedPermissions(db)
	if err != nil {
		return err
	}
	roleIDs, err := seedRoles(db, permIDs)
	if err != nil {
		return err
	}
	if err := seedSuperAdmin(db, cfg, roleIDs[models.RoleSuperAdmin], log); err != nil {
		return err
	}
	if err := seedTranslations(db); err != nil {
		return err
	}
	log.Info("seeds applied")
	return nil
}

func seedLanguages(db *gorm.DB) error {
	langs := []models.Language{
		{Code: "ar", Name: "Arabic", NativeName: "العربية", IsRTL: true, IsActive: true, IsDefault: true, SortOrder: 1},
		{Code: "fr", Name: "French", NativeName: "Français", IsRTL: false, IsActive: true, SortOrder: 2},
		{Code: "en", Name: "English", NativeName: "English", IsRTL: false, IsActive: true, SortOrder: 3},
	}
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "code"}}, DoNothing: true}).
		Create(&langs).Error
}

func seedSettings(db *gorm.DB, cfg *config.Config) error {
	settings := []models.Setting{
		{Key: models.SettingShopEnabled, Value: "false", Type: models.SettingBool},
		{Key: models.SettingReviewsEnabled, Value: "true", Type: models.SettingBool},
		{Key: models.SettingBoostEnabled, Value: "true", Type: models.SettingBool},
		{Key: models.SettingBookingOTP, Value: "true", Type: models.SettingBool},
		{Key: models.SettingDefaultCurrency, Value: "USD", Type: models.SettingString},
	}
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "key"}}, DoNothing: true}).
		Create(&settings).Error
}

// AllPermissions is the canonical permission catalog.
func AllPermissions() []string {
	return []string{
		// agency self-service
		"agency.update", "agency.member.invite", "agency.member.remove",
		"tour.create", "tour.update", "tour.publish", "tour.delete",
		"tour.image.manage", "departure.manage", "price.manage",
		"booking.view", "analytics.view",
		"boost.purchase", "boost.view",
		"product.manage",
		// platform admin
		"admin.agency.approve", "admin.agency.suspend",
		"admin.user.manage", "admin.category.manage", "admin.region.manage",
		"admin.language.manage", "admin.translation.manage",
		"admin.boost.manage", "admin.banner.manage",
		"admin.settings.manage", "admin.review.moderate",
		"admin.analytics.view",
	}
}

func seedPermissions(db *gorm.DB) (map[string]uint, error) {
	keys := AllPermissions()
	perms := make([]models.Permission, len(keys))
	for i, k := range keys {
		perms[i] = models.Permission{Key: k}
	}
	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "key"}}, DoNothing: true}).
		Create(&perms).Error; err != nil {
		return nil, err
	}

	var all []models.Permission
	if err := db.Find(&all).Error; err != nil {
		return nil, err
	}
	byKey := make(map[string]uint, len(all))
	for _, p := range all {
		byKey[p.Key] = p.ID
	}
	return byKey, nil
}

func seedRoles(db *gorm.DB, permIDs map[string]uint) (map[string]uint, error) {
	// permission sets per role
	owner := []string{
		"agency.update", "agency.member.invite", "agency.member.remove",
		"tour.create", "tour.update", "tour.publish", "tour.delete",
		"tour.image.manage", "departure.manage", "price.manage",
		"booking.view", "analytics.view", "boost.purchase", "boost.view",
		"product.manage",
	}
	staff := []string{
		"tour.create", "tour.update", "tour.image.manage",
		"departure.manage", "price.manage", "booking.view",
	}
	admin := []string{} // super admin gets everything
	for k := range permIDs {
		admin = append(admin, k)
	}

	roleDefs := []struct {
		key, name string
		perms     []string
	}{
		{models.RoleSuperAdmin, "Super Admin", admin},
		{models.RoleAgencyOwner, "Agency Owner", owner},
		{models.RoleAgencyStaff, "Agency Staff", staff},
		{models.RoleTourist, "Tourist", []string{}},
	}

	roleIDs := map[string]uint{}
	for _, rd := range roleDefs {
		role := models.Role{Key: rd.key, Name: rd.name}
		if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "key"}}, DoNothing: true}).
			Create(&role).Error; err != nil {
			return nil, err
		}
		// reload to get a stable ID (OnConflict DoNothing leaves role.ID zero)
		if err := db.Where("`key` = ?", rd.key).First(&role).Error; err != nil {
			return nil, err
		}
		roleIDs[rd.key] = role.ID

		// attach permissions
		links := make([]models.RolePermission, 0, len(rd.perms))
		for _, pk := range rd.perms {
			if pid, ok := permIDs[pk]; ok {
				links = append(links, models.RolePermission{RoleID: role.ID, PermissionID: pid})
			}
		}
		if len(links) > 0 {
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&links).Error; err != nil {
				return nil, err
			}
		}
	}
	return roleIDs, nil
}

func seedSuperAdmin(db *gorm.DB, cfg *config.Config, superRoleID uint, log *slog.Logger) error {
	var existing models.User
	err := db.Where("email = ?", cfg.App.SuperEmail).First(&existing).Error
	if err == nil {
		return nil // already present
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	pw, err := hash.Password(cfg.App.SuperPass)
	if err != nil {
		return err
	}
	now := time.Now()
	admin := models.User{
		Name:            "Super Admin",
		Email:           cfg.App.SuperEmail,
		PasswordHash:    &pw,
		EmailVerifiedAt: &now,
		Locale:          cfg.App.DefaultLang,
		Status:          models.UserStatusActive,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	if err := db.Create(&models.UserRole{UserID: admin.ID, RoleID: superRoleID}).Error; err != nil {
		return err
	}
	log.Warn("super admin created — change the password", "email", cfg.App.SuperEmail)
	return nil
}

func seedTranslations(db *gorm.DB) error {
	rows := []models.Translation{
		{Locale: "en", Namespace: "common", Key: "welcome", Value: "Welcome to Zourleb"},
		{Locale: "ar", Namespace: "common", Key: "welcome", Value: "أهلاً بكم في زورليبان"},
		{Locale: "fr", Namespace: "common", Key: "welcome", Value: "Bienvenue à Zourleb"},
		{Locale: "en", Namespace: "common", Key: "book_now", Value: "Book now"},
		{Locale: "ar", Namespace: "common", Key: "book_now", Value: "احجز الآن"},
		{Locale: "fr", Namespace: "common", Key: "book_now", Value: "Réserver"},
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "locale"}, {Name: "namespace"}, {Name: "key"}},
		DoNothing: true,
	}).Create(&rows).Error
}
