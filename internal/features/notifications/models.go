package notifications

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Notification represents a system notification, used for both business logic and data persistence.
type Notification struct {
	ID         string                 `db:"id"          json:"id"`
	ReceiverID structs.NullableString `db:"receiver_id" json:"receiverId"`
	ActorID    structs.NullableString `db:"actor_id"    json:"actorId"`
	TargetID   structs.NullableString `db:"target_id"   json:"targetId"`
	TargetType structs.NullableString `db:"target_type" json:"targetType"`
	Title      string                 `db:"title"       json:"title"`
	Message    string                 `db:"message"     json:"message"`
	Type       string                 `db:"type"        json:"type"`
	IsRead     bool                   `db:"is_read"     json:"isRead"`
	CreatedAt  time.Time              `db:"created_at"  json:"createdAt"`
	UpdatedAt  time.Time              `db:"updated_at"  json:"updatedAt"`
}
