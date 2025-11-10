package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/middleware"
	"github.com/ductran999/letobserv/services/orders/handler"
	"github.com/ductran999/letobserv/services/orders/repo"
	"github.com/ductran999/letobserv/services/orders/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func initTracer() func(context.Context) error {
	serviceName := os.Getenv("ORDER_SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure := os.Getenv("INSECURE_MODE")

	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	// Config exporter
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(secureOption, otlptracegrpc.WithEndpoint(collectorURL)),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	// Config resource
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKLanguageGo),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

func main() {
	// Initialize logger with the multi-writer
	if err := logger.New(
		logger.ServiceInfo{Name: "orders-service", Version: "v1.0.0", Env: "dev"},
		logger.EnableFileLogging(true),
	); err != nil {
		log.Fatalln("create logger error", err)
	}

	cleanup := initTracer()
	defer cleanup(context.Background())

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
