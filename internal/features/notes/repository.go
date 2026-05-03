package notes

import (
	"context"
	"fmt"

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

func (r *Repository) GetStudentSignificantNotes(
	ctx context.Context,
	iirID string,
) ([]SignificantNote, error) {
	query := `
		SELECT id, iir_id, appointment_id, admission_slip_id, note, remarks, created_at, updated_at
		FROM significant_notes
		WHERE iir_id = ?
	`

	var results []SignificantNote
	err := r.db.SelectContext(ctx, &results, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student significant notes: %w",
			err,
		)
	}

	return results, nil
}

func (r *Repository) CreateSignificantNote(
	ctx context.Context,
	sn *SignificantNote,
) (string, error) {
	return datastore.NewRunInTransaction(
		ctx,
		r.db,
		func(tx datastore.DB) (string, error) {
			query := `
				INSERT INTO significant_notes (
					id, iir_id, appointment_id, admission_slip_id, note, remarks, updated_at
				) VALUES (
					:id, :iir_id, :appointment_id, :admission_slip_id, :note, :remarks, NOW()
				)
			`

			_, err := tx.NamedExecContext(
				ctx,
				query,
				sn,
			)
			if err != nil {
				return "", fmt.Errorf(
					"failed to create significant note: %w",
					err,
				)
			}

			return sn.ID, nil
		},
	)
}

func (r *Repository) HasNoteForAppointment(
	ctx context.Context,
	appointmentID string,
) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM significant_notes WHERE appointment_id = ?"
	err := r.db.GetContext(ctx, &count, query, appointmentID)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check for existing note: %w",
			err,
		)
	}

	return count > 0, nil
}

