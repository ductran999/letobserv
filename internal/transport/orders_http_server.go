package transport

import (
	"errors"
	"log"
	"net"
	"net/http"

	generated "github.com/ductran999/letobserv/api/generated/orders"
	"github.com/ductran999/letobserv/internal/bootstrap"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/internal/transport/handler"
	"github.com/gin-gonic/gin"
)

func StartOrderHTTPServer(c *bootstrap.OrderContainer) error {
	if c.Env.ServiceEnv == consts.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	hdl := handler.NewOrderHandler(c.OrderUC)

	// Simple api place an order
	r := gin.New()
	generated.RegisterHandlers(r, hdl)

	// Start http server
	address := net.JoinHostPort("0.0.0.0", c.Env.ServicePort)
	log.Println("[INFO] order service serving on", address)
	if err := r.Run(address); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
