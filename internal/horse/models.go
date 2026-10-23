package horse

import "time"

var ValidHealthStatuses = map[string]bool{
	"active":  true,
	"injured": true,
	"retired": true,
}

type Horse struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Code             *string    `json:"code,omitempty"`
	Breed            *string    `json:"breed,omitempty"`
	DOB              *time.Time `json:"dob,omitempty"`
	Color            *string    `json:"color,omitempty"`
	OwnerName        *string    `json:"owner_name,omitempty"`
	CurrentTrainerID *int64     `json:"current_trainer_id,omitempty"`
	HealthStatus     string     `json:"health_status"`
	Notes            *string    `json:"notes,omitempty"`
	PhotoURL         *string    `json:"photo_url,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type UpsertHorseRequest struct {
	Name             string     `json:"name"`
	Code             *string    `json:"code"`
	Breed            *string    `json:"breed"`
	DOB              *time.Time `json:"dob"`
	Color            *string    `json:"color"`
	OwnerName        *string    `json:"owner_name"`
	CurrentTrainerID *int64     `json:"current_trainer_id"`
	HealthStatus     string     `json:"health_status"`
	Notes            *string    `json:"notes"`
	PhotoURL         *string    `json:"photo_url"`
}

type ListResponse struct {
	Horses []Horse `json:"horses"`
}
