// Package database manages the GORM connection and schema migration.
package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/models"
)

// Connect opens a GORM MySQL connection with a tuned pool.
func Connect(cfg *config.Config, log *slog.Logger) (*gorm.DB, error) {
	level := gormlogger.Warn
	if cfg.App.Env == "development" {
		level = gormlogger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{
		Logger:                                   gormlogger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	log.Info("database connected", "host", cfg.DB.Host, "name", cfg.DB.Name)
	return db, nil
}

// Migrate runs GORM AutoMigrate across all models. For production a versioned
// migration tool (golang-migrate) is preferable; AutoMigrate keeps dev fast.
func Migrate(db *gorm.DB, log *slog.Logger) error {
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	log.Info("schema migrated", "models", len(models.AllModels()))
	return nil
}
