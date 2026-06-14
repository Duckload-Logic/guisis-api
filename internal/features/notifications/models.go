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
	IsTouched  bool                   `db:"is_touched"  json:"isTouched"`
	CreatedAt  time.Time              `db:"created_at"  json:"createdAt"`
	UpdatedAt  time.Time              `db:"updated_at"  json:"updatedAt"`
}

// PushSubscription represents a browser push notification subscription.
type PushSubscription struct {
	ID        string    `db:"id"          json:"id"`
	UserID    string    `db:"user_id"     json:"userId"`
	Endpoint  string    `db:"endpoint"    json:"endpoint"`
	P256dhKey string    `db:"p256dh_key"  json:"p256dhKey"`
	AuthKey   string    `db:"auth_key"    json:"authKey"`
	CreatedAt time.Time `db:"created_at"  json:"createdAt"`
}
