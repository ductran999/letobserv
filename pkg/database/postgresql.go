package database

import (
	"fmt"

	"github.com/ductran999/dbkit"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PgConfig = dbkit.PostgreSQLConfig
type Config = dbkit.Config

func ConnectDB(pgConf PgConfig) (*gorm.DB, error) {
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := db.Use(tracing.NewPlugin(
		tracing.WithoutQueryVariables())); err != nil {
		return nil, fmt.Errorf("applied tracing failed: %w", err)
	}

	return db, nil
}

func MustConnect(cfg PgConfig) *gorm.DB {
	db, err := ConnectDB(cfg)
	if err != nil {
		panic(err)
	}

	return db
}
