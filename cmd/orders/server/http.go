package server

import (
	"errors"
	"log"
	"net"
	"net/http"

	generated "github.com/ductran999/letobserv/api/generated/orders"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/bootstrap"
	"github.com/ductran999/letobserv/internal/consts"
	handler "github.com/ductran999/letobserv/internal/transport/order"
	"github.com/gin-gonic/gin"
)

func RunHTTP(app *bootstrap.OrderBootstrap, env *configs.OrdersConfigEnv) error {
	if env.ServiceEnv == consts.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	handler := handler.NewOrderHandler(app.OrderUC)
	generated.RegisterHandlers(r, handler)

	address := net.JoinHostPort("0.0.0.0", env.ServicePort)

	log.Println("[INFO] order service serving on", address)

	if err := r.Run(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
