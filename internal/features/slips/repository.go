package slips

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

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

const (
	slipsBaseQuery = `
	SELECT
		slp.id AS id,
		u.id AS user_id,
		slp.iir_id AS iir_id,
		u.first_name AS user_first_name,
		u.middle_name AS user_middle_name,
		u.last_name AS user_last_name,
		u.email AS user_email,
		spi.student_number AS student_number,
		slp.reason AS reason,
		DATE_FORMAT(slp.date_of_absence, '%Y-%m-%d') AS date_of_absence,
		DATE_FORMAT(slp.date_needed, '%Y-%m-%d') AS date_needed,
		slp.admin_notes AS admin_notes,
		c.id AS category_id,
		c.name AS category_name,
		slp.status_id AS status_id,
		s.name AS status_name,
		t.ticket_code AS ticket_code,
		t.is_verified AS is_verified,
		t.verified_at AS verified_at,
		slp.created_at AS created_at,
		slp.updated_at AS updated_at
	FROM admission_slips slp
	JOIN iir_records ir ON slp.iir_id = ir.id
	JOIN student_personal_info spi ON ir.id = spi.iir_id
	JOIN users u ON ir.user_id = u.id
	JOIN admission_slip_categories c ON slp.category_id = c.id
	JOIN statuses s ON slp.status_id = s.id
	LEFT JOIN admission_tickets t ON slp.id = t.admission_slip_id
`
	orderSlipsCreatedDesc = " ORDER BY slp.date_needed ASC LIMIT ? OFFSET ?"
)

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) CreateSlip(
	ctx context.Context,
	tx datastore.DB,
	slip *Slip,
) (*SlipWithDetailsView, error) {
	query := `
		INSERT INTO admission_slips (
			id, iir_id, reason, date_of_absence, date_needed,
			category_id, status_id
		) VALUES (
			:id, :iir_id, :reason, :date_of_absence, :date_needed,
			:category_id, :status_id
		)
	`

	_, err := tx.NamedExecContext(ctx, query, slip)
	if err != nil {
		return nil, fmt.Errorf("failed to insert excuse slip: %w", err)
	}

	return r.GetSlipByIDWithDetails(ctx, tx, slip.ID)
}

func (r *Repository) SaveSlipAttachment(
	ctx context.Context,
	tx datastore.DB,
	attachment *SlipAttachment,
) error {
	query := `
			INSERT INTO slip_attachments (file_id, admission_slip_id, attachment_type)
			VALUES (?, ?, ?)
		`

	_, err := tx.ExecContext(
		ctx,
		query,
		attachment.FileID,
		attachment.SlipID,
		attachment.AttachmentType,
	)
	if err != nil {
		return fmt.Errorf("failed to insert slip attachment: %w", err)
	}

	return nil
}

