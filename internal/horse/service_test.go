package horse

import "testing"

func TestValidateRequiresName(t *testing.T) {
	req := &UpsertHorseRequest{Name: "  "}
	if err := validate(req); err == nil {
		t.Fatal("expected error for blank name")
	}
}

func TestValidateDefaultsHealthStatus(t *testing.T) {
	req := &UpsertHorseRequest{Name: "Thunder"}
	if err := validate(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.HealthStatus != "active" {
		t.Fatalf("expected default health_status 'active', got %q", req.HealthStatus)
	}
}

func TestValidateRejectsUnknownHealthStatus(t *testing.T) {
	req := &UpsertHorseRequest{Name: "Thunder", HealthStatus: "sick"}
	if err := validate(req); err == nil {
		t.Fatal("expected error for invalid health_status")
	}
}

func TestValidateAcceptsKnownHealthStatuses(t *testing.T) {
	for status := range ValidHealthStatuses {
		req := &UpsertHorseRequest{Name: "Thunder", HealthStatus: status}
		if err := validate(req); err != nil {
			t.Fatalf("unexpected error for status %q: %v", status, err)
		}
	}
}
