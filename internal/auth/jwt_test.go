package auth

import "testing"

func TestGenerateAndVerifyTokens(t *testing.T) {
	user := &User{ID: 42, Email: "trainer@example.com", Role: "trainer"}
	accessToken, refreshToken, err := GenerateTokens(user, "access-secret", "refresh-secret")
	if err != nil {
		t.Fatalf("generate tokens: %v", err)
	}

	claims, err := VerifyToken(accessToken, "access-secret", "access")
	if err != nil {
		t.Fatalf("verify access token: %v", err)
	}
	if claims.UserID != user.ID || claims.Email != user.Email || claims.Role != user.Role {
		t.Fatalf("unexpected claims: %+v", claims)
	}

	if _, err := VerifyToken(refreshToken, "refresh-secret", "refresh"); err != nil {
		t.Fatalf("verify refresh token: %v", err)
	}
}

func TestVerifyTokenRejectsWrongType(t *testing.T) {
	user := &User{ID: 1, Email: "a@b.com", Role: "admin"}
	accessToken, _, _ := GenerateTokens(user, "s1", "s2")

	if _, err := VerifyToken(accessToken, "s1", "refresh"); err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for wrong type, got %v", err)
	}
}

func TestVerifyTokenRejectsTamperedSignature(t *testing.T) {
	user := &User{ID: 1, Email: "a@b.com", Role: "admin"}
	accessToken, _, _ := GenerateTokens(user, "s1", "s2")

	if _, err := VerifyToken(accessToken, "wrong-secret", "access"); err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for bad secret, got %v", err)
	}
}

func TestVerifyTokenRejectsExpired(t *testing.T) {
	user := &User{ID: 1, Email: "a@b.com", Role: "admin"}
	token, err := signToken(TokenClaims{
		UserID: user.ID, Email: user.Email, Role: user.Role,
		TokenType: "access", IssuedAt: 0, ExpiresAt: 1,
	}, "secret")
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	if _, err := VerifyToken(token, "secret", "access"); err != ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestVerifyTokenRejectsMalformed(t *testing.T) {
	if _, err := VerifyToken("not-a-valid-token", "secret", "access"); err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for malformed token, got %v", err)
	}
}
