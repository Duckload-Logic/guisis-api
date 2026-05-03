package audit

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// Logger defines the interface for recording system logs.
type Logger interface {
	Record(ctx context.Context, tx datastore.DB, entry LogEntry)
}

// LogReader defines the interface for reading system logs.
type LogReader interface {
	ListLogs(
		ctx context.Context,
		req ListSystemLogsRequest,
	) (*ListSystemLogsDTO, error)
	GetStats(
		ctx context.Context,
		startDate, endDate string,
	) ([]LogStatsDTO, error)
	GetActivityStats(ctx context.Context) ([]LogActivityDTO, error)
}

// Notifier defines the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, notif NotificationEntry) error
}

// UserGetter defines the interface for fetching user IDs by role.
type UserGetter interface {
	GetUserIDsByRole(ctx context.Context, roleID int) ([]string, error)
}
