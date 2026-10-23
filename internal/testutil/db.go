package testutil

import (
	"context"
	"os"
	"testing"
	"time"

	"database/sql"

	"github.com/eri-co-id/app/internal/db"
)

// SetupDB connects to DATABASE_URL (falling back to the local dev default)
// and runs migrations, skipping the test if no database is reachable — so
// `go test ./...` stays safe to run anywhere, but is a real integration test
// wherever a DB is available (like this project's local dev Postgres).
func SetupDB(t *testing.T) *sql.DB {
	t.Helper()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://ericoid_user:ericoid_password@localhost:5433/ericoid_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, url)
	if err != nil {
		t.Skipf("skipping integration test: no reachable database: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}
