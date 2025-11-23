package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ductran999/letobserv/pkg/apm"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Load .env file
	if err := godotenv.Load("product.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize otel config and use it across the entire app
	apmAgent, err := apm.NewClickStackAPM(apm.AgentConfig{
		ExporterEndpoint: "http://localhost:4318",
		InsecureMode:     true,
		LogLevel:         "debug",
		ServiceName:      "some",
		ServiceVersion:   "v1.0.0",
		ServiceEnv:       "dev",
		ApiKey:           os.Getenv("APM_API_KEY"),
	})
	if err != nil {
		log.Fatalln("failed to start apm agent:", err)
	}
	defer apmAgent.Shutdown()

	if err := logger.New(logger.ServiceInfo{
		Name:    "some",
		Version: "v1.0.0",
		Env:     "prod",
	},
	); err != nil {
		log.Fatalln("failed to init logger:", err)
	}

	// Create a new Gin router
	router := gin.Default()

	router.Use(otelgin.Middleware("some-service"))

	// Define a route that responds to GET requests on the root URL
	router.GET("/", func(c *gin.Context) {
		// logger.Logger(c.Request.Context()).Error().Msg("Hello world")
		logger.Logger(c).Error().Ctx(c).Str("string", "string-value").Msg("Integrate")
		c.String(http.StatusOK, "Hello World!")
	})

	// Run the server on port 7777
	router.Run(":7777")
}
