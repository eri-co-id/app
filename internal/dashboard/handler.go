package dashboard

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/schema"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func writeDashboardError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, schema.ErrNotFound), errors.Is(err, ErrNotFound), errors.Is(err, ErrWidgetNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	dashboards, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListResponse{Dashboards: dashboards})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Get(r.Context(), r.PathValue("name"))
	if err != nil {
		writeDashboardError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) ExecuteWidget(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	widgetID := r.PathValue("widgetId")

	if r.URL.Query().Get("format") == "csv" {
		actor, ok := auth.UserFromContext(r.Context())
		if !ok || (actor.Role != "admin" && actor.Role != "super_admin") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	result, err := h.service.ExecuteWidget(r.Context(), objectName, widgetID)
	if err != nil {
		writeDashboardError(w, err)
		return
	}

	if r.URL.Query().Get("format") == "csv" {
		writeCSV(w, result)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func writeCSV(w http.ResponseWriter, result *WidgetResult) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, result.WidgetID))

	cw := csv.NewWriter(w)
	defer cw.Flush()

	_ = cw.Write([]string{"group_key", "value"})
	for _, row := range result.Rows {
		key := ""
		if row.GroupKey != nil {
			key = *row.GroupKey
		}
		_ = cw.Write([]string{key, fmt.Sprintf("%v", row.Value)})
	}
}