func (r *Repository) UpdateSlip(
	ctx context.Context,
	tx datastore.DB,
	slip *Slip,
) error {
	query := `
		UPDATE admission_slips
		SET reason = ?, date_of_absence = ?, date_needed = ?, category_id = ?, status_id = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		slip.Reason,
		slip.DateOfAbsence,
		slip.DateNeeded,
		slip.CategoryID,
		slip.StatusID,
		slip.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update slip: %w", err)
	}

	return nil
}

func (r *Repository) DeleteSlipAttachments(
	ctx context.Context,
	tx datastore.DB,
	slipID string,
) error {
	query := `DELETE FROM slip_attachments WHERE admission_slip_id = ?`
	_, err := tx.ExecContext(ctx, query, slipID)
	if err != nil {
		return fmt.Errorf("failed to delete slip attachments: %w", err)
	}
	return nil
}

func (r *Repository) GetSlipStatuses(
	ctx context.Context,
) ([]SlipStatus, error) {
	var statuses []SlipStatus
	query := `SELECT id, name FROM statuses WHERE status_type IN ('slip', 'both')`
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip statuses: %w", err)
	}

	return statuses, nil
}

func (r *Repository) GetSlipCategories(
	ctx context.Context,
) ([]SlipCategory, error) {
	var categories []SlipCategory
	query := `SELECT id, name FROM admission_slip_categories`
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip categories: %w", err)
	}

	return categories, nil
}

func (r *Repository) GetSlipStats(
	ctx context.Context,
	iirID *string,
	req *ListSlipsRequest,
) ([]SlipStatusCount, error) {
	filterConditions, args := r.applyFilters("1=1", nil, req, iirID)

	var counts []SlipStatusCount
	query := fmt.Sprintf(`
		SELECT
			s.id AS id,
			s.name AS name,
			COUNT(slp.id) AS count
		FROM statuses s
		LEFT JOIN admission_slips slp
			ON s.id = slp.status_id AND %s
		WHERE s.status_type IN ('slip', 'both')
		GROUP BY s.id, s.name
	`, filterConditions)
	err := r.db.SelectContext(ctx, &counts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}

	return counts, nil
}

func (r *Repository) applyFilters(
	query string,
	args []interface{},
	req *ListSlipsRequest,
	iirID *string,
) (string, []interface{}) {
	if args == nil {
		args = []interface{}{}
	}

	if req.StatusID != 0 {
		query += " AND slp.status_id = ?"
		args = append(args, req.StatusID)
	}

	if req.StartDate != "" {
		query += " AND slp.date_needed >= ?"
		args = append(args, req.StartDate)
	}

	if req.EndDate != "" {
		query += " AND slp.date_needed <= ?"
		args = append(args, req.EndDate)
	}

	if req.Search != "" {
		query += `
			AND (
				slp.reason LIKE ?
				OR u.first_name LIKE ?
				OR u.last_name LIKE ?
			)`
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if iirID != nil {
		query += " AND slp.iir_id = ?"
		args = append(args, iirID)
	}

	return query, args
}

func (r *Repository) GetTotalSlipsCount(
	ctx context.Context,
	req *ListSlipsRequest,
	iirID *string,
) (int, error) {
	query, args := r.applyFilters(
		`SELECT COUNT(*) FROM admission_slips slp
		 JOIN iir_records ir ON slp.iir_id = ir.id
		 JOIN users u ON ir.user_id = u.id
		 WHERE 1=1`,
		nil,
		req,
		iirID,
	)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count slips: %w", err)
	}

	return count, nil
}

func (r *Repository) GetTotalUrgentSlipsCount(
	ctx context.Context,
	req *ListSlipsRequest,
) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM admission_slips slp
        WHERE slp.status_id IN (1, 9)
    `
	var args []interface{}

	if req.StartDate != "" {
		query += " AND slp.date_needed >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		query += " AND slp.date_needed <= ?"
		args = append(args, req.EndDate)
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count urgent slips: %w", err)
	}

	return count, nil
}

