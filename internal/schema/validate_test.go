package schema

import (
	"encoding/json"
	"testing"
)

func field(name, fieldType string, required bool, metadata string) FieldDefinition {
	return FieldDefinition{Name: name, Type: fieldType, Required: required, Metadata: json.RawMessage(metadata)}
}

func TestValidateDataRequiresRequiredFields(t *testing.T) {
	fields := []FieldDefinition{field("horse", "relation", true, "{}")}
	_, err := ValidateData(fields, map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestValidateDataAppliesDefault(t *testing.T) {
	fields := []FieldDefinition{field("status", "select", true, `{"options":["completed","in_progress"],"default":"completed"}`)}
	result, err := ValidateData(fields, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "completed" {
		t.Fatalf("expected default applied, got %v", result["status"])
	}
}

func TestValidateDataSkipsReadonlyCalculated(t *testing.T) {
	fields := []FieldDefinition{field("current_stock", "number", false, `{"readonly":true,"calculated":true}`)}
	result, err := ValidateData(fields, map[string]any{"current_stock": 999.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, present := result["current_stock"]; present {
		t.Fatal("expected calculated field to be stripped from client-supplied data")
	}
}

func TestValidateDataSkipsFileFields(t *testing.T) {
	fields := []FieldDefinition{field("photos", "file", true, `{}`)}
	result, err := ValidateData(fields, map[string]any{})
	if err != nil {
		t.Fatalf("file fields should never block JSON validation, got: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected no fields in result, got %v", result)
	}
}

func TestValidateDataRejectsWrongType(t *testing.T) {
	fields := []FieldDefinition{field("duration_hours", "number", false, "{}")}
	_, err := ValidateData(fields, map[string]any{"duration_hours": "two"})
	if err == nil {
		t.Fatal("expected error for non-numeric value in number field")
	}
}

func TestValidateDataRejectsInvalidDate(t *testing.T) {
	fields := []FieldDefinition{field("date", "date", true, "{}")}
	_, err := ValidateData(fields, map[string]any{"date": "not-a-date"})
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestValidateDataRejectsSelectOutsideOptions(t *testing.T) {
	fields := []FieldDefinition{field("area", "select", true, `{"options":["stable","arena"]}`)}
	_, err := ValidateData(fields, map[string]any{"area": "kitchen"})
	if err == nil {
		t.Fatal("expected error for value outside select options")
	}
}

func TestValidateDataAllowsSelectWithEmptyOptions(t *testing.T) {
	fields := []FieldDefinition{field("grass_type", "select", true, `{"options":[]}`)}
	result, err := ValidateData(fields, map[string]any{"grass_type": "anything"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["grass_type"] != "anything" {
		t.Fatalf("expected value passed through, got %v", result["grass_type"])
	}
}

func TestValidateDataAcceptsValidMultiselect(t *testing.T) {
	fields := []FieldDefinition{field("equipment_used", "multiselect", false, `{"options":["saddle","bridle"]}`)}
	result, err := ValidateData(fields, map[string]any{"equipment_used": []any{"saddle"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := result["equipment_used"].([]string)
	if len(got) != 1 || got[0] != "saddle" {
		t.Fatalf("unexpected multiselect result: %v", got)
	}
}

func TestValidateDataRejectsMultiselectOutsideOptions(t *testing.T) {
	fields := []FieldDefinition{field("equipment_used", "multiselect", false, `{"options":["saddle","bridle"]}`)}
	_, err := ValidateData(fields, map[string]any{"equipment_used": []any{"spurs"}})
	if err == nil {
		t.Fatal("expected error for multiselect value outside options")
	}
}

func TestValidateDataOptionalFieldOmittedIsFine(t *testing.T) {
	fields := []FieldDefinition{field("notes", "text", false, "{}")}
	result, err := ValidateData(fields, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected no fields in result, got %v", result)
	}
}
