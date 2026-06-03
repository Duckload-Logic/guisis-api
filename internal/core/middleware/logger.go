package middleware

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

// SecurityLogger is an interface for recording security log events.
// This interface is implemented by logs.Service and is used to break
// the import cycle between middleware and logs packages.
type SecurityLogger interface {
	RecordSecurity(
		ctx context.Context,
		action, message, userEmail, ipAddress, userAgent string,
	)
	RecordEntry(ctx context.Context, entry audit.LogEntry)
}

// SecurityLoggerContextKey is the gin context key used to store the security
// logger.
const SecurityLoggerContextKey = "securityLogger"