func (r *Repository) GetUrgentSlips(
	ctx context.Context,
	req *ListSlipsRequest,
) ([]SlipWithDetailsView, error) {
	query := strings.Replace(
		slipsBaseQuery,
		"slp.updated_at AS updated_at",
		`slp.updated_at AS updated_at, `+
			r.getUrgencyScoreSQL()+" AS urgency_score",
		1,
	)

	query += " WHERE slp.status_id IN (1, 9)"
	var args []interface{}

	if req.StartDate != "" {
		query += " AND slp.date_needed >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		query += " AND slp.date_needed <= ?"
		args = append(args, req.EndDate)
	}

	query += `
		ORDER BY
			slp.date_needed ASC,
			urgency_score DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, req.PageSize, req.GetOffset())

	var slips []SlipWithDetailsView
	err := r.db.SelectContext(ctx, &slips, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get urgent slips: %w", err)
	}

	return slips, nil
}

func (r *Repository) getUrgencyScoreSQL() string {
	return `(
		(1000 - DATEDIFF(slp.date_needed, CURRENT_DATE)) * 10
		+
		CASE WHEN slp.category_id = 1 THEN 500 ELSE 0 END
	)`
}

func (r *Repository) GetAll(
	ctx context.Context,
	req *ListSlipsRequest,
) ([]SlipWithDetailsView, error) {
	query, args := r.applyFilters(
		slipsBaseQuery+" WHERE 1=1",
		nil,
		req,
		nil,
	)

	orderCol, orderDir := sanitizeSort(req.OrderBy, req.SortOrder)
	query += fmt.Sprintf(
		" ORDER BY %s %s LIMIT ? OFFSET ?",
		orderCol,
		orderDir,
	)
	args = append(args, req.PageSize, req.GetOffset())

	var slips []SlipWithDetailsView
	err := r.db.SelectContext(ctx, &slips, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get excuse slips: %w", err)
	}

	return slips, nil
}

func (r *Repository) GetByIIRID(
	ctx context.Context,
	iirID string,
	req *ListSlipsRequest,
) ([]SlipWithDetailsView, error) {
	query, args := r.applyFilters(
		slipsBaseQuery+" WHERE slp.iir_id = ?",
		[]interface{}{iirID},
		req,
		nil,
	)

	orderCol, orderDir := sanitizeSort(req.OrderBy, req.SortOrder)
	query += fmt.Sprintf(
		" ORDER BY %s %s LIMIT ? OFFSET ?",
		orderCol,
		orderDir,
	)
	args = append(args, req.PageSize, req.GetOffset())

	var slips []SlipWithDetailsView
	err := r.db.SelectContext(ctx, &slips, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get slips for IIR: %w", err)
	}

	return slips, nil
}

func (r *Repository) GetSlipByID(
	ctx context.Context,
	id string,
) (*Slip, error) {
	var slip Slip
	query := `SELECT id, iir_id, reason, date_of_absence, date_needed, admin_notes, category_id, status_id, created_at, updated_at FROM admission_slips WHERE id = ?`
	err := r.db.GetContext(ctx, &slip, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get slip: %w", err)
	}
	return &slip, nil
}

func (r *Repository) GetSlipByIDWithDetails(
	ctx context.Context,
	db datastore.DB,
	id string,
) (*SlipWithDetailsView, error) {
	var slip SlipWithDetailsView
	query := slipsBaseQuery + " WHERE slp.id = ?"
	err := db.GetContext(ctx, &slip, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get slip with details: %w", err)
	}
	return &slip, nil
}

func (r *Repository) GetUserIDBySlipID(
	ctx context.Context,
	id string,
) (string, error) {
	var userID string
	query := `
		SELECT ir.user_id
		FROM admission_slips slp
		JOIN iir_records ir ON slp.iir_id = ir.id
		WHERE slp.id = ?
	`
	err := r.db.GetContext(ctx, &userID, query, id)
	return userID, err
}

func (r *Repository) GetSlipAttachments(
	ctx context.Context,
	slipID string,
) ([]SlipAttachment, error) {
	var attachments []SlipAttachment
	query := `
		SELECT
			sa.file_id,
			sa.admission_slip_id,
			sa.attachment_type,
			f.file_name,
			f.file_url
		FROM slip_attachments sa
		JOIN files f ON sa.file_id = f.id
		WHERE sa.admission_slip_id = ?
	`
	err := r.db.SelectContext(ctx, &attachments, query, slipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip attachments: %w", err)
	}

	return attachments, nil
}

func (r *Repository) GetAttachmentByID(
	ctx context.Context,
	attachmentID string,
) (*SlipAttachment, error) {
	var attachment SlipAttachment
	query := `
		SELECT
			sa.file_id,
			sa.admission_slip_id,
			sa.attachment_type,
			f.file_name,
			f.file_url
		FROM slip_attachments sa
		JOIN files f ON sa.file_id = f.id
		WHERE sa.file_id = ?
	`
	err := r.db.GetContext(ctx, &attachment, query, attachmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	return &attachment, nil
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	tx datastore.DB,
	id string,
	statusName string,
	adminNotes string,
) error {
	// First, get the status ID from the status name
	var statusID int
	query := `SELECT id FROM statuses WHERE name = ?`
	err := tx.GetContext(ctx, &statusID, query, statusName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("status '%s' not found", statusName)
		}
		return fmt.Errorf("failed to get status ID: %w", err)
	}

	// Now update the slip with the status ID and admin notes
	updateQuery := `
		UPDATE admission_slips
		SET status_id = ?, admin_notes = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := tx.ExecContext(ctx, updateQuery, statusID, adminNotes, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// This usually means the ID doesn't exist
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) CreateTicket(
	ctx context.Context,
	tx datastore.DB,
	ticket *AdmissionTicket,
) error {
	query := `
		INSERT INTO admission_tickets (
			id, admission_slip_id, ticket_code
		) VALUES (
			:id, :admission_slip_id, :ticket_code
		)
	`
	_, err := tx.NamedExecContext(ctx, query, ticket)
	if err != nil {
		return fmt.Errorf("failed to insert ticket: %w", err)
	}
	return nil
}

func (r *Repository) GetTicketByCode(
	ctx context.Context,
	code string,
) (*AdmissionTicket, error) {
	var ticket AdmissionTicket
	query := `SELECT id, admission_slip_id, ticket_code, is_verified,
			  verified_at, verified_by, created_at, updated_at
			  FROM admission_tickets WHERE ticket_code = ?`
	err := r.db.GetContext(ctx, &ticket, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	return &ticket, nil
}

func (r *Repository) UpdateTicketVerification(
	ctx context.Context,
	tx datastore.DB,
	ticketID string,
	verifiedBy string,
) error {
	query := `
		UPDATE admission_tickets
		SET is_verified = TRUE, verified_at = NOW(), verified_by = ?
		WHERE id = ?
	`
	_, err := tx.ExecContext(ctx, query, verifiedBy, ticketID)
	if err != nil {
		return fmt.Errorf("failed to verify ticket: %w", err)
	}
	return nil
}

func (r *Repository) GetTicketBySlipID(
	ctx context.Context,
	slipID string,
) (*AdmissionTicket, error) {
	var ticket AdmissionTicket
	query := `SELECT id, admission_slip_id, ticket_code, is_verified,
			  verified_at, verified_by, created_at, updated_at
			  FROM admission_tickets WHERE admission_slip_id = ?`
	err := r.db.GetContext(ctx, &ticket, query, slipID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ticket by slip id: %w", err)
	}
	return &ticket, nil
}

func (r *Repository) GetSlipByTicketCode(
	ctx context.Context,
	code string,
) (*SlipWithDetailsView, error) {
	var slip SlipWithDetailsView
	query := slipsBaseQuery + ` WHERE t.ticket_code = ?`
	err := r.db.GetContext(ctx, &slip, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get slip by ticket code: %w", err)
	}
	return &slip, nil
}

func (r *Repository) HasNoteForAdmissionSlip(
	ctx context.Context,
	admissionSlipID string,
) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM significant_notes
		WHERE admission_slip_id = ?
	`
	err := r.db.GetContext(ctx, &count, query, admissionSlipID)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check for existing note: %w",
			err,
		)
	}
	return count > 0, nil
}

func sanitizeSort(orderBy, sortOrder string) (string, string) {
	allowed := map[string]string{
		"dateNeeded":      "slp.date_needed",
		"date_needed":     "slp.date_needed",
		"dateOfAbsence":   "slp.date_of_absence",
		"date_of_absence": "slp.date_of_absence",
		"createdAt":       "slp.created_at",
		"created_at":      "slp.created_at",
	}

	col, ok := allowed[orderBy]
	if !ok {
		col = "slp.date_needed"
	}

	dir := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		dir = "DESC"
	}

	return col, dir
}
