package file

import "time"

type File struct {
	ID           int64     `json:"id"`
	RecordID     int64     `json:"record_id"`
	FileName     string    `json:"file_name"`
	FileURL      string    `json:"file_url"`
	FileType     string    `json:"file_type"`
	FileSize     int64     `json:"file_size"`
	UploadedByID *int64    `json:"uploaded_by_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListResponse struct {
	Files []File `json:"files"`
}
