package di

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/internal/infrastructure/persistent"
	"github.com/ductran999/letobserv/internal/transport/handler"
	"github.com/ductran999/letobserv/pkg/dbconn"
	"github.com/ductran999/letobserv/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductContainer struct {
	Env *configs.ProductServiceEnv
	DB  *gorm.DB

	repo    repository.ProductRepository
	useCase usecase.ProductUseCase
}

func NewProductContainer(env *configs.ProductServiceEnv) (*ProductContainer, error) {
	container := &ProductContainer{
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

func (c *ProductContainer) ConnectDB() error {
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

func (c *ProductContainer) InitRepo() {
	c.repo = persistent.NewProductRepo(c.DB)
}

func (c *ProductContainer) InitUseCase() {
	c.useCase = usecase.NewProductUseCase(c.repo)
}

func (c *ProductContainer) StartHTTPServer() error {
	if c.Env.ServiceEnv == consts.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	hdl := handler.NewProductHandler(c.useCase)

	// Simple reduce stock API
	r := gin.Default()
	r.PATCH("/api/products/:id/reduce-stock", hdl.ReduceProductStock)
	r.GET("/health", hdl.CheckHealth)

	// Start http server
	address := net.JoinHostPort("0.0.0.0", c.Env.ServicePort)
	log.Println("[INFO] product service serving on", address)
	if err := r.Run(address); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
