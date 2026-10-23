package auth

import "time"

type User struct {
	ID                  int64      `json:"id"`
	Email               string     `json:"email"`
	Name                string     `json:"name"`
	Role                string     `json:"role"`
	Status              string     `json:"status"`
	TokensInvalidatedAt *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type OTPSession struct {
	ID           int64
	Email        string
	OTPCode      string
	AttemptCount int
	MaxAttempts  int
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SendOTPRequest struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Message            string `json:"message"`
	ExpiresInSeconds   int    `json:"expires_in_seconds"`
}

type VerifyOTPRequest struct {
	Email   string `json:"email"`
	OTPCode string `json:"otp_code"`
}

type VerifyOTPResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MeResponse struct {
	User User `json:"user"`
}
