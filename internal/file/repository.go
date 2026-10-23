package file

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const columns = `id, record_id, file_name, file_url, file_type, file_size, uploaded_by_id, created_at`

func scanFile(row interface{ Scan(...any) error }) (*File, error) {
	f := &File{}
	err := row.Scan(&f.ID, &f.RecordID, &f.FileName, &f.FileURL, &f.FileType, &f.FileSize, &f.UploadedByID, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Repository) Create(ctx context.Context, recordID int64, fileName, fileURL, fileType string, fileSize int64, uploadedByID *int64) (*File, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO files (record_id, file_name, file_url, file_type, file_size, uploaded_by_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+columns,
		recordID, fileName, fileURL, fileType, fileSize, uploadedByID,
	)
	f, err := scanFile(row)
	if err != nil {
		return nil, fmt.Errorf("insert file: %w", err)
	}
	return f, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*File, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+columns+` FROM files WHERE id = $1 AND deleted_at IS NULL`, id)
	f, err := scanFile(row)
	if err != nil {
		return nil, fmt.Errorf("query file: %w", err)
	}
	return f, nil
}

func (r *Repository) ListByRecordID(ctx context.Context, recordID int64) ([]File, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+columns+` FROM files WHERE record_id = $1 AND deleted_at IS NULL ORDER BY created_at`, recordID)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}
	defer rows.Close()

	files := []File{}
	for rows.Next() {
		f, err := scanFile(rows)
		if err != nil {
			return nil, fmt.Errorf("scan file: %w", err)
		}
		files = append(files, *f)
	}
	return files, rows.Err()
}

func (r *Repository) CountByRecordID(ctx context.Context, recordID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*) FROM files WHERE record_id = $1 AND deleted_at IS NULL`, recordID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count files: %w", err)
	}
	return count, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) (*File, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE files SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING `+columns,
		id,
	)
	f, err := scanFile(row)
	if err != nil {
		return nil, fmt.Errorf("soft-delete file: %w", err)
	}
	return f, nil
}
