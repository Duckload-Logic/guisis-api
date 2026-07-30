package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// M2MAuditMiddleware logs every integration endpoint access to the
// system log. It must run AFTER AuthMiddleware so that m2mClientID
// and clientName are already set in context.
//
// It records:
//   - Which M2M client made the request (clientID + clientName)
//   - Which resource was accessed (method + path)
//   - The HTTP response status
//   - IP address and User-Agent
func M2MAuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		svcVal, exists := c.Get(SecurityLoggerContextKey)
		if !exists {
			log.Printf(
				"[M2MAuditMiddleware] {Context}: " +
					"logger not found in context",
			)
			return
		}

		svc, ok := svcVal.(SecurityLogger)
		if !ok || svc == nil {
			log.Printf(
				"[M2MAuditMiddleware] {Context}: " +
					"logger type assertion failed",
			)
			return
		}

		clientID := ginStringVal(c, "m2mClientID")
		clientName := ginStringVal(c, "clientName")
		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		ua := c.Request.UserAgent()

		action := audit.ActionM2MDataAccess
		level := audit.LevelInfo
		if status >= http.StatusForbidden {
			action = audit.ActionM2MDataAccessDenied
			level = audit.LevelWarning
		}

		msg := fmt.Sprintf(
			"M2M client '%s' (%s) %s %s -> %d",
			clientName,
			clientID,
			method,
			path,
			status,
		)

		svc.RecordEntry(
			c.Request.Context(),
			audit.LogEntry{
				Level:    level,
				Category: audit.CategoryM2M,
				Action:   action,
				Message:  msg,
				UserID: structs.StringToNullableString(
					clientID,
				),
				UserEmail: structs.StringToNullableString(
					clientName,
				),
				IPAddress: structs.StringToNullableString(ip),
				UserAgent: structs.StringToNullableString(ua),
				Metadata: &audit.LogMetadata{
					EntityType: "m2m_client",
					EntityID:   clientID,
				},
			},
		)
	}
}

func ginStringVal(c *gin.Context, key string) string {
	val, _ := c.Get(key)
	s, _ := val.(string)
	return s
}
