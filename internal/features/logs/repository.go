package logs

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) Record(
	ctx context.Context,
	tx datastore.DB,
	log *SystemLog,
) error {
	query := `
		INSERT INTO system_logs (
			level, category, action, message, user_id, target_id, target_type,
			user_email, target_email, ip_address, user_agent, metadata, trace_id
		) VALUES (
			:level, :category, :action, :message, :user_id, :target_id, :target_type,
			:user_email, :target_email, :ip_address, :user_agent, :metadata, :trace_id
		)
	`

	exec := tx
	if exec == nil {
		exec = r.db
	}

	_, err := exec.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("failed to insert system log: %w", err)
	}

	return nil
}

func (r *Repository) List(
	ctx context.Context, offset, limit int,
	category, action, userEmail, targetType, targetEmail,
	search, startDate, endDate, orderBy string,
) ([]SystemLog, error) {
	query, args := r.applyLogFilters(
		"SELECT id, level, category, action, message, user_id, target_id, target_type, user_email, target_email, ip_address, user_agent, metadata, trace_id, created_at FROM system_logs WHERE 1=1",
		nil,
		category,
		action,
		userEmail,
		targetType,
		targetEmail,
		search,
		startDate,
		endDate,
	)

	if orderBy == "" {
		orderBy = "created_at"
	}

	query += fmt.Sprintf(" ORDER BY %s DESC LIMIT ? OFFSET ?", orderBy)
	args = append(args, limit, offset)

	var logs []SystemLog
	err := r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list system logs: %w", err)
	}

	return logs, nil
}

// GetTotalCount returns the total count of system log entries matching filters
func (r *Repository) GetTotalCount(
	ctx context.Context,
	category, action, userEmail, targetType, targetEmail,
	search, startDate, endDate string,
) (int, error) {
	query, args := r.applyLogFilters(
		"SELECT COUNT(*) FROM system_logs WHERE 1=1",
		nil,
		category, action, userEmail, targetType, targetEmail,
		search, startDate, endDate,
	)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count system logs: %w", err)
	}

	return count, nil
}

func (r *Repository) applyLogFilters(
	query string,
	args []interface{},
	category, action, userEmail, targetType, targetEmail,
	search, startDate, endDate string,
) (string, []interface{}) {
	if args == nil {
		args = []interface{}{}
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if action != "" {
		query += " AND action = ?"
		args = append(args, action)
	}

	if userEmail != "" {
		query += " AND (user_email = ? OR target_email = ?)"
		args = append(args, userEmail, userEmail)
	}

	if targetType != "" {
		query += " AND target_type = ?"
		args = append(args, targetType)
	}

	if targetEmail != "" {
		query += " AND target_email = ?"
		args = append(args, targetEmail)
	}

	if search != "" {
		query += " AND (message LIKE ? OR action LIKE ? OR user_email LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if startDate != "" {
		query += " AND created_at >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND created_at <= ?"
		args = append(args, endDate+" 23:59:59")
	}

	return query, args
}

// GetStats returns log counts grouped by category
func (r *Repository) GetStats(
	ctx context.Context,
	startDate, endDate string,
) ([]audit.LogStatsDTO, error) {
	query, args := r.applyLogFilters(
		"SELECT category, COUNT(*) as count FROM system_logs WHERE 1=1",
		nil,
		"", "", "", "", "", "", startDate, endDate,
	)

	query += " GROUP BY category ORDER BY category"

	var stats []audit.LogStatsDTO
	err := r.db.SelectContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get log stats: %w", err)
	}

	return stats, nil
}

// GetActivityStats returns log counts grouped by hour for the last 24 hours
func (r *Repository) GetActivityStats(
	ctx context.Context,
) ([]audit.LogActivityDTO, error) {
	query := `
		SELECT
			DATE_FORMAT(created_at, '%Y-%m-%d %H:00') as time,
			COUNT(CASE WHEN level != 'ERROR' THEN 1 END) as requests,
			COUNT(CASE WHEN level = 'ERROR' THEN 1 END) as errors
		FROM system_logs
		WHERE created_at >= NOW() - INTERVAL 24 HOUR
		GROUP BY time
		ORDER BY time ASC
	`

	var stats []audit.LogActivityDTO
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get log activity stats: %w", err)
	}

	return stats, nil
}

func (r *Repository) DeleteLogsOlderThan(
	ctx context.Context,
	days int,
) (int64, error) {
	query := `
		DELETE FROM system_logs
		WHERE created_at < DATE_SUB(NOW(), INTERVAL ? DAY)
	`
	res, err := r.db.ExecContext(ctx, query, days)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old logs: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

