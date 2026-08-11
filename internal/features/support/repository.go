package support

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTicket(
	ctx context.Context,
	ticket *SupportTicket,
) error {
	query := `
		INSERT INTO support_tickets (
			id, user_id, guest_name, guest_email, status
		) VALUES (
			:id, :user_id, :guest_name, :guest_email, :status
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		return fmt.Errorf("failed to create support ticket: %w", err)
	}
	return nil
}

func (r *Repository) CreateMessage(
	ctx context.Context,
	msg *SupportMessage,
) error {
	query := `
		INSERT INTO support_messages (
			id, ticket_id, sender_id, sender_name, message
		) VALUES (
			:id, :ticket_id, :sender_id, :sender_name, :message
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, msg)
	if err != nil {
		return fmt.Errorf("failed to create support message: %w", err)
	}
	return nil
}

func (r *Repository) GetTicket(
	ctx context.Context,
	id string,
) (*SupportTicket, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM support_tickets WHERE id = ?",
		datastore.GetColumns(SupportTicket{}),
	)
	var ticket SupportTicket
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get support ticket: %w", err)
	}
	return &ticket, nil
}

type TicketQueryResult struct {
	ID                 string                 `db:"id"`
	UserID             structs.NullableString `db:"user_id"`
	GuestName          structs.NullableString `db:"guest_name"`
	GuestEmail         structs.NullableString `db:"guest_email"`
	Status             string                 `db:"status"`
	CreatedAt          time.Time              `db:"created_at"`
	UpdatedAt          time.Time              `db:"updated_at"`
	UserFirstName      structs.NullableString `db:"user_first_name"`
	UserLastName       structs.NullableString `db:"user_last_name"`
	UserEmail          structs.NullableString `db:"user_email"`
	UserProfilePicture structs.NullableString `db:"user_profile_picture"`
	LastSenderID       structs.NullableString `db:"last_sender_id"`
	LastSenderName     structs.NullableString `db:"last_sender_name"`
	LastMessage        structs.NullableString `db:"last_message"`
	IsRead             bool                   `db:"is_read"`
}

func (r *Repository) GetTickets(
	ctx context.Context,
	staffUserID string,
	status string,
	limit int,
	offset int,
) ([]TicketQueryResult, error) {
	statusFilter := ""
	var args []interface{}
	args = append(args, staffUserID)

	if status != "" {
		statusFilter = "WHERE t.status = ?"
		args = append(args, status)
	}
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT 
			t.id, t.user_id, t.guest_name, t.guest_email, t.status, 
			t.created_at, t.updated_at,
			u.first_name AS user_first_name, u.last_name AS user_last_name, 
			u.email AS user_email,
			f.file_url AS user_profile_picture,
			m.sender_id AS last_sender_id, m.sender_name AS last_sender_name, 
			m.message AS last_message,
			COALESCE(tr.read_at >= t.updated_at, FALSE) AS is_read
		FROM support_tickets t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN profile_pictures pp ON pp.user_id = u.id
		LEFT JOIN files f ON f.id = pp.file_id
		LEFT JOIN support_ticket_reads tr 
			ON t.id = tr.ticket_id AND tr.user_id = ?
		LEFT JOIN support_messages m ON m.id = (
			SELECT id FROM support_messages 
			WHERE ticket_id = t.id 
			ORDER BY created_at DESC 
			LIMIT 1
		)
		%s
		ORDER BY t.updated_at DESC
		LIMIT ? OFFSET ?
	`, statusFilter)

	var tickets []TicketQueryResult
	err := r.db.SelectContext(
		ctx,
		&tickets,
		query,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get support tickets: %w", err)
	}
	return tickets, nil
}

func (r *Repository) GetTicketsCount(
	ctx context.Context,
	status string,
) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM support_tickets"
	var args []interface{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count support tickets: %w", err)
	}
	return count, nil
}

func (r *Repository) MarkTicketAsRead(
	ctx context.Context,
	ticketID string,
	userID string,
) error {
	query := `
		INSERT INTO support_ticket_reads (ticket_id, user_id, read_at)
		VALUES (?, ?, NOW())
		ON DUPLICATE KEY UPDATE read_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, ticketID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark ticket as read: %w", err)
	}
	return nil
}

func (r *Repository) GetTicketsByUserID(
	ctx context.Context,
	userID string,
) ([]SupportTicket, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM support_tickets WHERE user_id = ? "+
			"ORDER BY updated_at DESC",
		datastore.GetColumns(SupportTicket{}),
	)
	var tickets []SupportTicket
	err := r.db.SelectContext(ctx, &tickets, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get support tickets by user id: %w",
			err,
		)
	}
	return tickets, nil
}

func (r *Repository) GetMessagesByTicketID(
	ctx context.Context,
	ticketID string,
) ([]SupportMessage, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM support_messages "+
			"WHERE ticket_id = ? ORDER BY created_at ASC",
		datastore.GetColumns(SupportMessage{}),
	)
	var messages []SupportMessage
	err := r.db.SelectContext(ctx, &messages, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get support messages: %w", err)
	}
	return messages, nil
}

func (r *Repository) UpdateTicketStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	query := `
		UPDATE support_tickets
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}
	return nil
}
