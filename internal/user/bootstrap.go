package user

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// SeedSuperAdmin creates the first super_admin from SUPER_ADMIN_EMAIL /
// SUPER_ADMIN_NAME env vars, if no super_admin exists yet. Idempotent: safe
// to call on every boot. Without those env vars set, it's a no-op (with a
// warning if the instance has no super_admin at all yet) rather than a
// hardcoded identity baked into a migration file.
func SeedSuperAdmin(ctx context.Context, repo *Repository) error {
	count, err := repo.CountByRole(ctx, "super_admin")
	if err != nil {
		return fmt.Errorf("count super admins: %w", err)
	}
	if count > 0 {
		return nil
	}

	email := strings.TrimSpace(os.Getenv("SUPER_ADMIN_EMAIL"))
	name := strings.TrimSpace(os.Getenv("SUPER_ADMIN_NAME"))
	if email == "" || name == "" {
		log.Println("WARNING: no super_admin exists and SUPER_ADMIN_EMAIL/SUPER_ADMIN_NAME are not set — nobody can sign in yet")
		return nil
	}

	_, err = repo.Create(ctx, &CreateUserRequest{Email: email, Name: name, Role: "super_admin"}, nil)
	if err != nil {
		return fmt.Errorf("create bootstrap super admin: %w", err)
	}
	log.Printf("Bootstrapped super_admin: %s <%s>", name, email)
	return nil
}
