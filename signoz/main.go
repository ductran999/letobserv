package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func tracerApp() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	r := gin.Default()

	// Add OpenTelemetry middleware to instrument all routes
	r.Use(otelgin.Middleware("my-gin-service"))

	r.GET("/", func(c *gin.Context) {
		// Extract the span created by otelgin middleware
		ctx := c.Request.Context()
		span := trace.SpanFromContext(ctx)

		span.SetAttributes(attribute.String("custom.key", "value"))
		span.AddEvent("manual.event")

		c.JSON(200, gin.H{"message": "Hello traced Gin!"})
	})

	log.Println("Server starting on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalln(err)
	}
}
func metricApp() {
	ctx := context.Background()

	// ---------------------------------------------------------
	// BƯỚC 1: Khởi tạo Meter Provider (CHỈ LÀM 1 LẦN DUY NHẤT)
	// ---------------------------------------------------------
	shutdown, err := initMeterProvider(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize meter provider: %v", err)
	}
	// Đảm bảo flush metric khi ứng dụng dừng
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}()

	if err := runtime.Start(
		runtime.WithMeterProvider(otel.GetMeterProvider()),
		runtime.WithMinimumReadMemStatsInterval(2*time.Second), 
	); err != nil {
		log.Fatal(err)
	}
	log.Println("Runtime metrics collection started (go.memory.used, etc)...")

	meter := otel.Meter("my-app-meter")

	reqCounter, err := meter.Int64Counter(
		"http.requests",
		metric.WithDescription("Total request count"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Simulate work
		time.Sleep(50 * time.Millisecond)
		w.Write([]byte("ok"))

		// Ghi nhận metric sau khi xử lý xong
		reqCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("route", "/hello"),
			attribute.Int("status", 200),
		))

		log.Printf("Request completed in %v", time.Since(start))
	})

	log.Println("Server listening on :8081...")

	// ListenAndServe là hàm blocking, nó sẽ chạy mãi ở đây cho đến khi lỗi
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	metricApp()
}
