package file

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/record"
	"github.com/eri-co-id/app/internal/schema"
	"github.com/eri-co-id/app/internal/storage"
)

const (
	MaxPhotoBytes   = 10 * 1024 * 1024
	MaxVideoBytes   = 50 * 1024 * 1024
	DefaultMaxFiles = 5
)

var (
	ErrNotFound        = errors.New("file not found")
	ErrUnsupportedType = errors.New("unsupported file type: must be image/* or video/*")
	ErrTooLarge        = errors.New("file exceeds the maximum allowed size")
	ErrTooManyFiles    = errors.New("record already has the maximum number of files")
)

type Service struct {
	repo      *Repository
	recordSvc *record.Service
	schemaSvc *schema.Service
	storage   storage.Storage
}

func NewService(repo *Repository, recordSvc *record.Service, schemaSvc *schema.Service, store storage.Storage) *Service {
	return &Service{repo: repo, recordSvc: recordSvc, schemaSvc: schemaSvc, storage: store}
}

// Upload validates and stores a file against a record. Reuses record.Service.Get
// for existence + ownership (non-admin can only attach files to their own records).
func (s *Service) Upload(ctx context.Context, objectName string, recordID int64, actor *auth.User, fileName, mimeType string, data []byte) (*File, error) {
	if _, err := s.recordSvc.Get(ctx, objectName, recordID, actor); err != nil {
		return nil, err
	}

	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return nil, err
	}

	maxFiles := DefaultMaxFiles
	maxPhotoBytes := MaxPhotoBytes
	for _, f := range object.Fields {
		if f.Type != "file" {
			continue
		}
		meta := f.ParseMetadata()
		if meta.Max > 0 {
			maxFiles = meta.Max
		}
		if meta.MaxSizeMB > 0 {
			maxPhotoBytes = meta.MaxSizeMB * 1024 * 1024
		}
	}

	var fileType string
	var maxBytes int
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		fileType, maxBytes = "photo", maxPhotoBytes
	case strings.HasPrefix(mimeType, "video/"):
		fileType, maxBytes = "video", MaxVideoBytes
	default:
		return nil, ErrUnsupportedType
	}
	if len(data) > maxBytes {
		return nil, ErrTooLarge
	}

	count, err := s.repo.CountByRecordID(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("count files: %w", err)
	}
	if count >= maxFiles {
		return nil, ErrTooManyFiles
	}

	key, err := uploadKey(objectName, recordID, fileName)
	if err != nil {
		return nil, err
	}
	url, err := s.storage.Save(ctx, key, data)
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	f, err := s.repo.Create(ctx, recordID, fileName, url, fileType, int64(len(data)), &actor.ID)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	return f, nil
}

func (s *Service) List(ctx context.Context, objectName string, recordID int64, actor *auth.User) ([]File, error) {
	if _, err := s.recordSvc.Get(ctx, objectName, recordID, actor); err != nil {
		return nil, err
	}
	files, err := s.repo.ListByRecordID(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}
	return files, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	f, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}
	if f == nil {
		return ErrNotFound
	}
	if err := s.storage.Delete(ctx, f.FileURL); err != nil {
		return fmt.Errorf("delete stored file: %w", err)
	}
	return nil
}

// uploadKey builds the storage key. fileName is client-controlled, so it's
// reduced to its base name (strips any path separators / ".." traversal)
// before being used in a filesystem path.
func uploadKey(objectName string, recordID int64, fileName string) (string, error) {
	safeName := filepath.Base(fileName)
	if safeName == "." || safeName == "/" || safeName == "" {
		safeName = "upload"
	}

	suffix := make([]byte, 8)
	if _, err := rand.Read(suffix); err != nil {
		return "", fmt.Errorf("generate upload key: %w", err)
	}
	return fmt.Sprintf("%s/%d/%s-%s", objectName, recordID, hex.EncodeToString(suffix), safeName), nil
}
