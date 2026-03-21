package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ductran999/letobserv/pkg/apm"
	"github.com/ductran999/letobserv/pkg/database"
	"github.com/ductran999/letobserv/services/inventory/internal/config"

	"gorm.io/gorm"
)

type App struct {
	env *config.Config
	db  *gorm.DB

	APMAgent apm.APMAgent

}

func NewApp(cfg *config.Config) (*App, error) {
	// Init db connection
	db := database.MustConnect(cfg.PgConfig())
	slog.Info("connect pg db successfully!")

	var apmAgent apm.APMAgent
	var err error
	if cfg.ApmEnable {
		apmAgent, err = apm.NewAgent(context.Background(), cfg.APMConfig())
		if err != nil {
			return nil, fmt.Errorf("failed to start apm agent: %w", err)
		}
		slog.Info("apm agent running!")
	}

	// if err = newInventoryLogger(env); err != nil {
	// 	return nil, fmt.Errorf("init logger: %w", err)
	// }
	// log.Println("[INFO] initialize logger successfully!")

	app := &App{
		env:      cfg,
		db:       db,
		APMAgent: apmAgent,
	}

	return app, nil
}

// func newInventoryLogger(env *configs.InventoryEnv) error {
// 	serviceInfo := logger.ServiceInfo{
// 		Name:    env.ServiceName,
// 		Version: env.ServiceVersion,
// 		Env:     env.ServiceEnv,
// 	}

// 	if !env.ApmEnable {
// 		return logger.New(serviceInfo)
// 	}

// 	apmConfig := logger.APMConfig{
// 		Enable:           env.ApmEnable,
// 		ExporterEndpoint: env.ApmExporterEndpoint,
// 		ApiKey:           env.ApmApiKey,
// 	}

// 	return logger.New(serviceInfo, logger.WithAPM(apmConfig))
// }
