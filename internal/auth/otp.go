package auth

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	OTPLength    = 6
	OTPExpiry    = 15 * time.Minute
	MaxAttempts  = 5
)

// GenerateOTP creates a cryptographically random 6-digit code
func GenerateOTP() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("crypto rand: %w", err)
	}
	n := binary.BigEndian.Uint32(b) % 1000000
	return fmt.Sprintf("%06d", n), nil
}
