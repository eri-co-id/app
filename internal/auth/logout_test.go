package auth

import (
	"testing"
	"time"
)

func TestTokenRevokedNilNeverRevoked(t *testing.T) {
	user := &User{TokensInvalidatedAt: nil}
	if tokenRevoked(user, time.Now().Unix()) {
		t.Fatal("expected no revocation when tokens_invalidated_at is nil")
	}
}

func TestTokenRevokedRejectsOlderIssuedAt(t *testing.T) {
	cutoff := time.Now()
	user := &User{TokensInvalidatedAt: &cutoff}
	issuedBefore := cutoff.Add(-1 * time.Minute).Unix()
	if !tokenRevoked(user, issuedBefore) {
		t.Fatal("expected token issued before logout to be revoked")
	}
}

func TestTokenRevokedAllowsNewerIssuedAt(t *testing.T) {
	cutoff := time.Now()
	user := &User{TokensInvalidatedAt: &cutoff}
	issuedAfter := cutoff.Add(1 * time.Minute).Unix()
	if tokenRevoked(user, issuedAfter) {
		t.Fatal("expected token issued after logout to remain valid")
	}
}
