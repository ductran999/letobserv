package main

import (
	"log"
	"os"

	"github.com/ductran999/letobserv/services/orders/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	hdl := handler.NewOrderHandler(nil)

	// Simple api place an order
	r := gin.Default()
	r.GET("/health", hdl.CheckHealth)

	// Start http server
	address := "0.0.0.0:" + os.Getenv("ORDER_SERVICE_PORT")
	log.Println("[INFO] product service serving on", os.Getenv("ORDER_SERVICE_PORT"))
	if err := r.Run(address); err != nil {
		log.Println("[ERR] start product service error", err)
	}
}
