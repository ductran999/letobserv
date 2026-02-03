package apm

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type otelAgent struct {
	tp *sdktrace.TracerProvider
	mp *metric.MeterProvider
}

func (a *otelAgent) Shutdown(ctx context.Context) error {
	errMetric := a.mp.Shutdown(ctx)
	errTrace := a.tp.Shutdown(ctx)

	if errTrace != nil {
		return fmt.Errorf("trace shutdown error: %w", errTrace)
	}
	if errMetric != nil {
		return fmt.Errorf("metric shutdown error: %w", errMetric)
	}

	return nil
}

type noopAgent struct{}

func (a *noopAgent) Shutdown(ctx context.Context) error {
	return nil
}

func NewAgent(ctx context.Context, cfg AgentConfig) (APMAgent, error) {
	if !cfg.Enabled {
		return &noopAgent{}, nil
	}

	res, err := initResource(ctx, cfg.ServiceInfo)
	if err != nil {
		return nil, fmt.Errorf("new resource failed: %w", err)
	}

	conn, err := grpc.NewClient(cfg.ExporterEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		sdktrace.WithResource(res),
	)

	// Makes the tracer available to instrumentation libraries
	// Must set to allow user create new span any where in app	otel.Tracer("")
	otel.SetTracerProvider(tp)

	// Propagates trace context across service boundaries using W3C standards
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("new otlp metric grpc exporter failed: %w", err)
	}
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(30*time.Second))),
	)

	// Set the global MeterProvider.
	// Must set to allow user expose metric any where in app by otel.Meter("")
	otel.SetMeterProvider(mp)

	return &otelAgent{
		tp: tp,
		mp: mp,
	}, nil
}

func initResource(ctx context.Context, serviceInfo ServiceInfo) (*resource.Resource, error) {
	// Create the resource with attributes from environment (OTEL_RESOURCE_ATTRIBUTES)
	// and default host/process attributes
	return resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithAttributes(
			semconv.ServiceName(serviceInfo.Name),
			semconv.DeploymentEnvironment(serviceInfo.Env),
			semconv.ServiceVersion(serviceInfo.Version),
			semconv.ProcessRuntimeName("go"),
			semconv.ProcessRuntimeVersion(runtime.Version()),
		),
	)
}
