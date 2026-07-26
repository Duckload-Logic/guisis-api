package support

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type Service struct {
	repo     *Repository
	notifSvc *notifications.Service
	usersSvc *users.Service
}

func NewService(
	repo *Repository,
	notifSvc *notifications.Service,
	usersSvc *users.Service,
) *Service {
	return &Service{
		repo:     repo,
		notifSvc: notifSvc,
		usersSvc: usersSvc,
	}
}

func (s *Service) OpenTicket(
	ctx context.Context,
	req CreateTicketRequest,
	authUserID string,
) (*TicketResponse, error) {
	ticketID := uuid.New().String()
	ticket := &SupportTicket{
		ID:     ticketID,
		Status: "OPEN",
	}

	var senderName string
	if authUserID != "" {
		ticket.UserID = structs.StringToNullableString(authUserID)
		user, err := s.usersSvc.GetUserByID(ctx, authUserID)
		if err == nil && user != nil {
			senderName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		} else {
			senderName = "Student"
		}
	} else {
		if req.GuestName != nil {
			ticket.GuestName = structs.PointerToNullableString(req.GuestName)
			senderName = *req.GuestName
		} else {
			senderName = "Guest"
		}
		if req.GuestEmail != nil {
			ticket.GuestEmail = structs.PointerToNullableString(req.GuestEmail)
		}
	}

	err := s.repo.CreateTicket(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to open ticket: %w", err)
	}

	// Create first message
	messageID := uuid.New().String()
	firstMsg := &SupportMessage{
		ID:         messageID,
		TicketID:   ticketID,
		SenderName: senderName,
		Message:    req.Message,
	}
	if authUserID != "" {
		firstMsg.SenderID = structs.StringToNullableString(authUserID)
	}

	err = s.repo.CreateMessage(ctx, firstMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial message: %w", err)
	}

	// Notify admins/superadmins about the new ticket
	go s.notifyAdminsOfNewTicket(context.Background(), ticketID, senderName)

	return s.mapTicketToResponse(ticket), nil
}

func (s *Service) notifyAdminsOfNewTicket(
	ctx context.Context,
	ticketID string,
	senderName string,
) {
	adminIDs, err := s.usersSvc.GetUserIDsByRole(
		ctx, int(constants.AdminRoleID),
	)
	if err == nil {
		s.sendNotificationsToUserIDs(ctx, adminIDs, ticketID, senderName)
	}

	saIDs, err := s.usersSvc.GetUserIDsByRole(
		ctx, int(constants.SuperAdminRoleID),
	)
	if err == nil {
		s.sendNotificationsToUserIDs(ctx, saIDs, ticketID, senderName)
	}
}

func (s *Service) sendNotificationsToUserIDs(
	ctx context.Context,
	userIDs []string,
	ticketID string,
	senderName string,
) {
	for _, uid := range userIDs {
		notif := audit.NotificationEntry{
			ReceiverID: structs.StringToNullableString(uid),
			Title:      "New Support Chat",
			Message:    fmt.Sprintf("New support ticket from %s", senderName),
			Type:       "System",
			TargetID:   structs.StringToNullableString(ticketID),
			TargetType: structs.StringToNullableString("SupportTicket"),
		}
		_ = s.notifSvc.Send(ctx, notif)
	}
}

func (s *Service) AddMessage(
	ctx context.Context,
	ticketID string,
	req CreateMessageRequest,
	senderID string,
) (*MessageResponse, error) {
	ticket, err := s.repo.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	var senderName string
	if senderID != "" {
		user, err := s.usersSvc.GetUserByID(ctx, senderID)
		if err == nil && user != nil {
			senderName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		} else {
			senderName = "Staff"
		}
	} else {
		if ticket.GuestName.Valid && ticket.GuestName.String != "" {
			senderName = ticket.GuestName.String
		} else {
			senderName = "Guest"
		}
	}

	msgID := uuid.New().String()
	msg := &SupportMessage{
		ID:         msgID,
		TicketID:   ticketID,
		SenderName: senderName,
		Message:    req.Message,
	}
	if senderID != "" {
		msg.SenderID = structs.StringToNullableString(senderID)
	}

	err = s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update ticket updated_at
	_ = s.repo.UpdateTicketStatus(ctx, ticketID, ticket.Status)

	return s.mapMessageToResponse(msg), nil
}

func (s *Service) GetTicketMessages(
	ctx context.Context,
	ticketID string,
) ([]MessageResponse, error) {
	messages, err := s.repo.GetMessagesByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	dtos := make([]MessageResponse, len(messages))
	for i, m := range messages {
		dtos[i] = *s.mapMessageToResponse(&m)
	}
	return dtos, nil
}

