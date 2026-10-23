package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type TokenClaims struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

const (
	AccessTokenExpiry  = 24 * time.Hour
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")

// GenerateTokens issues a signed access/refresh token pair for the user.
// Tokens are HMAC-SHA256 signed (header.payload.signature, base64url, no padding) —
// a minimal JWT-compatible compact serialization, no external JWT library needed.
func GenerateTokens(user *User, accessSecret, refreshSecret string) (accessToken, refreshToken string, err error) {
	now := time.Now()

	accessToken, err = signToken(TokenClaims{
		UserID: user.ID, Email: user.Email, Role: user.Role,
		TokenType: "access", IssuedAt: now.Unix(), ExpiresAt: now.Add(AccessTokenExpiry).Unix(),
	}, accessSecret)
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err = signToken(TokenClaims{
		UserID: user.ID, Email: user.Email, Role: user.Role,
		TokenType: "refresh", IssuedAt: now.Unix(), ExpiresAt: now.Add(RefreshTokenExpiry).Unix(),
	}, refreshSecret)
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// VerifyToken checks the signature and expiry, and that the token is of wantType
// ("access" or "refresh"). Returns the decoded claims on success.
func VerifyToken(tokenString, secret, wantType string) (*TokenClaims, error) {
	payloadB64, sigB64, ok := splitToken(tokenString)
	if !ok {
		return nil, ErrInvalidToken
	}

	expectedSig := hmac.New(sha256.New, []byte(secret))
	expectedSig.Write([]byte(payloadB64))
	expected := base64.RawURLEncoding.EncodeToString(expectedSig.Sum(nil))

	gotSig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, ErrInvalidToken
	}
	expectedSigBytes, _ := base64.RawURLEncoding.DecodeString(expected)
	if !hmac.Equal(gotSig, expectedSigBytes) {
		return nil, ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims TokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != wantType {
		return nil, ErrInvalidToken
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}

func signToken(claims TokenClaims, secret string) (string, error) {
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payloadB64))
	sigB64 := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return payloadB64 + "." + sigB64, nil
}

func splitToken(token string) (payloadB64, sigB64 string, ok bool) {
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '.' {
			return token[:i], token[i+1:], true
		}
	}
	return "", "", false
}
