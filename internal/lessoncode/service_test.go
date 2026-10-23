package lessoncode

import "testing"

func TestValidateRequiresCodeAndType(t *testing.T) {
	if err := validate(&UpsertLessonCodeRequest{Type: "Dressage"}); err == nil {
		t.Fatal("expected error for blank code")
	}
	if err := validate(&UpsertLessonCodeRequest{Code: "DR-01"}); err == nil {
		t.Fatal("expected error for blank type")
	}
}

func TestValidateDefaultsActive(t *testing.T) {
	req := &UpsertLessonCodeRequest{Code: "DR-01", Type: "Dressage"}
	if err := validate(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Active == nil || !*req.Active {
		t.Fatal("expected default active=true")
	}
}
