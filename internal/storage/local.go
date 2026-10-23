package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage writes files under baseDir and serves them back at
// urlPrefix + key (main.go mounts an http.FileServer at urlPrefix for this).
type LocalStorage struct {
	baseDir   string
	urlPrefix string
}

func NewLocalStorage(baseDir, urlPrefix string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir, urlPrefix: urlPrefix}
}

func (s *LocalStorage) Save(_ context.Context, key string, data []byte) (string, error) {
	path := filepath.Join(s.baseDir, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return strings.TrimSuffix(s.urlPrefix, "/") + "/" + key, nil
}

func (s *LocalStorage) Delete(_ context.Context, url string) error {
	key := strings.TrimPrefix(url, strings.TrimSuffix(s.urlPrefix, "/")+"/")
	path := filepath.Join(s.baseDir, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove file: %w", err)
	}
	return nil
}
