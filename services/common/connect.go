package common

import (
	"fmt"

	"github.com/ductran999/dbkit"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func ConnectDB(pgConf dbkit.PostgreSQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		pgConf.Host,
		pgConf.Username,
		pgConf.Password,
		pgConf.Database,
		pgConf.Port,
		pgConf.SSLMode,
		pgConf.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := db.Use(tracing.NewPlugin(
		tracing.WithoutQueryVariables())); err != nil {
		return nil, fmt.Errorf("applied tracing failed: %w", err)
	}

	return db, nil
}
