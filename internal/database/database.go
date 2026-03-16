package database

import (
	"fmt"
	"log/slog"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/config"
	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	slog.Info("Connected to PostgreSQL", "host", cfg.DBHost, "port", cfg.DBPort, "dbname", cfg.DBName)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	slog.Info("Running database migrations...")
	err := db.AutoMigrate(
		&models.Employee{},
		&models.Client{},
		&models.Permission{},
		&models.Token{},
		// Sprint 2 models
		&models.Currency{},
		&models.SifraDelatnosti{},
		&models.SifraPlacanja{},
		&models.Firma{},
		&models.Account{},
		&models.Transfer{},
		&models.PaymentRecipient{},
		&models.Payment{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	slog.Info("Migrations complete")
	return nil
}

// SeedPermissions inserts default permissions if they don't already exist
func SeedPermissions(db *gorm.DB) error {
	if err := db.Model(&models.Permission{}).
		Where("subject_type = '' OR subject_type IS NULL").
		Update("subject_type", models.PermissionSubjectEmployee).Error; err != nil {
		return fmt.Errorf("failed to backfill permission subject types: %w", err)
	}

	for _, perm := range models.DefaultPermissions {
		p := perm
		result := db.Where(models.Permission{Name: p.Name}).Assign(models.Permission{
			Description: p.Description,
			SubjectType: p.SubjectType,
		}).FirstOrCreate(&p)
		if result.Error != nil {
			return fmt.Errorf("failed to seed permission %q: %w", p.Name, result.Error)
		}
	}
	slog.Info("Permissions seeded")
	return nil
}
