package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

// AuditContextMiddleware enriches the request context with audit metadata
// (user email, IP address, User-Agent). Must be placed after AuthMiddleware
// so that userEmail is available.
func AuditContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := ""
		if id, exists := c.Get("userID"); exists {
			if e, ok := id.(string); ok {
				userID = e
			}
		}

		userEmail, exists := c.Get("userEmail")
		if !exists {
			userEmail = ""
		} else if e, ok := userEmail.(string); ok {
			userEmail = e
		}

		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// Extract traceID set by TraceMiddleware
		traceID := ""
		if val, exists := c.Get(string(audit.TraceIDKey)); exists {
			if t, ok := val.(string); ok {
				traceID = t
			}
		}

		roles := ""
		if val, exists := c.Get("roleIDs"); exists {
			if rids, ok := val.([]int); ok {
				roles = fmt.Sprintf("%v", rids)
			}
		}

		ctx := audit.WithContext(
			c.Request.Context(),
			ipAddress,
			userAgent,
			userID,
			userEmail.(string),
			roles,
			traceID,
		)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
