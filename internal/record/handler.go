package record

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/schema"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func writeRecordError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, schema.ErrNotFound):
		http.Error(w, "object not found", http.StatusNotFound)
	case errors.Is(err, ErrNotFound):
		http.Error(w, "record not found", http.StatusNotFound)
	case errors.Is(err, ErrForbidden):
		http.Error(w, "forbidden", http.StatusForbidden)
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	actor, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	records, err := h.service.List(r.Context(), objectName, actor, limit, offset)
	if err != nil {
		writeRecordError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListResponse{Records: records})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	actor, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpsertRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	rec, err := h.service.Create(r.Context(), objectName, &req, actor)
	if err != nil {
		writeRecordError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rec)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	actor, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	rec, err := h.service.Get(r.Context(), objectName, id, actor)
	if err != nil {
		writeRecordError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rec)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpsertRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	rec, err := h.service.Update(r.Context(), objectName, id, &req)
	if err != nil {
		writeRecordError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rec)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), objectName, id); err != nil {
		writeRecordError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
