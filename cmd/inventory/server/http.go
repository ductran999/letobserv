package server

import (
	"errors"
	"log"
	"net"
	"net/http"

	generated "github.com/ductran999/letobserv/api/generated/inventory"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/bootstrap"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/internal/transport/inventory"
	"github.com/gin-gonic/gin"
)

func RunHTTP(env *configs.InventoryEnv, app *bootstrap.InventoryBootstrap) error {
	if env.ServiceEnv == consts.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	hdl := inventory.NewInventoryHandler(app.InventoryUC)

	r := gin.Default()
	generated.RegisterHandlers(r, hdl)

	// Start http server
	address := net.JoinHostPort("0.0.0.0", env.ServicePort)
	log.Println("[INFO] product service serving on", address)
	if err := r.Run(address); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