func (s *Service) GetTickets(
	ctx context.Context,
	staffUserID string,
	req structs.PaginationRequest,
) (*ListTicketsResponse, error) {
	total, err := s.repo.GetTicketsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count support tickets: %w", err)
	}

	tickets, err := s.repo.GetTickets(
		ctx,
		staffUserID,
		req.PageSize,
		req.GetOffset(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list support tickets: %w", err)
	}

	dtos := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		dtos[i] = *s.mapQueryResultToResponse(&t, staffUserID)
	}

	return &ListTicketsResponse{
		Tickets: dtos,
		Meta:    structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetTicketsByUserID(
	ctx context.Context,
	userID string,
) ([]TicketResponse, error) {
	tickets, err := s.repo.GetTicketsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		dtos[i] = *s.mapTicketToResponse(&t)
	}
	return dtos, nil
}

func (s *Service) CloseTicket(
	ctx context.Context,
	ticketID string,
) error {
	return s.repo.UpdateTicketStatus(ctx, ticketID, "CLOSED")
}

func (s *Service) GetTicket(
	ctx context.Context,
	id string,
) (*SupportTicket, error) {
	return s.repo.GetTicket(ctx, id)
}

func (s *Service) mapTicketToResponse(t *SupportTicket) *TicketResponse {
	var userID *string
	if t.UserID.Valid && t.UserID.String != "" {
		userID = &t.UserID.String
	}
	var guestName *string
	if t.GuestName.Valid && t.GuestName.String != "" {
		guestName = &t.GuestName.String
	}
	var guestEmail *string
	if t.GuestEmail.Valid && t.GuestEmail.String != "" {
		guestEmail = &t.GuestEmail.String
	}

	return &TicketResponse{
		ID:         t.ID,
		UserID:     userID,
		GuestName:  guestName,
		GuestEmail: guestEmail,
		Status:     t.Status,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func (s *Service) mapQueryResultToResponse(
	q *TicketQueryResult,
	staffUserID string,
) *TicketResponse {
	var userID *string
	if q.UserID.Valid && q.UserID.String != "" {
		userID = &q.UserID.String
	}
	var guestName *string
	if q.GuestName.Valid && q.GuestName.String != "" {
		guestName = &q.GuestName.String
	}
	var guestEmail *string
	if q.GuestEmail.Valid && q.GuestEmail.String != "" {
		guestEmail = &q.GuestEmail.String
	}

	var studentName *string
	var studentEmail *string
	if userID != nil {
		fullName := fmt.Sprintf(
			"%s %s",
			q.UserFirstName.String,
			q.UserLastName.String,
		)
		studentName = &fullName
		studentEmail = &q.UserEmail.String
	}

	var lastMsg *string
	if q.LastMessage.Valid && q.LastMessage.String != "" {
		var formattedMsg string
		if q.LastSenderID.Valid && q.LastSenderID.String == staffUserID {
			formattedMsg = fmt.Sprintf("You: %s", q.LastMessage.String)
		} else {
			formattedMsg = fmt.Sprintf(
				"%s: %s",
				q.LastSenderName.String,
				q.LastMessage.String,
			)
		}
		lastMsg = &formattedMsg
	}

	var profilePic *string
	if q.UserProfilePicture.Valid && q.UserProfilePicture.String != "" {
		profilePic = &q.UserProfilePicture.String
	}

	return &TicketResponse{
		ID:             q.ID,
		UserID:         userID,
		GuestName:      guestName,
		GuestEmail:     guestEmail,
		Status:         q.Status,
		CreatedAt:      q.CreatedAt,
		UpdatedAt:      q.UpdatedAt,
		StudentName:    studentName,
		StudentEmail:   studentEmail,
		ProfilePicture: profilePic,
		LastMessage:    lastMsg,
		IsRead:         q.IsRead,
	}
}

func (s *Service) MarkTicketAsRead(
	ctx context.Context,
	ticketID string,
	userID string,
) error {
	return s.repo.MarkTicketAsRead(ctx, ticketID, userID)
}

func (s *Service) mapMessageToResponse(m *SupportMessage) *MessageResponse {
	var senderID *string
	if m.SenderID.Valid && m.SenderID.String != "" {
		senderID = &m.SenderID.String
	}

	return &MessageResponse{
		ID:         m.ID,
		TicketID:   m.TicketID,
		SenderID:   senderID,
		SenderName: m.SenderName,
		Message:    m.Message,
		CreatedAt:  m.CreatedAt,
	}
}
