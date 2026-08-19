package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration, loaded from environment variables.
// Defaults are sane for local development; production must override secrets.
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Google   GoogleConfig
	OneSig   OneSignalConfig
	WhatsApp WhatsAppConfig
	SMS      SMSConfig
	Storage  StorageConfig
	OTP      OTPConfig
	SMTP     SMTPConfig
}

type AppConfig struct {
	Name        string
	Env         string // development | staging | production
	DefaultLang string
	SuperEmail  string // bootstrap super admin email
	SuperPass   string // bootstrap super admin password
}

type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	AllowOrigins    []string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Params   string
}

func (d DBConfig) DSN() string {
	// GORM MySQL DSN: user:pass@tcp(host:port)/dbname?params
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Params)
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

type GoogleConfig struct {
	ClientID string
}

type OneSignalConfig struct {
	AppID  string
	APIKey string
}

type WhatsAppConfig struct {
	Provider    string // meta | twilio
	PhoneID     string
	Token       string
	TemplateOTP string
}

type SMSConfig struct {
	Provider  string // twilio | vonage | local
	From      string
	APIKey    string
	APISecret string
}

type StorageConfig struct {
	Driver     string // local | s3
	LocalPath  string
	PublicURL  string
	S3Bucket   string
	S3Region   string
	S3Key      string
	S3Secret   string
	S3Endpoint string
}

type OTPConfig struct {
	TTL          time.Duration
	Length       int
	MaxAttempts  int
	ResendWindow time.Duration
	DailyCap     int
}

type SMTPConfig struct {
	Host     string
	Port     string
	From     string
	User     string
	Password string
	BaseURL  string
	TTL      time.Duration
}

// Load reads configuration from the environment, applying defaults.
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:        env("APP_NAME", "Zourleb"),
			Env:         env("APP_ENV", "development"),
			DefaultLang: env("APP_DEFAULT_LANG", "ar"),
			SuperEmail:  env("SUPER_ADMIN_EMAIL", "admin@zourleb.com"),
			SuperPass:   env("SUPER_ADMIN_PASSWORD", "ChangeMe123!"),
		},
		HTTP: HTTPConfig{
			Port:            env("HTTP_PORT", "8080"),
			ReadTimeout:     envDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    envDuration("HTTP_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: envDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			AllowOrigins:    envSlice("HTTP_ALLOW_ORIGINS", []string{"*"}),
		},
		DB: DBConfig{
			Host:     env("DB_HOST", "127.0.0.1"),
			Port:     env("DB_PORT", "3306"),
			User:     env("DB_USER", "root"),
			Password: env("DB_PASSWORD", ""),
			Name:     env("DB_NAME", "zourleb"),
			Params:   env("DB_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Redis: RedisConfig{
			Addr:     env("REDIS_ADDR", "127.0.0.1:6379"),
			Password: env("REDIS_PASSWORD", ""),
			DB:       envInt("REDIS_DB", 0),
			Enabled:  envBool("REDIS_ENABLED", false),
		},
		JWT: JWTConfig{
			Secret:     env("JWT_SECRET", "dev-insecure-secret-change-me"),
			AccessTTL:  envDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: envDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
			Issuer:     env("JWT_ISSUER", "zourleb"),
		},
		Google: GoogleConfig{
			ClientID: env("GOOGLE_CLIENT_ID", ""),
		},
		OneSig: OneSignalConfig{
			AppID:  env("ONESIGNAL_APP_ID", ""),
			APIKey: env("ONESIGNAL_API_KEY", ""),
		},
		WhatsApp: WhatsAppConfig{
			Provider:    env("WHATSAPP_PROVIDER", "meta"),
			PhoneID:     env("WHATSAPP_PHONE_ID", ""),
			Token:       env("WHATSAPP_TOKEN", ""),
			TemplateOTP: env("WHATSAPP_TEMPLATE_OTP", "otp_code"),
		},
		SMS: SMSConfig{
			Provider:  env("SMS_PROVIDER", "local"),
			From:      env("SMS_FROM", "Zourleb"),
			APIKey:    env("SMS_API_KEY", ""),
			APISecret: env("SMS_API_SECRET", ""),
		},
		Storage: StorageConfig{
			Driver:     env("STORAGE_DRIVER", "local"),
			LocalPath:  env("STORAGE_LOCAL_PATH", "./uploads"),
			PublicURL:  env("STORAGE_PUBLIC_URL", "http://localhost:8080/uploads"),
			S3Bucket:   env("S3_BUCKET", ""),
			S3Region:   env("S3_REGION", ""),
			S3Key:      env("S3_KEY", ""),
			S3Secret:   env("S3_SECRET", ""),
			S3Endpoint: env("S3_ENDPOINT", ""),
		},
		OTP: OTPConfig{
			TTL:          envDuration("OTP_TTL", 5*time.Minute),
			Length:       envInt("OTP_LENGTH", 6),
			MaxAttempts:  envInt("OTP_MAX_ATTEMPTS", 3),
			ResendWindow: envDuration("OTP_RESEND_WINDOW", 60*time.Second),
			DailyCap:     envInt("OTP_DAILY_CAP", 10),
		},
		SMTP: SMTPConfig{
			Host:     env("SMTP_HOST", ""),
			Port:     env("SMTP_PORT", "587"),
			From:     env("SMTP_FROM", "noreply@zourleb.com"),
			User:     env("SMTP_USER", ""),
			Password: env("SMTP_PASSWORD", ""),
			BaseURL:  env("SMTP_BASE_URL", "http://localhost:8080"),
			TTL:      envDuration("SMTP_TTL", 24*time.Hour),
		},
	}
}

func (c *Config) IsProduction() bool { return c.App.Env == "production" }

// --- helpers ---

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func envSlice(key string, def []string) []string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return def
}
