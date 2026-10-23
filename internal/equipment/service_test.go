package equipment

import "testing"

func TestValidateRequiresNameAndType(t *testing.T) {
	if err := validate(&UpsertEquipmentRequest{Type: "saddle"}); err == nil {
		t.Fatal("expected error for blank name")
	}
	if err := validate(&UpsertEquipmentRequest{Name: "Saddle #1"}); err == nil {
		t.Fatal("expected error for blank type")
	}
}

func TestValidateDefaultsConditionAndActive(t *testing.T) {
	req := &UpsertEquipmentRequest{Name: "Saddle #1", Type: "saddle"}
	if err := validate(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Condition != "good" {
		t.Fatalf("expected default condition 'good', got %q", req.Condition)
	}
	if req.Active == nil || !*req.Active {
		t.Fatal("expected default active=true")
	}
}

func TestValidateRejectsUnknownCondition(t *testing.T) {
	req := &UpsertEquipmentRequest{Name: "Saddle #1", Type: "saddle", Condition: "pristine"}
	if err := validate(req); err == nil {
		t.Fatal("expected error for invalid condition")
	}
}
