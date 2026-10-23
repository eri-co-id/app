package file

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/record"
	"github.com/eri-co-id/app/internal/schema"
)

const maxUploadBytes = MaxVideoBytes + (1 << 20) // cap the whole request slightly above the largest allowed file

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func writeFileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, schema.ErrNotFound), errors.Is(err, record.ErrNotFound), errors.Is(err, ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, record.ErrForbidden):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, ErrUnsupportedType), errors.Is(err, ErrTooLarge), errors.Is(err, ErrTooManyFiles):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	recordID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid record id", http.StatusBadRequest)
		return
	}

	actor, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		http.Error(w, "file too large or malformed upload", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	uploaded, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing 'file' field", http.StatusBadRequest)
		return
	}
	defer uploaded.Close()

	data, err := io.ReadAll(uploaded)
	if err != nil {
		http.Error(w, "failed to read upload", http.StatusBadRequest)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if detected := http.DetectContentType(data); detected != "application/octet-stream" {
		mimeType = detected
	}

	f, err := h.service.Upload(r.Context(), objectName, recordID, actor, header.Filename, mimeType, data)
	if err != nil {
		writeFileError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("name")
	recordID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid record id", http.StatusBadRequest)
		return
	}

	actor, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	files, err := h.service.List(r.Context(), objectName, recordID, actor)
	if err != nil {
		writeFileError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListResponse{Files: files})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeFileError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
