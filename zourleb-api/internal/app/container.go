// Package app wires dependencies (DI) and exposes them to the router.
package app

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/handler"
	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/cache"
	"github.com/zourleb/zourleb-api/pkg/mailer"
	"github.com/zourleb/zourleb-api/pkg/oauth"
	"github.com/zourleb/zourleb-api/pkg/onesignal"
	"github.com/zourleb/zourleb-api/pkg/otp"
	"github.com/zourleb/zourleb-api/pkg/payment"
	"github.com/zourleb/zourleb-api/pkg/sms"
	"github.com/zourleb/zourleb-api/pkg/storage"
	"github.com/zourleb/zourleb-api/pkg/token"
	"github.com/zourleb/zourleb-api/pkg/whatsapp"
)

// Container holds constructed handlers, services and middleware for the router.
type Container struct {
	Cfg *config.Config
	Log *slog.Logger

	// shared services
	I18n     *service.I18nService
	Settings *service.SettingsService
	Notify   *service.NotificationService

	// middleware
	Auth   *middleware.Auth
	Locale *service.I18nService

	// handlers
	Health  *handler.HealthHandler
	AuthH   *handler.AuthHandler
	Account *handler.AccountHandler
	I18nH   *handler.I18nHandler
	OTP     *handler.OTPHandler
	Catalog *handler.CatalogHandler
	Booking *handler.BookingHandler
	Engage  *handler.EngagementHandler
	Agency  *handler.AgencyHandler
	Boost   *handler.BoostHandler
	Admin   *handler.AdminHandler
	Upload  *handler.UploadHandler
	Notif   *handler.NotificationHandler
	Shop    *handler.ShopHandler
}

// New constructs the full dependency graph.
func New(cfg *config.Config, db *gorm.DB, log *slog.Logger) *Container {
	// --- repositories ---
	userRepo := repository.NewUserRepository(db)
	i18nRepo := repository.NewI18nRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	catalogRepo := repository.NewCatalogRepository(db)
	tourRepo := repository.NewTourRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	engageRepo := repository.NewEngagementRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	agencyRepo := repository.NewAgencyRepository(db)
	boostRepo := repository.NewBoostRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	shopRepo := repository.NewShopRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// --- infrastructure ---
	tokens := token.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	google := oauth.NewGoogleVerifier(cfg.Google.ClientID)

	store := buildCache(cfg, log)
	mail := buildMailer(cfg, log)
	dev := !cfg.IsProduction()
	waClient := whatsapp.NewClient(cfg.WhatsApp.PhoneID, cfg.WhatsApp.Token, cfg.WhatsApp.TemplateOTP, dev, log)
	smsClient := sms.NewClient(cfg.SMS.Provider, cfg.SMS.From, cfg.SMS.APIKey, cfg.SMS.APISecret, dev, log)
	dispatcher := otp.NewDispatcher(waClient, smsClient)
	pushClient := onesignal.NewClient(cfg.OneSig.AppID, cfg.OneSig.APIKey, log)
	payReg := payment.NewRegistry("manual", payment.NewManualProvider())
	store0 := buildStorage(cfg, log)

	// --- services ---
	settingsSvc := service.NewSettingsService(i18nRepo)
	auditSvc := service.NewAuditService(auditRepo)
	authSvc := service.NewAuthService(userRepo, tokens, google, cfg, store, mail)
	accountSvc := service.NewAccountService(userRepo, deviceRepo)
	i18nSvc := service.NewI18nService(i18nRepo, cfg.App.DefaultLang)
	otpSvc := service.NewOTPService(otpRepo, userRepo, store, dispatcher, cfg.OTP)
	catalogSvc := service.NewCatalogService(catalogRepo, settingsSvc, cfg.App.DefaultLang)
	engageSvc := service.NewEngagementService(engageRepo, settingsSvc)
	notifySvc := service.NewNotificationService(notifRepo, pushClient)
	bookingSvc := service.NewBookingService(bookingRepo, tourRepo, userRepo, paymentRepo, payReg, settingsSvc, notifySvc, cfg.App.DefaultLang)
	agencySvc := service.NewAgencyService(agencyRepo, tourRepo, bookingRepo, userRepo, notifySvc, cfg.App.DefaultLang)
	boostSvc := service.NewBoostService(boostRepo, paymentRepo, agencyRepo, payReg, notifySvc)
	uploadSvc := service.NewUploadService(store0)
	adminSvc := service.NewAdminService(adminRepo, agencyRepo, userRepo, i18nRepo, engageRepo, boostRepo, paymentRepo, settingsSvc, boostSvc, auditSvc)
	shopSvc := service.NewShopService(shopRepo, cfg.App.DefaultLang)

	// --- middleware ---
	authMw := middleware.NewAuth(tokens, userRepo)

	return &Container{
		Cfg:      cfg,
		Log:      log,
		I18n:     i18nSvc,
		Settings: settingsSvc,
		Notify:   notifySvc,
		Auth:     authMw,
		Locale:   i18nSvc,
		Health:   handler.NewHealthHandler(db),
		AuthH:    handler.NewAuthHandler(authSvc),
		Account:  handler.NewAccountHandler(accountSvc),
		I18nH:    handler.NewI18nHandler(i18nSvc),
		OTP:      handler.NewOTPHandler(otpSvc),
		Catalog:  handler.NewCatalogHandler(catalogSvc),
		Booking:  handler.NewBookingHandler(bookingSvc),
		Engage:   handler.NewEngagementHandler(engageSvc),
		Agency:   handler.NewAgencyHandler(agencySvc, uploadSvc),
		Boost:    handler.NewBoostHandler(boostSvc),
		Admin:    handler.NewAdminHandler(adminSvc),
		Upload:   handler.NewUploadHandler(uploadSvc),
		Notif:    handler.NewNotificationHandler(notifySvc),
		Shop:     handler.NewShopHandler(shopSvc),
	}
}

// buildCache returns a Redis store when enabled, else an in-memory store.
func buildCache(cfg *config.Config, log *slog.Logger) cache.Store {
	if cfg.Redis.Enabled {
		if s, err := cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB); err == nil {
			log.Info("cache: redis", "addr", cfg.Redis.Addr)
			return s
		}
		log.Warn("cache: redis unavailable, falling back to in-memory")
	}
	log.Info("cache: in-memory")
	return cache.NewMemory()
}

// buildStorage returns the configured storage driver (local by default).
func buildStorage(cfg *config.Config, log *slog.Logger) storage.Storage {
	s, err := storage.NewLocal(cfg.Storage.LocalPath, cfg.Storage.PublicURL)
	if err != nil {
		log.Error("storage init failed", "err", err)
	}
	return s
}

// buildMailer returns an SMTP client when configured, otherwise a logging sender.
func buildMailer(cfg *config.Config, log *slog.Logger) mailer.Sender {
	if cfg.SMTP.Host != "" {
		return mailer.NewSMTPClient(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From, cfg.SMTP.User, cfg.SMTP.Password)
	}
	return mailer.NewLogSender(log)
}
