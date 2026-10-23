package dashboard

import (
	"encoding/json"
	"testing"
)

func TestParseWidgets(t *testing.T) {
	d := DashboardDefinition{Widgets: json.RawMessage(`[{"id":"total","type":"stat","title":"Total","metric":"count"}]`)}
	widgets, err := d.ParseWidgets()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(widgets) != 1 || widgets[0].ID != "total" || widgets[0].Metric != "count" {
		t.Fatalf("unexpected widgets: %+v", widgets)
	}
}

func TestParseWidgetsEmpty(t *testing.T) {
	d := DashboardDefinition{}
	widgets, err := d.ParseWidgets()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(widgets) != 0 {
		t.Fatalf("expected no widgets, got %+v", widgets)
	}
}
