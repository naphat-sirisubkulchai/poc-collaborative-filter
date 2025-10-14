package database

import (
	"fmt"

	"poc-collaborative-filter/internal/config"
	"poc-collaborative-filter/internal/domain"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.Database.GetDSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.GetDSN())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info("Connected to database", zap.String("driver", cfg.Database.Driver))

	return db, nil
}

// AutoMigrate runs automatic migrations for all domain models
func AutoMigrate(db *gorm.DB, log *zap.Logger) error {
	log.Info("Running database migrations...")

	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Warn("Failed to create uuid-ossp extension (may already exist or not be PostgreSQL)", zap.Error(err))
	}

	if err := db.AutoMigrate(
		&domain.Customer{},
		&domain.Recommendation{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("Database migrations completed successfully")
	return nil
}
