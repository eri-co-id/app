package lessoncode

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

const columns = `id, code, type, default_fee, duration_hours, description, active, created_at, updated_at`

func scanLessonCode(row interface{ Scan(...any) error }) (*LessonCode, error) {
	lc := &LessonCode{}
	err := row.Scan(
		&lc.ID, &lc.Code, &lc.Type, &lc.DefaultFee, &lc.DurationHours,
		&lc.Description, &lc.Active, &lc.CreatedAt, &lc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return lc, nil
}

func (r *Repository) Create(ctx context.Context, req *UpsertLessonCodeRequest) (*LessonCode, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO lesson_codes (code, type, default_fee, duration_hours, description, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+columns,
		req.Code, req.Type, req.DefaultFee, req.DurationHours, req.Description, req.Active,
	)
	lc, err := scanLessonCode(row)
	if err != nil {
		return nil, fmt.Errorf("insert lesson code: %w", err)
	}
	return lc, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*LessonCode, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+columns+` FROM lesson_codes WHERE id = $1`, id)
	lc, err := scanLessonCode(row)
	if err != nil {
		return nil, fmt.Errorf("query lesson code: %w", err)
	}
	return lc, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]LessonCode, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+columns+`
		FROM lesson_codes
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list lesson codes: %w", err)
	}
	defer rows.Close()

	items := []LessonCode{}
	for rows.Next() {
		lc, err := scanLessonCode(rows)
		if err != nil {
			return nil, fmt.Errorf("scan lesson code: %w", err)
		}
		items = append(items, *lc)
	}
	return items, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id int64, req *UpsertLessonCodeRequest) (*LessonCode, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE lesson_codes SET
			code = $1, type = $2, default_fee = $3, duration_hours = $4,
			description = $5, active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING `+columns,
		req.Code, req.Type, req.DefaultFee, req.DurationHours, req.Description, req.Active, id,
	)
	lc, err := scanLessonCode(row)
	if err != nil {
		return nil, fmt.Errorf("update lesson code: %w", err)
	}
	return lc, nil
}

// Deactivate marks a lesson code inactive. Lesson codes are kept forever
// (spec: "deactivate, keep history") rather than soft-deleted like other master data.
func (r *Repository) Deactivate(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lesson_codes SET active = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("deactivate lesson code: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return rows > 0, nil
}
