package support

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type SupportTicket struct {
	ID        string                 `db:"id"          json:"id"`
	UserID    structs.NullableString `db:"user_id"     json:"userId"`
	GuestName structs.NullableString `db:"guest_name"  json:"guestName"`
	GuestEmail structs.NullableString `db:"guest_email" json:"guestEmail"`
	Status    string                 `db:"status"      json:"status"`
	CreatedAt time.Time              `db:"created_at"  json:"createdAt"`
	UpdatedAt time.Time              `db:"updated_at"  json:"updatedAt"`
}

type SupportMessage struct {
	ID         string                 `db:"id"          json:"id"`
	TicketID   string                 `db:"ticket_id"   json:"ticketId"`
	SenderID   structs.NullableString `db:"sender_id"   json:"senderId"`
	SenderName string                 `db:"sender_name" json:"senderName"`
	Message    string                 `db:"message"     json:"message"`
	CreatedAt  time.Time              `db:"created_at"  json:"createdAt"`
}
