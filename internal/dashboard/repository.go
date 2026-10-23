package dashboard

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// identifierPattern is the last line of defense before a field/group_by
// string is embedded into a query: Postgres can't parameterize identifiers,
// so this rejects anything that isn't already-validated-by-the-service plain
// snake_case (the shape every field_definitions.name actually has).
var identifierPattern = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

func (r *Repository) ListAll(ctx context.Context) ([]DashboardDefinition, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.id, d.object_id, o.name, d.name, d.label, d.widgets, d.created_at, d.updated_at
		FROM dashboard_definitions d JOIN object_definitions o ON o.id = d.object_id
		WHERE d.deleted_at IS NULL ORDER BY o.name`)
	if err != nil {
		return nil, fmt.Errorf("list dashboards: %w", err)
	}
	defer rows.Close()

	dashboards := []DashboardDefinition{}
	for rows.Next() {
		var d DashboardDefinition
		if err := rows.Scan(&d.ID, &d.ObjectID, &d.ObjectName, &d.Name, &d.Label, &d.Widgets, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan dashboard: %w", err)
		}
		dashboards = append(dashboards, d)
	}
	return dashboards, rows.Err()
}

func (r *Repository) GetByObjectName(ctx context.Context, objectName string) (*DashboardDefinition, error) {
	var d DashboardDefinition
	err := r.db.QueryRowContext(ctx, `
		SELECT d.id, d.object_id, o.name, d.name, d.label, d.widgets, d.created_at, d.updated_at
		FROM dashboard_definitions d JOIN object_definitions o ON o.id = d.object_id
		WHERE o.name = $1 AND d.deleted_at IS NULL`, objectName,
	).Scan(&d.ID, &d.ObjectID, &d.ObjectName, &d.Name, &d.Label, &d.Widgets, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query dashboard: %w", err)
	}
	return &d, nil
}

// ExecuteWidget runs the widget's aggregation against records.data for
// objectID. field and groupBy must already be validated against the
// object's real field_definitions by the caller (schema.ValidateData's
// sibling check in service.go) — this is the defense-in-depth identifier
// check, not the primary one.
func (r *Repository) ExecuteWidget(ctx context.Context, objectID int64, metric, field, groupBy string) ([]WidgetRow, error) {
	if field != "" && !identifierPattern.MatchString(field) {
		return nil, fmt.Errorf("invalid field identifier")
	}
	if groupBy != "" && !identifierPattern.MatchString(groupBy) {
		return nil, fmt.Errorf("invalid group_by identifier")
	}

	var aggExpr string
	switch metric {
	case "count":
		aggExpr = "count(*)"
	case "sum":
		aggExpr = fmt.Sprintf("coalesce(sum((data->>'%s')::numeric), 0)", field)
	case "avg":
		aggExpr = fmt.Sprintf("coalesce(avg((data->>'%s')::numeric), 0)", field)
	default:
		return nil, fmt.Errorf("unsupported metric: %s", metric)
	}

	var query string
	if groupBy != "" {
		query = fmt.Sprintf(`
			SELECT data->>'%s' AS group_key, %s AS value
			FROM records WHERE object_id = $1 AND deleted_at IS NULL
			GROUP BY data->>'%s' ORDER BY value DESC`, groupBy, aggExpr, groupBy)
	} else {
		query = fmt.Sprintf(`
			SELECT NULL AS group_key, %s AS value
			FROM records WHERE object_id = $1 AND deleted_at IS NULL`, aggExpr)
	}

	rows, err := r.db.QueryContext(ctx, query, objectID)
	if err != nil {
		return nil, fmt.Errorf("execute widget query: %w", err)
	}
	defer rows.Close()

	results := []WidgetRow{}
	for rows.Next() {
		var row WidgetRow
		if err := rows.Scan(&row.GroupKey, &row.Value); err != nil {
			return nil, fmt.Errorf("scan widget row: %w", err)
		}
		results = append(results, row)
	}
	return results, rows.Err()
}
