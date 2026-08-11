package support

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type ListTicketsResponse struct {
	Tickets []TicketResponse           `json:"tickets"`
	Meta    structs.PaginationMetadata `json:"meta"`
}

type GetTicketsRequest struct {
	structs.PaginationRequest
	Status string `json:"status" form:"status"`
}

type CreateTicketRequest struct {
	GuestName  *string `json:"guestName"`
	GuestEmail *string `json:"guestEmail"`
	Message    string  `json:"message" binding:"required"`
}

type CreateMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type TicketResponse struct {
	ID           string    `json:"id"`
	UserID       *string   `json:"userId"`
	GuestName    *string   `json:"guestName"`
	GuestEmail   *string   `json:"guestEmail"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	StudentName    *string   `json:"studentName,omitempty"`
	StudentEmail   *string   `json:"studentEmail,omitempty"`
	ProfilePicture *string   `json:"profilePicture,omitempty"`
	LastMessage    *string   `json:"lastMessage,omitempty"`
	IsRead         bool      `json:"isRead"`
}

type MessageResponse struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticketId"`
	SenderID   *string   `json:"senderId"`
	SenderName string    `json:"senderName"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
}
