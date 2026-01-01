package bootstrap

import (
	"fmt"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/internal/infrastructure/persistent"
	"github.com/ductran999/letobserv/pkg/dbconn"
	"github.com/ductran999/letobserv/pkg/logger"

	"gorm.io/gorm"
)

type Inventory struct {
	Env *configs.InventoryEnv
	DB  *gorm.DB

	ProductRepo repository.ProductRepository
	InventoryUC usecase.InventoryUsecase
}

func NewInventory(env *configs.InventoryEnv) (*Inventory, error) {
	container := &Inventory{
		Env: env,
	}

	// Init db connection
	if err := container.ConnectDB(); err != nil {
		return nil, fmt.Errorf("init container failed: connect DB error: %w", err)
	}

	if err := logger.New(logger.ServiceInfo{
		Name:    env.ServiceName,
		Version: env.ServiceVersion,
		Env:     env.ServiceEnv,
	}); err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	container.InitRepo()
	container.InitUseCase()

	return container, nil
}

func (c *Inventory) ConnectDB() error {
	var err error

	pgConf := dbkit.PostgreSQLConfig{
		Config: dbkit.Config{
			Host:     c.Env.PgHost,
			Port:     c.Env.PgPort,
			Username: c.Env.PgUser,
			Password: c.Env.PgPass,
			Database: c.Env.PgDB,
			TimeZone: "UTC",
		},
		SSLMode: "disable",
	}

	c.DB, err = dbconn.ConnectDB(pgConf)
	return err
}

func (c *Inventory) InitRepo() {
	c.ProductRepo = persistent.NewProductRepo(c.DB)
}

func (c *Inventory) InitUseCase() {
	c.InventoryUC = usecase.NewInventoryUsecase(c.ProductRepo)
}
