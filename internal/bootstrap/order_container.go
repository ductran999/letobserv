package bootstrap

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/internal/infrastructure/persistent"
	"github.com/ductran999/letobserv/internal/transport/handler"
	"github.com/ductran999/letobserv/pkg/httpclient"
	"github.com/gin-gonic/gin"
)

type OrderContainer struct {
	Env        *configs.OrdersConfigEnv
	httpClient httpclient.Client

	orderUC   usecase.OrderUseCase
	orderRepo repository.OrderRepository
}

func NewOrderContainer(env *configs.OrdersConfigEnv) (*OrderContainer, error) {
	container := &OrderContainer{
		Env:        env,
		httpClient: httpclient.New(),
	}

	container.InitRepo()
	container.InitUseCase()

	return container, nil
}

func (c *OrderContainer) InitRepo() {
	c.orderRepo = persistent.NewOrderRepo()
}

func (c *OrderContainer) InitUseCase() {
	c.orderUC = usecase.NewOrderUseCase(c.orderRepo, c.httpClient)
}

func (c *OrderContainer) StartHTTPServer() error {
	if c.Env.ServiceEnv == consts.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	hdl := handler.NewOrderHandler(c.orderUC)
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
