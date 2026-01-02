package bootstrap

import (
	"fmt"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/inventory"
	infraInventory "github.com/ductran999/letobserv/internal/infrastructure/inventory"
	infraProduct "github.com/ductran999/letobserv/internal/infrastructure/product"
	"github.com/ductran999/letobserv/pkg/dbconn"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/txmanager"

	"gorm.io/gorm"
)

type InventoryBootstrap struct {
	env *configs.InventoryEnv
	db  *gorm.DB

	InventoryUC inventory.InventoryUsecase
}

func NewInventory(env *configs.InventoryEnv) (*InventoryBootstrap, error) {
	// Init db connection
	conn, err := connectInventoryDB(env)
	if err != nil {
		return nil, fmt.Errorf("init container failed: connect DB error: %w", err)
	}

	// Init logger
	if err := logger.New(logger.ServiceInfo{
		Name:    env.ServiceName,
		Version: env.ServiceVersion,
		Env:     env.ServiceEnv,
	}); err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	txmanager := txmanager.NewTransactionManager(conn)
	productRepo := infraProduct.NewProductRepo(conn)
	inventoryRepo := infraInventory.NewInventoryRepo(conn)
	inventoryUC := inventory.NewInventoryUsecase(txmanager, productRepo, inventoryRepo)

	app := &InventoryBootstrap{
		env:         env,
		db:          conn,
		InventoryUC: inventoryUC,
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
