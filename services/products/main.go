package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/services/common"
	"github.com/ductran999/letobserv/services/products/handler"
	"github.com/ductran999/letobserv/services/products/repo"
	"github.com/ductran999/letobserv/services/products/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc/credentials"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

var (
	serviceName  = os.Getenv("PRODUCT_SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func initTracer() func(context.Context) error {
	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.TelemetrySDKLanguageGo,
		),
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
	cleanup := initTracer()
	defer cleanup(context.Background())

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
	r.Use(otelgin.Middleware(serviceName))
	r.PATCH("/api/products/:id/reduce-stock", hdl.ReduceProductStock)
	r.GET("/health", hdl.CheckHealth)

	// Start http server
	address := "0.0.0.0:" + os.Getenv("PRODUCT_SERVICE_PORT")
	log.Println("[INFO] product service serving on", os.Getenv("PRODUCT_SERVICE_PORT"))
	if err := r.Run(address); err != nil {
		log.Println("[ERR] start product service error", err)
	}
}
