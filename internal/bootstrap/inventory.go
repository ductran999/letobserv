package bootstrap

import (
	"fmt"
	"log"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/usecase/inventory"
	infraInventory "github.com/ductran999/letobserv/internal/infrastructure/inventory"
	infraProduct "github.com/ductran999/letobserv/internal/infrastructure/product"
	"github.com/ductran999/letobserv/pkg/apm"
	"github.com/ductran999/letobserv/pkg/dbconn"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/txmanager"

	"gorm.io/gorm"
)

type InventoryBootstrap struct {
	env *configs.InventoryEnv
	db  *gorm.DB

	APMAgent    apm.APMAgent
	InventoryUC inventory.InventoryUsecase
}

func NewInventory(env *configs.InventoryEnv) (*InventoryBootstrap, error) {
	// Init db connection
	conn, err := connectInventoryDB(env)
	if err != nil {
		return nil, fmt.Errorf("init container failed: connect DB error: %w", err)
	}
	log.Println("[INFO] connect pg db successfully!")

	// Init logger
	var apmAgent apm.APMAgent
	if env.ApmEnable {
		apmAgent, err = newInventoryServiceAPMAgent(env)
		if err != nil {
			return nil, fmt.Errorf("failed to start apm agent: %w", err)
		}
		log.Println("[INFO] start apm agent successfully!")
	}

	if err = newInventoryLogger(env); err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	log.Println("[INFO] initialize logger successfully!")

	txmanager := txmanager.NewTransactionManager(conn)
	productRepo := infraProduct.NewProductRepo(conn)
	inventoryRepo := infraInventory.NewInventoryRepo(conn)
	inventoryUC := inventory.NewInventoryUsecase(txmanager, productRepo, inventoryRepo)

	app := &InventoryBootstrap{
		env:         env,
		db:          conn,
		InventoryUC: inventoryUC,
		APMAgent:    apmAgent,
	}

	return app, nil
}

func connectInventoryDB(env *configs.InventoryEnv) (*gorm.DB, error) {
	pgConf := dbkit.PostgreSQLConfig{
		Config: dbkit.Config{
			Host:     env.PgHost,
			Port:     env.PgPort,
			Username: env.PgUser,
			Password: env.PgPass,
			Database: env.PgDB,
			TimeZone: "UTC",
		},
		SSLMode: "disable",
	}

	return dbconn.ConnectDB(pgConf)
}

func newInventoryServiceAPMAgent(env *configs.InventoryEnv) (apm.APMAgent, error) {
	config := apm.AgentConfig{
		ServiceInfo: apm.ServiceInfo{
			Name:    env.ServiceName,
			Version: env.ServiceVersion,
			Env:     env.ServiceEnv,
		},
		APIKey:           env.ApmApiKey,
		InsecureMode:     env.ApmInsecureMode,
		ExporterEndpoint: env.ApmExporterEndpoint,
	}
	agent, err := apm.NewClickStackAPM(config)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func newInventoryLogger(env *configs.InventoryEnv) error {
	serviceInfo := logger.ServiceInfo{
		Name:    env.ServiceName,
		Version: env.ServiceVersion,
		Env:     env.ServiceEnv,
	}

	if !env.ApmEnable {
		return logger.New(serviceInfo)
	}

	apmConfig := logger.APMConfig{
		Enable:           env.ApmEnable,
		ExporterEndpoint: env.ApmExporterEndpoint,
		ApiKey:           env.ApmApiKey,
	}

	return logger.New(serviceInfo, logger.WithAPM(apmConfig))
}
