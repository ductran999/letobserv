package main

import (
	"log"
	"os"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/middleware"
	"github.com/ductran999/letobserv/services/orders/handler"
	"github.com/ductran999/letobserv/services/orders/repo"
	"github.com/ductran999/letobserv/services/orders/usecase"
	"github.com/gin-gonic/gin"
	"github.com/hyperdxio/otel-config-go/otelconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	config, err := configs.LoadConfig("order.env")
	if err != nil {
		log.Fatalln("failed to load product config:", err)
	}

	// Initialize logger with the multi-writer
	if err := logger.New(
		logger.ServiceInfo{Name: config.ServiceName, Version: config.ServiceVersion, Env: config.ServiceEnv},
		logger.EnableFileLogging(true),
	); err != nil {
		log.Fatalln("create logger error", err)
	}

	// Initialize otel config and use it across the entire app
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry()
	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	defer otelShutdown()

	rep := repo.NewOrderRepo()
	uc := usecase.NewOrderUseCase(rep)
	hdl := handler.NewOrderHandler(uc)
	// Simple api place an order
	r := gin.New()

	r.Use(gin.Recovery(), otelgin.Middleware(os.Getenv("ORDER_SERVICE_NAME")))
	r.Use(middleware.LoggingMiddleware())
	r.GET("/health", hdl.CheckHealth)
	r.POST("/orders", hdl.PlaceOrder)

	// Start http server
	address := "0.0.0.0:" + os.Getenv("ORDER_SERVICE_PORT")
	log.Println("[INFO] product service serving on", os.Getenv("ORDER_SERVICE_PORT"))
	if err := r.Run(address); err != nil {
		log.Println("[ERR] start product service error", err)
	}
}
