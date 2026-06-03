package files

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

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) GetFileByID(
	ctx context.Context,
	id string,
) (*File, error) {
	var file File
	query := `SELECT id, file_name, file_url, file_type, file_size, mime_type, created_at, updated_at, deleted_at FROM files WHERE id = ?`
	err := r.db.GetContext(ctx, &file, query, id)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *Repository) GetFilesByIDs(
	ctx context.Context,
	ids []string,
) ([]File, error) {
	var files []File
	query, args, err := sqlx.In(
		`SELECT id, file_name, file_url, file_type, file_size, mime_type, created_at, updated_at, deleted_at FROM files WHERE id IN (?)`,
		ids,
	)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &files, query, args...)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *Repository) Create(
	ctx context.Context,
	tx datastore.DB,
	file File,
) (string, error) {
	query := `
		INSERT INTO files (
			id, file_name, file_url, file_type, file_size, mime_type, updated_at
		) VALUES (
			:id, :file_name, :file_url, :file_type, :file_size, :mime_type, NOW()
		)
	`
	_, err := tx.NamedExecContext(ctx, query, &file)
	if err != nil {
		return "", err
	}

	return file.ID, nil
}

func (r *Repository) CreateBulk(
	ctx context.Context,
	tx datastore.DB,
	files []File,
) ([]string, error) {
	var ids []string
	for _, file := range files {
		id, err := r.Create(ctx, tx, file)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repository) GetFileByURL(
	ctx context.Context,
	fileURL string,
) (*File, error) {
	var file File
	query := `SELECT id, file_name, file_url, file_type, file_size, mime_type, created_at, updated_at, deleted_at FROM files WHERE file_url = ?`
	err := r.db.GetContext(ctx, &file, query, fileURL)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	tx datastore.DB,
	id string,
) error {
	query := `DELETE FROM files WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) SaveOCRResult(
	ctx context.Context,
	tx datastore.DB,
	result OCRResult,
) error {
	query := `
		INSERT INTO file_ocr_results (
			file_id, raw_text, structured_data, engine_v, confidence
		) VALUES (
			:file_id, :raw_text, :structured_data, :engine_v, :confidence
		)
		ON DUPLICATE KEY UPDATE
			raw_text = VALUES(raw_text),
			structured_data = VALUES(structured_data),
			engine_v = VALUES(engine_v),
			confidence = VALUES(confidence)
	`
	_, err := tx.NamedExecContext(ctx, query, result)
	if err != nil {
		return fmt.Errorf("[OCRRepository] {SaveOCRResult}: %w", err)
	}
	return nil
}

func (r *Repository) GetOCRResult(
	ctx context.Context,
	fileID string,
) (*OCRResult, error) {
	var result OCRResult
	query := "SELECT file_id, raw_text, structured_data, engine_v, confidence, created_at FROM file_ocr_results WHERE file_id = ?"
	err := r.db.GetContext(ctx, &result, query, fileID)
	if err != nil {
		return nil, fmt.Errorf(
			"[OCRRepository] {GetOCRResult}: %w",
			err,
		)
	}
	return &result, nil
}
