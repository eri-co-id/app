package lessoncode

import "time"

type LessonCode struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Type          string    `json:"type"`
	DefaultFee    *int64    `json:"default_fee,omitempty"`
	DurationHours *int      `json:"duration_hours,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpsertLessonCodeRequest struct {
	Code          string  `json:"code"`
	Type          string  `json:"type"`
	DefaultFee    *int64  `json:"default_fee"`
	DurationHours *int    `json:"duration_hours"`
	Description   *string `json:"description"`
	Active        *bool   `json:"active"`
}

type ListResponse struct {
	LessonCodes []LessonCode `json:"lesson_codes"`
}
