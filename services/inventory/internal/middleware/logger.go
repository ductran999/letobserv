package middleware

import (
	"time"

	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/xcontext"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// LoggingMiddleware returns a gin middleware that logs requests with trace information.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		traceID := uuid.New().String()
		span := trace.SpanFromContext(c.Request.Context())
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
			c.Set(xcontext.TraceIDKey, traceID)
		}

		// Process request
		c.Next()

		// Log request details
		logger.Logger(c.Request.Context()).Info().
			Str("method", c.Request.Method).
			Str("trace_id", traceID).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Msg("Request processed")
	}
}
