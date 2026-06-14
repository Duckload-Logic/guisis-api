package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/webpush"
)

type Service struct {
	repo          *Repository
	webpushClient *webpush.Client
	mu            sync.RWMutex
	subscribers   map[string][]chan audit.NotificationEntry
}

func NewService(repo *Repository, wp *webpush.Client) *Service {
	return &Service{
		repo:          repo,
		webpushClient: wp,
		subscribers:   make(map[string][]chan audit.NotificationEntry),
	}
}

func (s *Service) Send(
	ctx context.Context,
	notif audit.NotificationEntry,
) error {
	notificationID := uuid.New().String()
	notif.ID = notificationID
	notif.IsRead = false
	notif.IsTouched = false
	notif.CreatedAt = time.Now()

	err := s.repo.Create(ctx, nil, Notification{
		ID:         notificationID,
		ReceiverID: notif.ReceiverID,
		ActorID:    notif.ActorID,
		TargetID:   notif.TargetID,
		TargetType: notif.TargetType,
		Title:      notif.Title,
		Message:    notif.Message,
		Type:       notif.Type,
		IsRead:     notif.IsRead,
		IsTouched:  notif.IsTouched,
	})
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	s.Broadcast(ctx, notif)

	go s.sendWebPush(context.Background(), notif)

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
	req ListNotificationsRequest,
) (ListNotificationsResponse, error) {
	models, err := s.repo.GetByUserID(
		ctx,
		userID,
		req.UnreadOnly,
		req.PageSize,
		req.GetOffset(),
	)
	if err != nil {
		return ListNotificationsResponse{}, fmt.Errorf(
			"failed to fetch notifications for user %s: %w",
			userID,
			err,
		)
	}

	total, err := s.repo.CountByUserID(ctx, userID, req.UnreadOnly)
	if err != nil {
		return ListNotificationsResponse{}, err
	}

	unreadCount, err := s.repo.CountUnreadByUserID(ctx, userID)
	if err != nil {
		return ListNotificationsResponse{}, err
	}

	untouchedCount, err := s.repo.CountUntouchedByUserID(ctx, userID)
	if err != nil {
		return ListNotificationsResponse{}, err
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
			IsTouched:  m.IsTouched,
			CreatedAt:  m.CreatedAt,
		})
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}

	return ListNotificationsResponse{
		Notifications:  dtos,
		Total:          total,
		Page:           req.Page,
		PageSize:       req.PageSize,
		TotalPages:     totalPages,
		UnreadCount:    unreadCount,
		UntouchedCount: untouchedCount,
	}, nil
}

func (s *Service) MarkAsRead(
	ctx context.Context,
	id string,
	userID string,
) error {
	return s.repo.MarkAsRead(ctx, nil, id, userID)
}

func (s *Service) MarkAllAsTouched(
	ctx context.Context,
	userID string,
) error {
	return s.repo.MarkAllAsTouched(ctx, nil, userID)
}

func (s *Service) MarkAllAsRead(
	ctx context.Context,
	userID string,
) error {
	return s.repo.MarkAllAsRead(ctx, nil, userID)
}

func (s *Service) DeleteOldNotifications(
	ctx context.Context,
	days int,
) (int64, error) {
	return s.repo.DeleteOldNotifications(ctx, days)
}

func (s *Service) sendWebPush(
	ctx context.Context,
	notif audit.NotificationEntry,
) {
	if !notif.ReceiverID.Valid || notif.ReceiverID.String == "" {
		return
	}

	userID := notif.ReceiverID.String
	subs, err := s.repo.GetPushSubscriptionsByUserID(ctx, userID)
	if err != nil {
		fmt.Printf(
			"[notifications.Service.sendWebPush] {GetPushSubs}: %v\n",
			err,
		)
		return
	}

	if len(subs) == 0 {
		return
	}

	payload, err := json.Marshal(map[string]string{
		"title":    notif.Title,
		"message":  notif.Message,
		"type":     notif.Type,
		"targetId": notif.TargetID.String,
	})
	if err != nil {
		fmt.Printf(
			"[notifications.Service.sendWebPush] {Marshal}: %v\n",
			err,
		)
		return
	}

	for _, sub := range subs {
		go func(sub PushSubscription) {
			sendCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			status, sendErr := s.webpushClient.Send(
				sendCtx,
				sub.Endpoint,
				sub.P256dhKey,
				sub.AuthKey,
				payload,
			)
			if sendErr != nil {
				fmt.Printf(
					"[notifications.Service.sendWebPush] {Send}: %v\n",
					sendErr,
				)
				return
			}

			if status == http.StatusGone || status == http.StatusNotFound {
				delErr := s.repo.DeletePushSubscription(
					context.Background(),
					nil,
					sub.Endpoint,
					userID,
				)
				if delErr != nil {
					fmt.Printf(
						"[notifications.Service.sendWebPush] {Delete}: %v\n",
						delErr,
					)
				}
			}
		}(sub)
	}
}
