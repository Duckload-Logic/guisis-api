package appointments

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

const (
	filterWhenDateGe = " AND a.when_date >= ?"
	filterWhenDateLe = " AND a.when_date <= ?"
)

const appointmentsBaseQuery = `
	SELECT
		a.id,
		COALESCE(u.id, '') AS user_id,
		COALESCE(ir.id, '') AS iir_id,
		COALESCE(spi.student_number, '') AS student_number,
		COALESCE(u.first_name, '') AS user_first_name,
		u.middle_name AS user_middle_name,
		COALESCE(u.last_name, '') AS user_last_name,
		COALESCE(u.email, '') AS user_email,
		a.reason AS reason,
		a.admin_notes AS admin_notes,
		DATE_FORMAT(a.when_date, '%Y-%m-%d') AS when_date,
		a.created_at AS created_at,
		a.updated_at AS updated_at,
		ts.id AS time_slot_id,
		ts.time AS time_slot_time,
		ac.id AS category_id,
		ac.name AS category_name,
		as2.id AS status_id,
		as2.name AS status_name,
		a.urgency_level AS urgency_level,
		a.urgency_score AS urgency_score
	FROM appointments a
	LEFT JOIN iir_records ir ON a.iir_id = ir.id
	LEFT JOIN users u ON ir.user_id = u.id
	LEFT JOIN student_personal_info spi ON ir.id = spi.iir_id
	JOIN time_slots ts ON a.time_slot_id = ts.id
	JOIN appointment_categories ac ON
		a.appointment_category_id = ac.id
	JOIN statuses as2 ON a.status_id = as2.id
`

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}


func (r *Repository) GetCategories(
	ctx context.Context,
) ([]AppointmentCategory, error) {
	query := `SELECT id, name FROM appointment_categories`

	var categories []AppointmentCategory
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get concern categories: %w", err)
	}

	return categories, nil
}

func (r *Repository) GetAppointment(
	ctx context.Context,
	db datastore.DB,
	id string,
) (*AppointmentWithDetailsView, error) {
	query := fmt.Sprintf(`
		%s
		WHERE a.id = ?
	`, appointmentsBaseQuery)

	var appt AppointmentWithDetailsView
	err := db.GetContext(ctx, &appt, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return &appt, nil
}

func (r *Repository) GetUserIDByAppointmentID(
	ctx context.Context,
	id string,
) (string, error) {
	var userID string
	query := `
		SELECT ir.user_id
		FROM appointments a
		JOIN iir_records ir ON a.iir_id = ir.id
		WHERE a.id = ?
	`
	err := r.db.GetContext(ctx, &userID, query, id)
	return userID, err
}

func (r *Repository) GetDailyStatusCount(
	ctx context.Context,
	startDate, endDate string,
) ([]DailyStatusCount, error) {
	query := `
		SELECT
			DATE(a.when_date) as date,
			COUNT(CASE WHEN s.name = 'Pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN s.name = 'Scheduled' THEN 1 END) as scheduled_count,
			COUNT(CASE WHEN s.name = 'Rescheduled' THEN 1 END) as rescheduled_count
		FROM appointments a
		JOIN statuses s ON a.status_id = s.id
		WHERE when_date BETWEEN ? AND ?
		GROUP BY DATE(a.when_date);
	`

	var dsc []DailyStatusCount
	err := r.db.SelectContext(ctx, &dsc, query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return dsc, nil
}

func (r *Repository) GetTotalAppointmentsCount(
	ctx context.Context,
	statusID, startDate, endDate string,
	iirID *string,
) (int, error) {
	query, args := r.applyFilters(
		"SELECT COUNT(*) FROM appointments a WHERE 1=1",
		nil,
		statusID,
		startDate,
		endDate,
		iirID,
	)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get appointment count: %w",
			err,
		)
	}
	return count, nil
}

func (r *Repository) applyFilters(
	query string,
	args []interface{},
	statusID, startDate, endDate string,
	iirID *string,
) (string, []interface{}) {
	if args == nil {
		args = []interface{}{}
	}

	if statusID != "" {
		query += " AND a.status_id = ?"
		args = append(args, statusID)
	}
	if startDate != "" {
		query += filterWhenDateGe
		args = append(args, startDate)
	}
	if endDate != "" {
		query += filterWhenDateLe
		args = append(args, endDate)
	}
	if iirID != nil {
		query += " AND a.iir_id = ?"
		args = append(args, *iirID)
	}

	return query, args
}

