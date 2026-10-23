package schema

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

func (r *Repository) ListObjects(ctx context.Context) ([]ObjectDefinition, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, label, plural_label, icon, description
		FROM object_definitions WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list objects: %w", err)
	}
	defer rows.Close()

	objects := []ObjectDefinition{}
	for rows.Next() {
		var o ObjectDefinition
		if err := rows.Scan(&o.ID, &o.Name, &o.Label, &o.PluralLabel, &o.Icon, &o.Description); err != nil {
			return nil, fmt.Errorf("scan object: %w", err)
		}
		objects = append(objects, o)
	}
	return objects, rows.Err()
}

func (r *Repository) GetObjectByName(ctx context.Context, name string) (*ObjectDefinition, error) {
	var o ObjectDefinition
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, label, plural_label, icon, description
		FROM object_definitions WHERE name = $1 AND deleted_at IS NULL`, name,
	).Scan(&o.ID, &o.Name, &o.Label, &o.PluralLabel, &o.Icon, &o.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query object: %w", err)
	}
	return &o, nil
}

func (r *Repository) ListFieldsByObjectID(ctx context.Context, objectID int64) ([]FieldDefinition, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, object_id, name, label, type, required, position, metadata
		FROM field_definitions WHERE object_id = $1 ORDER BY position`, objectID)
	if err != nil {
		return nil, fmt.Errorf("list fields: %w", err)
	}
	defer rows.Close()

	fields := []FieldDefinition{}
	for rows.Next() {
		var f FieldDefinition
		if err := rows.Scan(&f.ID, &f.ObjectID, &f.Name, &f.Label, &f.Type, &f.Required, &f.Position, &f.Metadata); err != nil {
			return nil, fmt.Errorf("scan field: %w", err)
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}
