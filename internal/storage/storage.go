package storage

import "context"

// Storage persists uploaded file bytes and returns a URL to fetch them.
// LocalStorage is the dev-mode implementation (no Garage/S3 credentials in
// this environment); swap in an S3-compatible client for Garage in production.
type Storage interface {
	Save(ctx context.Context, key string, data []byte) (url string, err error)
	Delete(ctx context.Context, url string) error
}
