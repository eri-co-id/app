package user

import "testing"

func TestValidRolesCoversAllSixRoles(t *testing.T) {
	want := []string{"super_admin", "admin", "trainer", "groom", "arit", "cleaner"}
	if len(ValidRoles) != len(want) {
		t.Fatalf("expected %d valid roles, got %d", len(want), len(ValidRoles))
	}
	for _, role := range want {
		if !ValidRoles[role] {
			t.Fatalf("expected %q to be a valid role", role)
		}
	}
}

func TestValidRolesRejectsUnknown(t *testing.T) {
	if ValidRoles["operator"] {
		t.Fatal("'operator' should not be a valid role")
	}
}

func TestValidateCreateRequiresEmailNameRole(t *testing.T) {
	cases := []*CreateUserRequest{
		{Name: "Trainer", Role: "trainer"},
		{Email: "t@example.com", Role: "trainer"},
		{Email: "t@example.com", Name: "Trainer", Role: "operator"},
	}
	for _, req := range cases {
		if err := validateCreate(req); err == nil {
			t.Fatalf("expected error for %+v", req)
		}
	}
}

func TestValidateCreateTrimsWhitespace(t *testing.T) {
	req := &CreateUserRequest{Email: "  t@example.com  ", Name: "  Trainer  ", Role: "trainer"}
	if err := validateCreate(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Email != "t@example.com" || req.Name != "Trainer" {
		t.Fatalf("expected trimmed fields, got email=%q name=%q", req.Email, req.Name)
	}
}
