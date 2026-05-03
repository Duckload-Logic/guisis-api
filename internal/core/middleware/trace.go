package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		} else if len(traceID) > 64 {
			traceID = traceID[:64]
		}

		// Set in Gin context for potential use in handlers
		c.Set(string(audit.TraceIDKey), traceID)

		// Set in Request context for propagation to service layers
		ctx := context.WithValue(c.Request.Context(), audit.TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
