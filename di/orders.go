package di

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/services/common"
	"github.com/ductran999/letobserv/services/orders/handler"
	"github.com/ductran999/letobserv/services/orders/port"
	"github.com/ductran999/letobserv/services/orders/repo"
	"github.com/ductran999/letobserv/services/orders/usecase"
	"github.com/gin-gonic/gin"
)

type OrderContainer struct {
	Env *configs.OrdersConfigEnv

	repo    port.OrderRepo
	useCase port.OrderUseCase
}

func NewOrderContainer(env *configs.OrdersConfigEnv) (*OrderContainer, error) {
	container := &OrderContainer{
		Env: env,
	}
	// Initialize logger with the multi-writer
	if err := logger.New(
		logger.ServiceInfo{
			Name:    env.ServiceName,
			Version: env.ServiceVersion,
			Env:     env.ServiceEnv,
		},
	); err != nil {
		return nil, err
	}

	container.InitRepo()
	container.InitUseCase()

	return container, nil
}

func (c *OrderContainer) InitRepo() {
	c.repo = repo.NewOrderRepo()
}

func (c *OrderContainer) InitUseCase() {
	c.useCase = usecase.NewOrderUseCase(c.repo)
}

func (c *OrderContainer) StartHTTPServer() error {
	if c.Env.ServiceEnv == common.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	hdl := handler.NewOrderHandler(c.useCase)
	// Simple api place an order
	r := gin.New()

	r.GET("/health", hdl.CheckHealth)
	r.POST("/orders", hdl.PlaceOrder)

	// Start http server
	address := net.JoinHostPort("0.0.0.0", c.Env.ServicePort)
	log.Println("[INFO] order service serving on", address)
	if err := r.Run(address); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
