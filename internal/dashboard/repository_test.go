package dashboard

import "testing"

func TestIdentifierPatternAcceptsValidFieldNames(t *testing.T) {
	for _, name := range []string{"fee", "code_lesson", "activity_type", "_private"} {
		if !identifierPattern.MatchString(name) {
			t.Fatalf("expected %q to be a valid identifier", name)
		}
	}
}

func TestIdentifierPatternRejectsInjectionAttempts(t *testing.T) {
	for _, name := range []string{
		"fee'; DROP TABLE records; --",
		"fee) UNION SELECT password FROM users --",
		"fee OR 1=1",
		"1fee",
		"",
		"fee bar",
	} {
		if identifierPattern.MatchString(name) {
			t.Fatalf("expected %q to be rejected as an invalid identifier", name)
		}
	}
}
