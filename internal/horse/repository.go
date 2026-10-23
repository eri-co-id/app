package horse

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

func scanHorse(row interface{ Scan(...any) error }) (*Horse, error) {
	h := &Horse{}
	err := row.Scan(
		&h.ID, &h.Name, &h.Code, &h.Breed, &h.DOB, &h.Color, &h.OwnerName,
		&h.CurrentTrainerID, &h.HealthStatus, &h.Notes, &h.PhotoURL,
		&h.CreatedAt, &h.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return h, nil
}

const horseColumns = `id, name, code, breed, dob, color, owner_name,
	current_trainer_id, health_status, notes, photo_url, created_at, updated_at`

func (r *Repository) Create(ctx context.Context, req *UpsertHorseRequest) (*Horse, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO horses (name, code, breed, dob, color, owner_name, current_trainer_id, health_status, notes, photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING `+horseColumns,
		req.Name, req.Code, req.Breed, req.DOB, req.Color, req.OwnerName,
		req.CurrentTrainerID, req.HealthStatus, req.Notes, req.PhotoURL,
	)
	h, err := scanHorse(row)
	if err != nil {
		return nil, fmt.Errorf("insert horse: %w", err)
	}
	return h, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Horse, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+horseColumns+`
		FROM horses WHERE id = $1 AND deleted_at IS NULL`, id)
	h, err := scanHorse(row)
	if err != nil {
		return nil, fmt.Errorf("query horse: %w", err)
	}
	return h, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Horse, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+horseColumns+`
		FROM horses WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list horses: %w", err)
	}
	defer rows.Close()

	horses := []Horse{}
	for rows.Next() {
		h, err := scanHorse(rows)
		if err != nil {
			return nil, fmt.Errorf("scan horse: %w", err)
		}
		horses = append(horses, *h)
	}
	return horses, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id int64, req *UpsertHorseRequest) (*Horse, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE horses SET
			name = $1, code = $2, breed = $3, dob = $4, color = $5, owner_name = $6,
			current_trainer_id = $7, health_status = $8, notes = $9, photo_url = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND deleted_at IS NULL
		RETURNING `+horseColumns,
		req.Name, req.Code, req.Breed, req.DOB, req.Color, req.OwnerName,
		req.CurrentTrainerID, req.HealthStatus, req.Notes, req.PhotoURL, id,
	)
	h, err := scanHorse(row)
	if err != nil {
		return nil, fmt.Errorf("update horse: %w", err)
	}
	return h, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE horses SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false, fmt.Errorf("soft-delete horse: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return rows > 0, nil
}
