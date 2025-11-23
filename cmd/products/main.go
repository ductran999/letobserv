package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/services/common"
	"github.com/ductran999/letobserv/services/products/handler"
	"github.com/ductran999/letobserv/services/products/repo"
	"github.com/ductran999/letobserv/services/products/usecase"
	"github.com/gin-gonic/gin"
	"github.com/hyperdxio/otel-config-go/otelconfig"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Load .env file
	if err := godotenv.Load("product.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize otel config and use it across the entire app
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry()
	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	defer otelShutdown()

	// ConnectDB
	port, err := strconv.Atoi(os.Getenv("DB_PRODUCT_PORT"))
	if err != nil || port == 0 {
		log.Fatalln("Invalid DB_PRODUCT_PORT:", err)
	}
	pgConf := dbkit.PostgreSQLConfig{
		Config: dbkit.Config{
			Host:     os.Getenv("DB_PRODUCT_HOST"),
			Port:     port,
			Username: os.Getenv("DB_PRODUCT_USERNAME"),
			Password: os.Getenv("DB_PRODUCT_PASSWORD"),
			Database: os.Getenv("DB_PRODUCT_DATABASE"),
			TimeZone: "UTC",
		},
		SSLMode: "disable",
	}
	conn, err := common.ConnectDB(pgConf)
	if err != nil {
		log.Fatalln("failed when connecting to db", err)
	}

	if os.Getenv("SERVICE_ENV") == common.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	repo := repo.NewProductRepo(conn)
	uc := usecase.NewProductUseCase(repo)
	hdl := handler.NewProductHandler(uc)

	// Simple reduce stock API
	r := gin.Default()
	r.Use(otelgin.Middleware("product-service"))
	r.PATCH("/api/products/:id/reduce-stock", hdl.ReduceProductStock)
	r.GET("/health", hdl.CheckHealth)

	// Start http server
	address := "0.0.0.0:" + os.Getenv("PRODUCT_SERVICE_PORT")
	log.Println("[INFO] product service serving on", os.Getenv("PRODUCT_SERVICE_PORT"))
	if err := r.Run(address); err != nil {
		log.Println("[ERR] start product service error", err)
	}
}
