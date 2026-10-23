package record

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const columns = `r.id, r.data, r.created_by_id, u.name, r.created_at, r.updated_at`
const fromClause = `records r LEFT JOIN users u ON u.id = r.created_by_id`

func scanRecord(row interface{ Scan(...any) error }) (*Record, error) {
	rec := &Record{}
	err := row.Scan(&rec.ID, &rec.Data, &rec.CreatedByID, &rec.CreatedByName, &rec.CreatedAt, &rec.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *Repository) Create(ctx context.Context, objectID int64, data map[string]any, createdByID int64) (*Record, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal data: %w", err)
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO records (object_id, data, created_by_id)
		VALUES ($1, $2, $3)
		RETURNING id, data, created_by_id, (SELECT name FROM users WHERE id = $3), created_at, updated_at`,
		objectID, dataJSON, createdByID,
	)
	rec, err := scanRecord(row)
	if err != nil {
		return nil, fmt.Errorf("insert record: %w", err)
	}
	return rec, nil
}

func (r *Repository) GetByID(ctx context.Context, objectID, id int64) (*Record, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+columns+` FROM `+fromClause+`
		WHERE r.id = $1 AND r.object_id = $2 AND r.deleted_at IS NULL`, id, objectID)
	rec, err := scanRecord(row)
	if err != nil {
		return nil, fmt.Errorf("query record: %w", err)
	}
	return rec, nil
}

// List returns records for objectID, optionally filtered to one creator
// (used to scope non-admin roles to their own records — pass nil to see all).
func (r *Repository) List(ctx context.Context, objectID int64, createdByID *int64, limit, offset int) ([]Record, error) {
	query := `SELECT ` + columns + ` FROM ` + fromClause + ` WHERE r.object_id = $1 AND r.deleted_at IS NULL`
	args := []any{objectID}

	if createdByID != nil {
		query += fmt.Sprintf(" AND r.created_by_id = $%d", len(args)+1)
		args = append(args, *createdByID)
	}

	query += fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list records: %w", err)
	}
	defer rows.Close()

	records := []Record{}
	for rows.Next() {
		rec, err := scanRecord(rows)
		if err != nil {
			return nil, fmt.Errorf("scan record: %w", err)
		}
		records = append(records, *rec)
	}
	return records, rows.Err()
}

func (r *Repository) Update(ctx context.Context, objectID, id int64, data map[string]any) (*Record, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal data: %w", err)
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE records r SET data = $1, updated_at = CURRENT_TIMESTAMP
		WHERE r.id = $2 AND r.object_id = $3 AND r.deleted_at IS NULL
		RETURNING r.id, r.data, r.created_by_id, (SELECT name FROM users WHERE id = r.created_by_id), r.created_at, r.updated_at`,
		dataJSON, id, objectID,
	)
	rec, err := scanRecord(row)
	if err != nil {
		return nil, fmt.Errorf("update record: %w", err)
	}
	return rec, nil
}

func (r *Repository) SoftDelete(ctx context.Context, objectID, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE records SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND object_id = $2 AND deleted_at IS NULL`, id, objectID)
	if err != nil {
		return false, fmt.Errorf("soft-delete record: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return rowsAffected > 0, nil
}
