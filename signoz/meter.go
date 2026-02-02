package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func initMeterProvider(ctx context.Context) (func(context.Context) error, error) {
	// Create the OTLP exporter
	// It will automatically use OTEL_EXPORTER_OTLP_METRICS_ENDPOINT and headers from env
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("new otlp metric grpc exporter failed: %w", err)
	}

	// Create the resource with attributes from environment (OTEL_RESOURCE_ATTRIBUTES)
	// and default host/process attributes
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithAttributes(
			semconv.ServiceName("auth-service"),
			semconv.DeploymentEnvironment("dev"),
			semconv.ServiceVersion("v1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new resource failed: %w", err)
	}

	// Create the MeterProvider with the exporter and resource
	// Set a periodic reader to export metrics every 10 seconds
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(5*time.Second))),
	)

	// Set the global MeterProvider
	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}
