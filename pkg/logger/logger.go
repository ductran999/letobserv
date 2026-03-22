package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperdxio/opentelemetry-go/otelzerolog"
	"github.com/hyperdxio/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/hyperdxio/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdk "github.com/hyperdxio/opentelemetry-logs-go/sdk/logs"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/natefinch/lumberjack.v2"
)

var gLogger zerolog.Logger

// New initializes and configures the global logger instance.
// By default, logging to file is disabled.
// It sets up the global logger `gLogger` with timestamp and service context fields.
//
// Example default log output:
//
//	{
//	  "level":"info",
//	  "service_name":"order-service",
//	  "service_version":"v1.0.0",
//	  "service_env":"dev",
//	  "time":"2025-11-09T17:04:40+07:00",
//	  "message":"Application started"
//	}
//
// Options can override defaults, e.g., enabling file logging, rotation, or compression.
func New(serviceInfo ServiceInfo, options ...ConfigOption) (*zerolog.Logger, error) {
	// Validate basic service metadata
	if err := serviceInfo.Validate(); err != nil {
		return nil, err
	}

	// Initialize configuration with default values
	conf := getDefaultConfig(serviceInfo)

	// Apply user-provided configuration options
	for _, opt := range options {
		if err := opt(conf); err != nil {
			return nil, err
		}
	}

	// Prepare the log writers (e.g., console, file)
	writers, err := setupWriter(conf)
	if err != nil {
		return nil, err
	}

	// Create zerolog instance with context fields
	l := zerolog.New(writers).With().
		Timestamp().
		Str("service_name", conf.ServiceInfo.Name).
		Str("service_version", conf.ServiceInfo.Version).
		Str("service_env", conf.ServiceInfo.Env).
		Logger()

	if conf.EnableAPM {
		hook, err := setupAPM(conf)
		if err != nil {
			return nil, err
		}
		l.Hook(hook) // Add telemetry Hook
		log.Println("[INFO] hook apm to logger")
	}

	gLogger = l

	return &l, nil
}

// Logger returns a logger derived from the global logger gLogger,
// with trace information attached from the context if available (useful for correlating logs with OpenTelemetry traces).
//
// Trace fields added to the log if present:
//   - trace_id
//   - span_id
//
// If the context does not contain a span, the logger returned is just the global logger.
func Logger(ctx context.Context) *zerolog.Logger {
	logger := gLogger.With()

	// Extract span from context and attach trace info if available
	if span := trace.SpanFromContext(ctx); span != nil {
		sc := span.SpanContext()
		if sc.HasTraceID() {
			logger = logger.Str("trace_id", sc.TraceID().String())
		}
		if sc.HasSpanID() {
			logger = logger.Str("span_id", sc.SpanID().String())
		}
	}

	l := logger.Logger()
	return &l
}

func setupWriter(conf *config) (io.Writer, error) {
	// Setup zerolog global config
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	var consoleWriter io.Writer
	if conf.ServiceInfo.Env == DevelopmentEnv {
		// Dev mode: human-readable logs
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		consoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		// Prod mode: JSON logs
		consoleWriter = os.Stdout
	}

	// Always include console writer
	writers := []io.Writer{consoleWriter}

	// If file logging is enabled, add file writer
	if conf.EnableFileLog && conf.FilePath != "" {
		dir := filepath.Dir(conf.FilePath)
		if err := os.MkdirAll(dir, 0o744); err != nil { //nolint:gosec
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Config file rotate
		file := &lumberjack.Logger{
			Filename:   conf.FilePath,
			MaxSize:    conf.MaxAge, // megabytes
			MaxBackups: conf.MaxBackups,
			MaxAge:     conf.MaxAge, // days
			Compress:   conf.Compress,
		}

		writers = append(writers, file)
	}

	// Combine all writers
	return io.MultiWriter(writers...), nil
}

func setupAPM(conf *config) (*otelzerolog.Hook, error) {
	// configure opentelemetry logger provider
	exporterClient := otlplogshttp.NewClient(
		otlplogshttp.WithEndpoint(conf.ApmExporterEndpoint),
		otlplogshttp.WithInsecure(),
		otlplogshttp.WithHeaders(map[string]string{
			"authorization": conf.ApmApiKey,
		}),
	)

	// Default resource
	rc := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.DeploymentEnvironment(conf.ServiceInfo.Env),
		semconv.ServiceName(conf.ServiceInfo.Name),
		semconv.ServiceVersion(conf.ServiceInfo.Version),
	)

	exporter, err := otlplogs.NewExporter(context.Background(), otlplogs.WithClient(exporterClient))
	if err != nil {
		return nil, err
	}

	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(rc),
	)
	hook := otelzerolog.NewHook(loggerProvider)

	return hook, nil
}
