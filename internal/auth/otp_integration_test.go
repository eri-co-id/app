package auth

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/eri-co-id/app/internal/testutil"
)

// captureMailer records the OTP instead of sending it, so the test can read
// it back the same way a real user would from their inbox.
type captureMailer struct{ code string }

func (m *captureMailer) SendOTP(_ context.Context, _, code string) error {
	m.code = code
	return nil
}

func seedTestUser(t *testing.T, database *sql.DB, prefix string) (email string, cleanup func()) {
	t.Helper()
	email = fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano())
	_, err := database.ExecContext(context.Background(), `
		INSERT INTO users (email, name, role, status) VALUES ($1, 'Integration Test', 'trainer', 'active')`, email)
	if err != nil {
		t.Fatalf("seed test user: %v", err)
	}
	return email, func() {
		database.ExecContext(context.Background(), `DELETE FROM users WHERE email = $1`, email)
	}
}

func TestSendAndVerifyOTPIntegration(t *testing.T) {
	database := testutil.SetupDB(t)
	ctx := context.Background()

	email, cleanup := seedTestUser(t, database, "otp-happy")
	defer cleanup()

	mailer := &captureMailer{}
	svc := NewService(NewRepository(database), "test-access-secret-aaaaaaaaaaaaaaaaaaaaaaaa", "test-refresh-secret-aaaaaaaaaaaaaaaaaaaaaaaa", mailer)

	if err := svc.SendOTP(ctx, email); err != nil {
		t.Fatalf("send otp: %v", err)
	}
	if mailer.code == "" {
		t.Fatal("expected mailer to capture a code")
	}

	user, err := svc.VerifyOTP(ctx, email, mailer.code)
	if err != nil {
		t.Fatalf("verify otp with correct code: %v", err)
	}
	if user.Email != email {
		t.Fatalf("expected verified user %q, got %q", email, user.Email)
	}

	// OTP is single-use — the same code must not verify twice.
	if _, err := svc.VerifyOTP(ctx, email, mailer.code); err == nil {
		t.Fatal("expected error reusing an already-verified OTP code")
	}
}

func TestVerifyOTPLocksOutAfterFiveWrongAttemptsIntegration(t *testing.T) {
	database := testutil.SetupDB(t)
	ctx := context.Background()

	email, cleanup := seedTestUser(t, database, "otp-lockout")
	defer cleanup()

	mailer := &captureMailer{}
	svc := NewService(NewRepository(database), "test-access-secret-aaaaaaaaaaaaaaaaaaaaaaaa", "test-refresh-secret-aaaaaaaaaaaaaaaaaaaaaaaa", mailer)

	if err := svc.SendOTP(ctx, email); err != nil {
		t.Fatalf("send otp: %v", err)
	}

	for i := 0; i < 5; i++ {
		if _, err := svc.VerifyOTP(ctx, email, "000000"); err == nil {
			t.Fatalf("attempt %d: expected error for wrong code", i+1)
		}
	}

	// 6th attempt, even with the correct code, must be locked out.
	if _, err := svc.VerifyOTP(ctx, email, mailer.code); err == nil {
		t.Fatal("expected lockout after 5 failed attempts")
	}
}

func TestVerifyOTPRejectsExpiredIntegration(t *testing.T) {
	database := testutil.SetupDB(t)
	ctx := context.Background()

	email, cleanup := seedTestUser(t, database, "otp-expired")
	defer cleanup()

	mailer := &captureMailer{}
	svc := NewService(NewRepository(database), "test-access-secret-aaaaaaaaaaaaaaaaaaaaaaaa", "test-refresh-secret-aaaaaaaaaaaaaaaaaaaaaaaa", mailer)

	if err := svc.SendOTP(ctx, email); err != nil {
		t.Fatalf("send otp: %v", err)
	}

	// Force the session into the past instead of waiting out the real 15-minute expiry.
	_, err := database.ExecContext(ctx, `
		UPDATE otp_sessions SET expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute' WHERE email = $1`, email)
	if err != nil {
		t.Fatalf("force-expire otp session: %v", err)
	}

	if _, err := svc.VerifyOTP(ctx, email, mailer.code); err == nil {
		t.Fatal("expected error verifying an expired OTP")
	}
}
