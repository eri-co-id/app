package equipment

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

const columns = `id, name, type, code, serial_batch, condition, current_horse_id,
	location, purchase_date, notes, photo_url, active, created_at, updated_at`

func scanEquipment(row interface{ Scan(...any) error }) (*Equipment, error) {
	e := &Equipment{}
	err := row.Scan(
		&e.ID, &e.Name, &e.Type, &e.Code, &e.SerialBatch, &e.Condition, &e.CurrentHorseID,
		&e.Location, &e.PurchaseDate, &e.Notes, &e.PhotoURL, &e.Active, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *Repository) Create(ctx context.Context, req *UpsertEquipmentRequest) (*Equipment, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO equipment (name, type, code, serial_batch, condition, current_horse_id,
			location, purchase_date, notes, photo_url, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING `+columns,
		req.Name, req.Type, req.Code, req.SerialBatch, req.Condition, req.CurrentHorseID,
		req.Location, req.PurchaseDate, req.Notes, req.PhotoURL, req.Active,
	)
	e, err := scanEquipment(row)
	if err != nil {
		return nil, fmt.Errorf("insert equipment: %w", err)
	}
	return e, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Equipment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+columns+`
		FROM equipment WHERE id = $1 AND deleted_at IS NULL`, id)
	e, err := scanEquipment(row)
	if err != nil {
		return nil, fmt.Errorf("query equipment: %w", err)
	}
	return e, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Equipment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+columns+`
		FROM equipment WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list equipment: %w", err)
	}
	defer rows.Close()
	return collectEquipment(rows)
}

func (r *Repository) ListByHorse(ctx context.Context, horseID int64) ([]Equipment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+columns+`
		FROM equipment WHERE current_horse_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`, horseID)
	if err != nil {
		return nil, fmt.Errorf("list equipment by horse: %w", err)
	}
	defer rows.Close()
	return collectEquipment(rows)
}

func collectEquipment(rows *sql.Rows) ([]Equipment, error) {
	items := []Equipment{}
	for rows.Next() {
		e, err := scanEquipment(rows)
		if err != nil {
			return nil, fmt.Errorf("scan equipment: %w", err)
		}
		items = append(items, *e)
	}
	return items, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id int64, req *UpsertEquipmentRequest) (*Equipment, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE equipment SET
			name = $1, type = $2, code = $3, serial_batch = $4, condition = $5,
			current_horse_id = $6, location = $7, purchase_date = $8, notes = $9,
			photo_url = $10, active = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND deleted_at IS NULL
		RETURNING `+columns,
		req.Name, req.Type, req.Code, req.SerialBatch, req.Condition, req.CurrentHorseID,
		req.Location, req.PurchaseDate, req.Notes, req.PhotoURL, req.Active, id,
	)
	e, err := scanEquipment(row)
	if err != nil {
		return nil, fmt.Errorf("update equipment: %w", err)
	}
	return e, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE equipment SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false, fmt.Errorf("soft-delete equipment: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return rows > 0, nil
}
