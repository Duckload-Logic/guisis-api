package notifications

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

type Service struct {
	repo        *Repository
	mu          sync.RWMutex
	subscribers map[string][]chan audit.NotificationEntry
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:        repo,
		subscribers: make(map[string][]chan audit.NotificationEntry),
	}
}

func (s *Service) Send(
	ctx context.Context,
	notif audit.NotificationEntry,
) error {
	err := s.repo.Create(ctx, nil, Notification{
		ID:         uuid.New().String(),
		ReceiverID: notif.ReceiverID,
		ActorID:    notif.ActorID,
		TargetID:   notif.TargetID,
		TargetType: notif.TargetType,
		Title:      notif.Title,
		Message:    notif.Message,
		Type:       notif.Type,
	})
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	s.Broadcast(ctx, notif)

	return nil
}

func (s *Service) Subscribe(
	ctx context.Context,
	userID string,
) (<-chan audit.NotificationEntry, func()) {
	ch := make(chan audit.NotificationEntry, 1)

	s.mu.Lock()
	s.subscribers[userID] = append(s.subscribers[userID], ch)
	s.mu.Unlock()

	return ch, func() { s.Unsubscribe(ctx, userID, ch) }
}

func (s *Service) Unsubscribe(
	ctx context.Context,
	userID string,
	ch chan audit.NotificationEntry,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, sub := range s.subscribers[userID] {
		if sub == ch {
			s.subscribers[userID] = append(
				s.subscribers[userID][:i],
				s.subscribers[userID][i+1:]...,
			)

			close(ch)
			break
		}
	}
}

func (s *Service) Broadcast(
	ctx context.Context,
	notif audit.NotificationEntry,
) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Send to all subscribers for this user
	for _, ch := range s.subscribers[notif.ReceiverID.String] {
		select {
		case ch <- notif:
		default:
			// Non-blocking send, skip if channel is full
		}
	}

	return nil
}

func (s *Service) GetUserNotifications(
	ctx context.Context,
	userID string,
) ([]audit.NotificationEntry, error) {
	models, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch notifications for user %s: %w",
			userID,
			err,
		)
	}

	dtos := make([]audit.NotificationEntry, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, audit.NotificationEntry{
			ID:         m.ID,
			ReceiverID: m.ReceiverID,
			ActorID:    m.ActorID,
			TargetID:   m.TargetID,
			TargetType: m.TargetType,
			Title:      m.Title,
			Message:    m.Message,
			Type:       m.Type,
			IsRead:     m.IsRead,
			CreatedAt:  m.CreatedAt,
		})
	}
	return dtos, nil
}

func (s *Service) MarkAsRead(
	ctx context.Context,
	id string,
	userID string,
) error {
	return s.repo.MarkAsRead(ctx, nil, id, userID)
}

func (s *Service) DeleteOldNotifications(
	ctx context.Context,
	days int,
) (int64, error) {
	return s.repo.DeleteOldNotifications(ctx, days)
}
