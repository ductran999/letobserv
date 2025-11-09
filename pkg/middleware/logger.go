package middleware

import (
	"time"

	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware returns a gin middleware that logs requests with trace information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		logger.Logger(c.Request.Context()).Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Msg("Request processed")
	}
}
