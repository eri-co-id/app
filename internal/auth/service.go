package auth

import (
	"context"
	"fmt"
	"time"
)

// Mailer delivers the OTP code to the user. See internal/mailer for the real
// SMTP implementation and the LogMailer dev-mode fallback.
type Mailer interface {
	SendOTP(ctx context.Context, email, code string) error
}

type Service struct {
	repo          *Repository
	accessSecret  string
	refreshSecret string
	mailer        Mailer
}

func NewService(repo *Repository, accessSecret, refreshSecret string, mailer Mailer) *Service {
	return &Service{repo: repo, accessSecret: accessSecret, refreshSecret: refreshSecret, mailer: mailer}
}

func (s *Service) SendOTP(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if user == nil || user.Status != "active" {
		// Don't leak whether email exists
		return nil
	}

	code, err := GenerateOTP()
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	expiresAt := time.Now().Add(OTPExpiry)
	if err := s.repo.CreateOTPSession(ctx, email, code, expiresAt); err != nil {
		return fmt.Errorf("create otp session: %w", err)
	}

	if err := s.mailer.SendOTP(ctx, email, code); err != nil {
		return fmt.Errorf("send otp email: %w", err)
	}
	return nil
}

func (s *Service) VerifyOTP(ctx context.Context, email, code string) (*User, error) {
	session, err := s.repo.GetOTPSession(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get otp session: %w", err)
	}

	if session == nil || session.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("otp expired or not found")
	}

	if session.AttemptCount >= session.MaxAttempts {
		return nil, fmt.Errorf("too many attempts")
	}

	if session.OTPCode != code {
		session.AttemptCount++
		_ = s.repo.UpdateOTPSession(ctx, session)
		return nil, fmt.Errorf("invalid otp")
	}

	// OTP valid, get user
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Delete OTP session (single-use)
	_ = s.repo.DeleteOTPSession(ctx, session.ID)

	return user, nil
}

// IssueTokens signs a fresh access/refresh token pair for the user.
func (s *Service) IssueTokens(user *User) (accessToken, refreshToken string, err error) {
	return GenerateTokens(user, s.accessSecret, s.refreshSecret)
}

// RefreshTokens validates the refresh token and issues a new access/refresh pair.
// The user is re-fetched so a deactivated account can't refresh its way back in.
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := VerifyToken(refreshToken, s.refreshSecret, "refresh")
	if err != nil {
		return "", "", fmt.Errorf("verify refresh token: %w", err)
	}

	user, err := s.repo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}
	if user == nil || user.Status != "active" || tokenRevoked(user, claims.IssuedAt) {
		return "", "", fmt.Errorf("user not found, inactive, or logged out")
	}

	return s.IssueTokens(user)
}

// ParseAccessToken validates an access token and re-fetches the current user,
// so a role change, deactivation, or logout takes effect before the token's
// natural expiry.
func (s *Service) ParseAccessToken(ctx context.Context, accessToken string) (*User, error) {
	claims, err := VerifyToken(accessToken, s.accessSecret, "access")
	if err != nil {
		return nil, fmt.Errorf("verify access token: %w", err)
	}

	user, err := s.repo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil || user.Status != "active" || tokenRevoked(user, claims.IssuedAt) {
		return nil, fmt.Errorf("user not found, inactive, or logged out")
	}

	return user, nil
}

// Logout invalidates every token issued for this user up to now (both access
// and refresh) — tokens are otherwise stateless JWTs that would keep verifying.
func (s *Service) Logout(ctx context.Context, userID int64) error {
	if err := s.repo.InvalidateTokens(ctx, userID); err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	return nil
}

func tokenRevoked(user *User, issuedAt int64) bool {
	return user.TokensInvalidatedAt != nil && issuedAt <= user.TokensInvalidatedAt.Unix()
}
