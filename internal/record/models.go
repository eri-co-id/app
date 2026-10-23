package record

import (
	"encoding/json"
	"time"
)

type Record struct {
	ID            int64           `json:"id"`
	ObjectName    string          `json:"object_name"`
	Data          json.RawMessage `json:"data"`
	CreatedByID   int64           `json:"created_by_id"`
	CreatedByName string          `json:"created_by_name"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type UpsertRecordRequest struct {
	Data map[string]any `json:"data"`
}

type ListResponse struct {
	Records []Record `json:"records"`
}
