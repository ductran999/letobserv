package app

// import (
// 	"errors"
// 	"net"
// 	"net/http"

// 	configs "github.com/ductran999/letobserv/services/placement/internal/config"
// 	"github.com/gin-gonic/gin"
// 	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
// )

// func RunHTTP(app *bootstrap.OrderBootstrap, env *configs.OrdersConfigEnv) error {
// 	if env.ServiceEnv == consts.ProductionEnv {
// 		gin.SetMode(gin.ReleaseMode)
// 	}

// 	router := gin.Default()
// 	if env.ApmEnable {
// 		router.Use(otelgin.Middleware(env.ServiceName))
// 	}
// 	router.Use(middleware.LoggingMiddleware())

// 	handler := handler.NewOrderHandler(app.OrderUC)
// 	generated.RegisterHandlers(router, handler)

// 	address := net.JoinHostPort("0.0.0.0", env.ServicePort)
// 	if err := router.Run(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 		return err
// 	}

// 	return nil
// }
