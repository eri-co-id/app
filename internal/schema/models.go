package schema

import "encoding/json"

type ObjectDefinition struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	PluralLabel string `json:"plural_label"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type FieldDefinition struct {
	ID       int64           `json:"id"`
	ObjectID int64           `json:"object_id"`
	Name     string          `json:"name"`
	Label    string          `json:"label"`
	Type     string          `json:"type"`
	Required bool            `json:"required"`
	Position int             `json:"position"`
	Metadata json.RawMessage `json:"metadata"`
}

type FieldMetadata struct {
	Options          []string `json:"options,omitempty"`
	RelationToObject string   `json:"relation_to_object,omitempty"`
	Readonly         bool     `json:"readonly,omitempty"`
	Calculated       bool     `json:"calculated,omitempty"`
	Default          string   `json:"default,omitempty"`
	Min              int      `json:"min,omitempty"`
	Max              int      `json:"max,omitempty"`
	MaxSizeMB        int      `json:"maxSizeMB,omitempty"`
}

func (f FieldDefinition) ParseMetadata() FieldMetadata {
	var m FieldMetadata
	if len(f.Metadata) > 0 {
		_ = json.Unmarshal(f.Metadata, &m)
	}
	return m
}

type ObjectWithFields struct {
	ObjectDefinition
	Fields []FieldDefinition `json:"fields"`
}

type ListResponse struct {
	Objects []ObjectDefinition `json:"objects"`
}