func (r *Repository) List(
	ctx context.Context,
	offset, limit int,
	search, orderBy, statusIDs, startDate, endDate string,
) ([]AppointmentWithDetailsView, error) {
	query := appointmentsBaseQuery + " WHERE 1=1"
	var args []interface{}

	if statusIDs != "" {
		statusIDList := strings.Split(statusIDs, ",")
		query += " AND a.status_id IN (?)"
		args = append(args, statusIDList)
	}
	query, args = r.applyFilters(
		query,
		args,
		"", // statusIDs handled separately with IN clause
		startDate,
		endDate,
		nil,
	)

	if search != "" {
		query += ` AND (u.first_name LIKE ? OR
			u.middle_name LIKE ? OR u.last_name LIKE ? OR
			u.email LIKE ?)`
		searchTerm := "%" + search + "%"
		args = append(
			args,
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query += fmt.Sprintf(
		" ORDER BY a.%s DESC LIMIT %d OFFSET %d",
		orderBy,
		limit,
		offset,
	)

	expandedQuery, expandedArgs, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	finalQuery := r.db.Rebind(expandedQuery)

	var appts []AppointmentWithDetailsView
	err = r.db.SelectContext(
		ctx,
		&appts,
		finalQuery,
		expandedArgs...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appts, nil
}


func (r *Repository) IsSlotAvailableForUpdate(
	ctx context.Context,
	tx datastore.DB,
	date string,
	timeSlotID int,
) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM appointments
		WHERE when_date = ?
			AND time_slot_id = ?
			AND status_id != (
				SELECT id
				FROM statuses
				WHERE name = 'Cancelled'
			)
		FOR UPDATE
	`
	err := tx.GetContext(ctx, &count, query, date, timeSlotID)
	return count == 0, err
}

func (r *Repository) GetAvailableTimeSlots(
	ctx context.Context,
	date string,
) ([]AvailableTimeSlotView, error) {
	query := `
		SELECT
			ts.id as time_slot_id,
            ts.time,
            (a.id IS NULL) as is_available
        FROM time_slots ts
        LEFT JOIN appointments a ON ts.id = a.time_slot_id
            AND a.when_date = ?
            AND a.status_id != (SELECT id FROM statuses WHERE name = 'Cancelled')
	`

	var slots []AvailableTimeSlotView
	err := r.db.SelectContext(ctx, &slots, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get available slots: %w", err)
	}

	return slots, nil
}

func (r *Repository) GetStatuses(
	ctx context.Context,
) ([]AppointmentStatus, error) {
	query := `
		SELECT id, name
		FROM statuses
		WHERE status_type IN ('appointment', 'both')
	`
	var statuses []AppointmentStatus
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}

	return statuses, nil
}


func (r *Repository) ListByIIRID(
	ctx context.Context,
	iirID string,
	offset, limit int,
	orderBy string,
	statusID, startDate, endDate string,
) ([]AppointmentWithDetailsView, error) {
	query, args := r.applyFilters(
		appointmentsBaseQuery+" WHERE a.iir_id = ?",
		[]interface{}{iirID},
		statusID,
		startDate,
		endDate,
		nil,
	)

	query += fmt.Sprintf(
		" ORDER BY a.%s DESC LIMIT %d OFFSET %d",
		orderBy,
		limit,
		offset,
	)

	var appts []AppointmentWithDetailsView
	err := r.db.SelectContext(ctx, &appts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}

	return appts, nil
}

func (r *Repository) GetAppointmentStats(
	ctx context.Context,
	statusID, startDate, endDate string,
	iirID *string,
) ([]StatusCount, error) {
	joinCondition := "a.status_id = as2.id"
	var args []interface{}

	if statusID != "" {
		joinCondition += " AND a.status_id = ?"
		args = append(args, statusID)
	}

	if startDate != "" {
		joinCondition += " AND a.when_date >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		joinCondition += " AND a.when_date <= ?"
		args = append(args, endDate)
	}

	if iirID != nil {
		joinCondition += " AND a.iir_id = ?"
		args = append(args, *iirID)
	}

	query := fmt.Sprintf(`
		SELECT
			as2.id AS id,
			as2.name AS name,
			COUNT(a.id) AS count
		FROM statuses as2
		LEFT JOIN appointments a ON %s
		WHERE as2.status_type IN ('appointment', 'both')
		GROUP BY as2.id, as2.name
		ORDER BY as2.id
	`, joinCondition)

	var counts []StatusCount
	err := r.db.SelectContext(ctx, &counts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return counts, nil
}

func (r *Repository) CreateAppointment(
	ctx context.Context,
	tx datastore.DB,
	appt *Appointment,
) (*AppointmentWithDetailsView, error) {
	query := `
		INSERT INTO appointments (
			id, iir_id, reason, admin_notes, when_date,
			time_slot_id, appointment_category_id, status_id,
			urgency_level, urgency_score
		) VALUES (
			:id, :iir_id, :reason, :admin_notes, :when_date,
			:time_slot_id, :appointment_category_id, :status_id,
			:urgency_level, :urgency_score
		)
	`

	_, err := tx.NamedExecContext(ctx, query, appt)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	newAppt, err := r.GetAppointment(ctx, tx, appt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created appointment: %w", err)
	}

	return newAppt, nil
}

func (r *Repository) UpdateAppointment(
	ctx context.Context,
	tx datastore.DB,
	appt Appointment,
) error {
	var args []interface{}
	var setQuery []string

	if appt.IIRID != "" {
		setQuery = append(setQuery, "iir_id = ?")
		args = append(args, appt.IIRID)
	}
	if appt.Reason.Valid {
		setQuery = append(setQuery, "reason = ?")
		args = append(args, appt.Reason.String)
	}
	if appt.AdminNotes.Valid {
		setQuery = append(setQuery, "admin_notes = ?")
		args = append(args, appt.AdminNotes.String)
	}
	if appt.WhenDate != "" {
		setQuery = append(setQuery, "when_date = ?")
		args = append(args, appt.WhenDate)
	}
	if appt.TimeSlotID != 0 {
		setQuery = append(setQuery, "time_slot_id = ?")
		args = append(args, appt.TimeSlotID)
	}
	if appt.StatusID != 0 {
		setQuery = append(setQuery, "status_id = ?")
		args = append(args, appt.StatusID)
	}

	if len(setQuery) == 0 {
		return nil
	}

	query := "UPDATE appointments SET " +
		strings.Join(setQuery, ", ") +
		" WHERE id = ?"
	args = append(args, appt.ID)

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
