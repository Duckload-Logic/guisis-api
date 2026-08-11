package support

import (
	"context"
	"fmt"

	goaway "github.com/TwiN/go-away"
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
	emailer  audit.Emailer
}

func NewService(
	repo *Repository,
	notifSvc *notifications.Service,
	usersSvc *users.Service,
	emailer audit.Emailer,
) *Service {
	return &Service{
		repo:     repo,
		notifSvc: notifSvc,
		usersSvc: usersSvc,
		emailer:  emailer,
	}
}

func (s *Service) OpenTicket(
	ctx context.Context,
	req CreateTicketRequest,
	authUserID string,
) (*TicketResponse, error) {
	req.Message = goaway.Censor(req.Message)

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
		s.sendNotificationsToUserIDs(
			ctx,
			adminIDs,
			ticketID,
			"New Support Chat",
			fmt.Sprintf("New support ticket from %s", senderName),
		)
	}

	saIDs, err := s.usersSvc.GetUserIDsByRole(
		ctx, int(constants.SuperAdminRoleID),
	)
	if err == nil {
		s.sendNotificationsToUserIDs(
			ctx,
			saIDs,
			ticketID,
			"New Support Chat",
			fmt.Sprintf("New support ticket from %s", senderName),
		)
	}

	adminEmails, err1 := s.usersSvc.GetEmailsByRole(
		ctx, int(constants.AdminRoleID),
	)
	saEmails, err2 := s.usersSvc.GetEmailsByRole(
		ctx, int(constants.SuperAdminRoleID),
	)

	emailsMap := make(map[string]bool)
	if err1 == nil {
		for _, email := range adminEmails {
			if email != "" {
				emailsMap[email] = true
			}
		}
	}
	if err2 == nil {
		for _, email := range saEmails {
			if email != "" {
				emailsMap[email] = true
			}
		}
	}

	for email := range emailsMap {
		go s.sendNewTicketEmail(
			context.Background(),
			email,
			ticketID,
			senderName,
		)
	}
}

func (s *Service) sendNotificationsToUserIDs(
	ctx context.Context,
	userIDs []string,
	ticketID string,
	title string,
	message string,
) {
	for _, uid := range userIDs {
		notif := audit.NotificationEntry{
			ReceiverID: structs.StringToNullableString(uid),
			Title:      title,
			Message:    message,
			Type:       "System",
			TargetID:   structs.StringToNullableString(ticketID),
			TargetType: structs.StringToNullableString("SupportTicket"),
		}
		_ = s.notifSvc.Send(ctx, notif)
	}
}

func (s *Service) notifyUserOfNewMessage(
	ctx context.Context,
	userID string,
	ticketID string,
	senderName string,
) {
	notif := audit.NotificationEntry{
		ReceiverID: structs.StringToNullableString(userID),
		Title:      "Support Message",
		Message: fmt.Sprintf(
			"%s replied to your support ticket",
			senderName,
		),
		Type:       "System",
		TargetID:   structs.StringToNullableString(ticketID),
		TargetType: structs.StringToNullableString("SupportTicket"),
	}
	_ = s.notifSvc.Send(ctx, notif)
}

func (s *Service) AddMessage(
	ctx context.Context,
	ticketID string,
	req CreateMessageRequest,
	senderID string,
) (*MessageResponse, error) {
	req.Message = goaway.Censor(req.Message)

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
			senderName = "Admin"
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

	isStudentReply := ticket.UserID.Valid &&
		senderID == ticket.UserID.String
	isGuestReply := !ticket.UserID.Valid && senderID == ""

	if ticket.UserID.Valid &&
		ticket.UserID.String != "" &&
		senderID != ticket.UserID.String {
		go s.notifyUserOfNewMessage(
			context.Background(),
			ticket.UserID.String,
			ticketID,
			senderName,
		)
		go s.sendReplyEmailToUser(
			context.Background(),
			ticket.UserID.String,
			ticketID,
			req.Message,
		)
	} else if !ticket.UserID.Valid && senderID != "" {
		if ticket.GuestEmail.Valid && ticket.GuestEmail.String != "" {
			recipientEmail := ticket.GuestEmail.String
			recipientName := "Guest"
			if ticket.GuestName.Valid && ticket.GuestName.String != "" {
				recipientName = ticket.GuestName.String
			}
			go s.sendReplyEmail(
				context.Background(),
				recipientEmail,
				recipientName,
				ticketID,
				req.Message,
			)
		}
	} else if isStudentReply || isGuestReply {
		go s.notifyAdminsOfReply(
			context.Background(),
			ticketID,
			senderName,
		)
	}

	newStatus := ticket.Status
	if ticket.Status == "CLOSED" || ticket.Status == "RESOLVED" {
		newStatus = "OPEN"
	}

	// Update ticket status and updated_at
	_ = s.repo.UpdateTicketStatus(ctx, ticketID, newStatus)

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
	status string,
) (*ListTicketsResponse, error) {
	total, err := s.repo.GetTicketsCount(ctx, status)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to count support tickets: %w",
			err,
		)
	}

	tickets, err := s.repo.GetTickets(
		ctx,
		staffUserID,
		status,
		req.PageSize,
		req.GetOffset(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list support tickets: %w",
			err,
		)
	}

	dtos := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		dtos[i] = *s.mapQueryResultToResponse(&t, staffUserID)
	}

	return &ListTicketsResponse{
		Tickets: dtos,
		Meta: structs.CalculateMetadata(
			total, req.Page, req.PageSize,
		),
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

func (s *Service) sendReplyEmailToUser(
	ctx context.Context,
	userID string,
	ticketID string,
	message string,
) {
	user, err := s.usersSvc.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		fmt.Printf(
			"[SupportService] {EmailToUser - GetUser}: %v\n",
			err,
		)
		return
	}

	name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	s.sendReplyEmail(ctx, user.Email, name, ticketID, message)
}

