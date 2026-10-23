package record

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/schema"
	"github.com/eri-co-id/app/internal/testutil"
)

func seedTestTrainer(t *testing.T, database *sql.DB, prefix string) (id int64, cleanup func()) {
	t.Helper()
	email := fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano())
	err := database.QueryRowContext(context.Background(), `
		INSERT INTO users (email, name, role, status) VALUES ($1, 'Integration Test', 'trainer', 'active')
		RETURNING id`, email).Scan(&id)
	if err != nil {
		t.Fatalf("seed trainer: %v", err)
	}
	return id, func() {
		database.ExecContext(context.Background(), `DELETE FROM records WHERE created_by_id = $1`, id)
		database.ExecContext(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	}
}

// TestRecordOwnershipIsolationIntegration is the CONSTRAINTS.md floor check:
// a trainer cannot fetch or list another trainer's records; admin sees all.
func TestRecordOwnershipIsolationIntegration(t *testing.T) {
	database := testutil.SetupDB(t)
	ctx := context.Background()

	idA, cleanupA := seedTestTrainer(t, database, "rec-itest-a")
	defer cleanupA()
	idB, cleanupB := seedTestTrainer(t, database, "rec-itest-b")
	defer cleanupB()

	schemaSvc := schema.NewService(schema.NewRepository(database))
	recordSvc := NewService(NewRepository(database), schemaSvc)

	trainerA := &auth.User{ID: idA, Role: "trainer"}
	trainerB := &auth.User{ID: idB, Role: "trainer"}
	adminViewer := &auth.User{ID: idA, Role: "admin"} // role is all the service reads; reusing idA is fine

	rec, err := recordSvc.Create(ctx, "riding", &UpsertRecordRequest{Data: map[string]any{
		"date": "2026-10-23", "rider": float64(idA), "horse": float64(1), "code_lesson": float64(1),
	}}, trainerA)
	if err != nil {
		t.Fatalf("create record as trainer A: %v", err)
	}

	if _, err := recordSvc.Get(ctx, "riding", rec.ID, trainerB); err != ErrForbidden {
		t.Fatalf("expected ErrForbidden for trainer B fetching trainer A's record, got %v", err)
	}
	if _, err := recordSvc.Get(ctx, "riding", rec.ID, trainerA); err != nil {
		t.Fatalf("expected trainer A to fetch their own record, got %v", err)
	}
	if _, err := recordSvc.Get(ctx, "riding", rec.ID, adminViewer); err != nil {
		t.Fatalf("expected admin to fetch any record, got %v", err)
	}

	listA, err := recordSvc.List(ctx, "riding", trainerA, 50, 0)
	if err != nil {
		t.Fatalf("list as trainer A: %v", err)
	}
	if !containsRecordID(listA, rec.ID) {
		t.Fatal("expected trainer A's list to include their own record")
	}

	listB, err := recordSvc.List(ctx, "riding", trainerB, 50, 0)
	if err != nil {
		t.Fatalf("list as trainer B: %v", err)
	}
	if containsRecordID(listB, rec.ID) {
		t.Fatal("trainer B's list must not include trainer A's record")
	}

	listAdmin, err := recordSvc.List(ctx, "riding", adminViewer, 50, 0)
	if err != nil {
		t.Fatalf("list as admin: %v", err)
	}
	if !containsRecordID(listAdmin, rec.ID) {
		t.Fatal("expected admin's list to include every trainer's record")
	}
}

func containsRecordID(records []Record, id int64) bool {
	for _, r := range records {
		if r.ID == id {
			return true
		}
	}
	return false
}
