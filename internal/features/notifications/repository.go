package notifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) GetByUserID(
	ctx context.Context,
	userID string,
	unreadOnly bool,
	limit int,
	offset int,
) ([]Notification, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT id, receiver_id, actor_id, target_id, target_type,
			title, message, type, is_read, is_touched,
			created_at, updated_at
		FROM notifications
		WHERE receiver_id = ?
	`)

	args := []interface{}{userID}
	if unreadOnly {
		query.WriteString(" AND is_read = FALSE")
	}

	query.WriteString(" ORDER BY created_at DESC LIMIT ? OFFSET ?")
	args = append(args, limit, offset)

	var results []Notification
	err := r.db.SelectContext(ctx, &results, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get notifications for user %s: %w",
			userID,
			err,
		)
	}

	return results, nil
}

func (r *Repository) CountByUserID(
	ctx context.Context,
	userID string,
	unreadOnly bool,
) (int, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT COUNT(*)
		FROM notifications
		WHERE receiver_id = ?
	`)

	args := []interface{}{userID}
	if unreadOnly {
		query.WriteString(" AND is_read = FALSE")
	}

	var total int
	err := r.db.GetContext(ctx, &total, query.String(), args...)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to count notifications for user %s: %w",
			userID,
			err,
		)
	}

	return total, nil
}

func (r *Repository) CountUnreadByUserID(
	ctx context.Context,
	userID string,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE receiver_id = ? AND is_read = FALSE
	`

	var total int
	err := r.db.GetContext(ctx, &total, query, userID)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to count unread notifications for user %s: %w",
			userID,
			err,
		)
	}

	return total, nil
}

func (r *Repository) CountUntouchedByUserID(
	ctx context.Context,
	userID string,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE receiver_id = ? AND is_touched = FALSE
	`

	var total int
	err := r.db.GetContext(ctx, &total, query, userID)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to count untouched notifications for user %s: %w",
			userID,
			err,
		)
	}

	return total, nil
}

func (r *Repository) MarkAsRead(
	ctx context.Context,
	tx datastore.DB,
	id string,
	userID string,
) error {
	if tx == nil {
		tx = r.db
	}

	query := `
		UPDATE notifications
		SET is_read = TRUE, is_touched = TRUE
		WHERE id = ? AND receiver_id = ?
	`
	res, err := tx.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf(
			"failed to mark notification %s as read: %w",
			id,
			err,
		)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("notification not found or unauthorized")
	}

	return nil
}

func (r *Repository) MarkAllAsTouched(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	if tx == nil {
		tx = r.db
	}

	query := `
		UPDATE notifications
		SET is_touched = TRUE
		WHERE receiver_id = ? AND is_touched = FALSE
	`
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf(
			"failed to mark notifications as touched for user %s: %w",
			userID,
			err,
		)
	}

	return nil
}

func (r *Repository) MarkAllAsRead(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	if tx == nil {
		tx = r.db
	}

	query := `
		UPDATE notifications
		SET is_read = TRUE, is_touched = TRUE
		WHERE receiver_id = ? AND (is_read = FALSE OR is_touched = FALSE)
	`
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf(
			"failed to mark notifications as read for user %s: %w",
			userID,
			err,
		)
	}

	return nil
}

func (r *Repository) Create(
	ctx context.Context,
	tx datastore.DB,
	notif Notification,
) error {
	query := `
        INSERT INTO notifications (
			id, receiver_id, actor_id, target_id, target_type,
			title, message, type, is_read, is_touched, updated_at
		) VALUES (
			:id, :receiver_id, :actor_id, :target_id, :target_type,
			:title, :message, :type, :is_read, :is_touched, NOW()
		)`

	if tx == nil {
		tx = r.db
	}

	_, err := tx.NamedExecContext(ctx, query, &notif)
	if err != nil {
		return fmt.Errorf(
			"failed to create notification for %s: %w",
			notif.ReceiverID.String,
			err,
		)
	}
	return nil
}

func (r *Repository) DeleteOldNotifications(
	ctx context.Context,
	days int,
) (int64, error) {
	query := `
		DELETE FROM notifications
		WHERE created_at < DATE_SUB(NOW(), INTERVAL ? DAY)
		AND is_read = TRUE
	`
	res, err := r.db.ExecContext(ctx, query, days)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old notifications: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

func (r *Repository) SavePushSubscription(
	ctx context.Context,
	tx datastore.DB,
	sub PushSubscription,
) error {
	if tx == nil {
		tx = r.db
	}

	query := `
		INSERT INTO push_subscriptions (
			id, user_id, endpoint, p256dh_key, auth_key, created_at
		) VALUES (
			:id, :user_id, :endpoint, :p256dh_key, :auth_key, NOW()
		) ON DUPLICATE KEY UPDATE
			user_id = VALUES(user_id),
			p256dh_key = VALUES(p256dh_key),
			auth_key = VALUES(auth_key)
	`

	_, err := tx.NamedExecContext(ctx, query, &sub)
	if err != nil {
		return fmt.Errorf("failed to save push sub: %w", err)
	}
	return nil
}

func (r *Repository) GetPushSubscriptionsByUserID(
	ctx context.Context,
	userID string,
) ([]PushSubscription, error) {
	query := `
		SELECT id, user_id, endpoint, p256dh_key, auth_key, created_at
		FROM push_subscriptions
		WHERE user_id = ?
	`

	var results []PushSubscription
	err := r.db.SelectContext(ctx, &results, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get push subs for user %s: %w",
			userID,
			err,
		)
	}

	return results, nil
}

func (r *Repository) DeletePushSubscription(
	ctx context.Context,
	tx datastore.DB,
	endpoint string,
	userID string,
) error {
	if tx == nil {
		tx = r.db
	}

	query := `
		DELETE FROM push_subscriptions
		WHERE endpoint = ? AND user_id = ?
	`
	_, err := tx.ExecContext(ctx, query, endpoint, userID)
	if err != nil {
		return fmt.Errorf("failed to delete push sub: %w", err)
	}
	return nil
}