func (s *Service) sendReplyEmail(
	ctx context.Context,
	email string,
	name string,
	ticketID string,
	message string,
) {
	ticketShort := ticketID
	if len(ticketShort) > 8 {
		ticketShort = ticketShort[:8]
	}

	body := fmt.Sprintf(`
<div style="font-family: sans-serif; padding: 20px; color: #333; `+
		`max-width: 600px; margin: 0 auto; border: 1px solid #eee; `+
		`border-radius: 8px;">
	<h2 style="color: #800000; border-bottom: 2px solid #800000; `+
		`padding-bottom: 10px; margin-top: 0;">GuiSIS Support</h2>
	<p>Hi %s,</p>
	<p>A support representative has replied to your ticket `+
		`(<strong>#%s</strong>):</p>
	<div style="background-color: #f9f9f9; border-left: 4px solid `+
		`#800000; padding: 12px 15px; margin: 15px 0; `+
		`font-style: italic; white-space: pre-wrap;">%s</div>
	<p>To view the full conversation or reply, please visit the `+
		`GuiSIS portal.</p>
	<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;" />
	<p style="font-size: 11px; color: #888; text-align: center; `+
		`margin-bottom: 0;">This is an automated notification. `+
		`Please do not reply directly to this email.</p>
</div>
`, name, ticketShort, message)

	emailEntry := audit.EmailEntry{
		To: []string{email},
		Subject: fmt.Sprintf(
			"GuiSIS Support - Ticket #%s Reply",
			ticketShort,
		),
		Body: body,
	}

	err := s.emailer.Send(ctx, emailEntry)
	if err != nil {
		fmt.Printf(
			"[SupportService] {SendReplyEmail - Send}: %v\n",
			err,
		)
	}
}

func (s *Service) notifyAdminsOfReply(
	ctx context.Context,
	ticketID string,
	senderName string,
) {
	title := "Support Message"
	message := fmt.Sprintf("%s replied to support ticket", senderName)

	adminIDs, err := s.usersSvc.GetUserIDsByRole(
		ctx, int(constants.AdminRoleID),
	)
	if err == nil {
		s.sendNotificationsToUserIDs(
			ctx,
			adminIDs,
			ticketID,
			title,
			message,
		)
	}

	saIDs, err := s.usersSvc.GetUserIDsByRole(
		ctx, int(constants.SuperAdminRoleID),
	)
	if err == nil {
		s.sendNotificationsToUserIDs(
			ctx,
			saIDs,
			ticketID,
			title,
			message,
		)
	}
}

func (s *Service) sendNewTicketEmail(
	ctx context.Context,
	email string,
	ticketID string,
	senderName string,
) {
	ticketShort := ticketID
	if len(ticketShort) > 8 {
		ticketShort = ticketShort[:8]
	}

	body := fmt.Sprintf(`
<div style="font-family: sans-serif; padding: 20px; color: #333; `+
		`max-width: 600px; margin: 0 auto; border: 1px solid #eee; `+
		`border-radius: 8px;">
	<h2 style="color: #800000; border-bottom: 2px solid #800000; `+
		`padding-bottom: 10px; margin-top: 0;">GuiSIS Support</h2>
	<p>Hello Admin,</p>
	<p>A new support ticket has been opened by <strong>%s</strong> `+
		`(<strong>#%s</strong>).</p>
	<p>To view and respond to this ticket, please log in to the `+
		`GuiSIS portal.</p>
	<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;" />
	<p style="font-size: 11px; color: #888; text-align: center; `+
		`margin-bottom: 0;">This is an automated notification. `+
		`Please do not reply directly to this email.</p>
</div>
`, senderName, ticketShort)

	emailEntry := audit.EmailEntry{
		To: []string{email},
		Subject: fmt.Sprintf(
			"GuiSIS Support - New Ticket #%s",
			ticketShort,
		),
		Body: body,
	}

	err := s.emailer.Send(ctx, emailEntry)
	if err != nil {
		fmt.Printf(
			"[SupportService] {SendNewTicketEmail - Send}: %v\n",
			err,
		)
	}
}
