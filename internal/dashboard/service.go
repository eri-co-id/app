package dashboard

import (
	"context"
	"errors"
	"fmt"

	"github.com/eri-co-id/app/internal/schema"
)

var (
	ErrNotFound       = errors.New("dashboard not found")
	ErrWidgetNotFound = errors.New("widget not found")
)

type Service struct {
	repo      *Repository
	schemaSvc *schema.Service
}

func NewService(repo *Repository, schemaSvc *schema.Service) *Service {
	return &Service{repo: repo, schemaSvc: schemaSvc}
}

func (s *Service) List(ctx context.Context) ([]DashboardDefinition, error) {
	dashboards, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list dashboards: %w", err)
	}
	return dashboards, nil
}

func (s *Service) Get(ctx context.Context, objectName string) (*DashboardDefinition, error) {
	d, err := s.repo.GetByObjectName(ctx, objectName)
	if err != nil {
		return nil, fmt.Errorf("get dashboard: %w", err)
	}
	if d == nil {
		return nil, ErrNotFound
	}
	return d, nil
}

// ExecuteWidget validates the widget's field/group_by against the object's
// real field_definitions (the primary SQL-injection defense — only a known
// column name ever reaches a query) before running it.
func (s *Service) ExecuteWidget(ctx context.Context, objectName, widgetID string) (*WidgetResult, error) {
	d, err := s.Get(ctx, objectName)
	if err != nil {
		return nil, err
	}

	widgets, err := d.ParseWidgets()
	if err != nil {
		return nil, fmt.Errorf("parse widgets: %w", err)
	}

	var widget *Widget
	for i := range widgets {
		if widgets[i].ID == widgetID {
			widget = &widgets[i]
			break
		}
	}
	if widget == nil {
		return nil, ErrWidgetNotFound
	}

	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return nil, err
	}
	validField := func(name string) bool {
		if name == "" {
			return true
		}
		for _, f := range object.Fields {
			if f.Name == name {
				return true
			}
		}
		return false
	}
	if !validField(widget.Field) {
		return nil, fmt.Errorf("widget references unknown field: %s", widget.Field)
	}
	if !validField(widget.GroupBy) {
		return nil, fmt.Errorf("widget references unknown group_by: %s", widget.GroupBy)
	}

	rows, err := s.repo.ExecuteWidget(ctx, d.ObjectID, widget.Metric, widget.Field, widget.GroupBy)
	if err != nil {
		return nil, fmt.Errorf("execute widget: %w", err)
	}

	return &WidgetResult{WidgetID: widget.ID, Title: widget.Title, Rows: rows}, nil
}
