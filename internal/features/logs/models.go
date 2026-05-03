package logs

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// SystemLog represents a system log entry, used for both business logic and data persistence.
type SystemLog struct {
	ID          int                    `db:"id"           json:"id"`
	Level       string                 `db:"level"        json:"level"`
	Category    string                 `db:"category"     json:"category"`
	Action      string                 `db:"action"       json:"action"`
	Message     string                 `db:"message"      json:"message"`
	UserID      structs.NullableString `db:"user_id"      json:"userId"`
	TargetID    structs.NullableString `db:"target_id"    json:"targetId"`
	TargetType  structs.NullableString `db:"target_type"  json:"targetType"`
	UserEmail   structs.NullableString `db:"user_email"   json:"userEmail"`
	TargetEmail structs.NullableString `db:"target_email" json:"targetEmail"`
	IPAddress   structs.NullableString `db:"ip_address"   json:"ipAddress"`
	UserAgent   structs.NullableString `db:"user_agent"   json:"userAgent"`
	Metadata    structs.NullableString `db:"metadata"     json:"metadata"`
	TraceID     structs.NullableString `db:"trace_id"     json:"traceId"`
	CreatedAt   time.Time              `db:"created_at"   json:"createdAt"`
}
