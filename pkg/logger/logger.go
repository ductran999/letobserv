package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
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
//	  "service_name":"orders-service",
//	  "service_version":"v1.0.0",
//	  "service_env":"dev",
//	  "time":"2025-11-09T17:04:40+07:00",
//	  "message":"Application started"
//	}
//
// Options can override defaults, e.g., enabling file logging, rotation, or compression.
func New(serviceInfo ServiceInfo, options ...ConfigOption) error {
	// Validate basic service metadata
	if err := serviceInfo.Validate(); err != nil {
		return err
	}

	// Initialize configuration with default values
	conf := getDefaultConfig(serviceInfo)

	// Apply user-provided configuration options
	for _, opt := range options {
		if err := opt(conf); err != nil {
			return err
		}
	}

	// Prepare the log writers (e.g., console, file)
	writers, err := setupWriter(conf)
	if err != nil {
		return err
	}

	// Create the zerolog instance with timestamp and context
	gLogger = zerolog.New(writers).With().
		Timestamp().
		Str("service_name", conf.ServiceInfo.Name).
		Str("service_version", conf.ServiceInfo.Version).
		Str("service_env", conf.ServiceInfo.Env).
		Logger()
	return nil
}

// GetLoggerWithTrace returns a logger derived from the global logger gLogger,
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
		if err := os.MkdirAll(dir, 0o744); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Config file rotate
		file := &lumberjack.Logger{
			Filename:   conf.FilePath,
			MaxSize:    conf.MaxAge, // megabytes
			MaxBackups: conf.MaxBackups,
			MaxAge:     conf.MaxAge, //days
			Compress:   conf.Compress,
		}

		writers = append(writers, file)
	}

	// Combine all writers
	return io.MultiWriter(writers...), nil
}
