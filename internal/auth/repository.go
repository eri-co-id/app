package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, role, status, tokens_invalidated_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`, email).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Status,
		&user.TokensInvalidatedAt, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	return user, nil
}

// InvalidateTokens powers Logout: every token issued before now is rejected
// by ParseAccessToken/RefreshTokens, even though JWTs are otherwise stateless.
func (r *Repository) InvalidateTokens(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET tokens_invalidated_at = CURRENT_TIMESTAMP WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("invalidate tokens: %w", err)
	}
	return nil
}

func (r *Repository) CreateOTPSession(ctx context.Context, email, code string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO otp_sessions (email, otp_code, attempt_count, max_attempts, expires_at)
		VALUES ($1, $2, 0, $3, $4)
	`, email, code, MaxAttempts, expiresAt)
	if err != nil {
		return fmt.Errorf("create otp session: %w", err)
	}
	return nil
}

func (r *Repository) GetOTPSession(ctx context.Context, email string) (*OTPSession, error) {
	session := &OTPSession{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, otp_code, attempt_count, max_attempts, expires_at, created_at, updated_at
		FROM otp_sessions
		WHERE email = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, email).Scan(&session.ID, &session.Email, &session.OTPCode, &session.AttemptCount, &session.MaxAttempts, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query otp session: %w", err)
	}
	return session, nil
}

func (r *Repository) UpdateOTPSession(ctx context.Context, session *OTPSession) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE otp_sessions
		SET attempt_count = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, session.AttemptCount, session.ID)
	if err != nil {
		return fmt.Errorf("update otp session: %w", err)
	}
	return nil
}

func (r *Repository) DeleteOTPSession(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM otp_sessions WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete otp session: %w", err)
	}
	return nil
}
