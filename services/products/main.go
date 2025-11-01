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
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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
	defer conn.Close() //nolint

	if os.Getenv("SERVICE_ENV") == common.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	repo := repo.NewProductRepo(conn.DB())
	uc := usecase.NewProductUseCase(repo)
	hdl := handler.NewProductHandler(uc)

	// Simple reduce stock API
	r := gin.Default()
	r.PATCH("/api/products/:id/reduce-stock", hdl.ReduceProductStock)
	r.GET("/health", hdl.CheckHealth)

	// Start http server
	address := "0.0.0.0:" + os.Getenv("PRODUCT_SERVICE_PORT")
	log.Println("[INFO] product service serving on", os.Getenv("PRODUCT_SERVICE_PORT"))
	if err := r.Run(address); err != nil {
		log.Println("[ERR] start product service error", err)
	}
}
