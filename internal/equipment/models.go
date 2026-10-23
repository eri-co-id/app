package equipment

import "time"

var ValidConditions = map[string]bool{
	"new":     true,
	"good":    true,
	"worn":    true,
	"damaged": true,
}

type Equipment struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	Code           *string    `json:"code,omitempty"`
	SerialBatch    *string    `json:"serial_batch,omitempty"`
	Condition      string     `json:"condition"`
	CurrentHorseID *int64     `json:"current_horse_id,omitempty"`
	Location       *string    `json:"location,omitempty"`
	PurchaseDate   *time.Time `json:"purchase_date,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	PhotoURL       *string    `json:"photo_url,omitempty"`
	Active         bool       `json:"active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UpsertEquipmentRequest struct {
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	Code           *string    `json:"code"`
	SerialBatch    *string    `json:"serial_batch"`
	Condition      string     `json:"condition"`
	CurrentHorseID *int64     `json:"current_horse_id"`
	Location       *string    `json:"location"`
	PurchaseDate   *time.Time `json:"purchase_date"`
	Notes          *string    `json:"notes"`
	PhotoURL       *string    `json:"photo_url"`
	Active         *bool      `json:"active"`
}

type ListResponse struct {
	Equipment []Equipment `json:"equipment"`
}
