package dashboard

import (
	"encoding/json"
	"time"
)

type DashboardDefinition struct {
	ID         int64           `json:"id"`
	ObjectID   int64           `json:"object_id"`
	ObjectName string          `json:"object_name"`
	Name       string          `json:"name"`
	Label      string          `json:"label"`
	Widgets    json.RawMessage `json:"widgets"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// Widget is a structured, whitelisted query descriptor — never a raw SQL
// string. Metric and GroupBy/Field are validated against the object's real
// field_definitions before being used to build a query (see validate.go).
type Widget struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // "stat" | "chart" | "table"
	Title     string `json:"title"`
	Metric    string `json:"metric"` // "count" | "sum" | "avg"
	Field     string `json:"field,omitempty"`
	GroupBy   string `json:"group_by,omitempty"`
	ChartType string `json:"chart_type,omitempty"`
}

func (d DashboardDefinition) ParseWidgets() ([]Widget, error) {
	var widgets []Widget
	if len(d.Widgets) == 0 {
		return widgets, nil
	}
	if err := json.Unmarshal(d.Widgets, &widgets); err != nil {
		return nil, err
	}
	return widgets, nil
}

type WidgetRow struct {
	GroupKey *string `json:"group_key,omitempty"`
	Value    float64 `json:"value"`
}

type WidgetResult struct {
	WidgetID string      `json:"widget_id"`
	Title    string      `json:"title"`
	Rows     []WidgetRow `json:"rows"`
}

type ListResponse struct {
	Dashboards []DashboardDefinition `json:"dashboards"`
}
