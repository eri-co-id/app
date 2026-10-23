package record

import (
	"context"
	"errors"
	"fmt"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/schema"
)

const (
	DefaultListLimit = 50
	MaxListLimit     = 200
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrForbidden = errors.New("not allowed to access this record")
)

func isAdmin(role string) bool {
	return role == "admin" || role == "super_admin"
}

type Service struct {
	repo      *Repository
	schemaSvc *schema.Service
}

func NewService(repo *Repository, schemaSvc *schema.Service) *Service {
	return &Service{repo: repo, schemaSvc: schemaSvc}
}

func (s *Service) Create(ctx context.Context, objectName string, req *UpsertRecordRequest, actor *auth.User) (*Record, error) {
	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return nil, err
	}

	validated, err := schema.ValidateData(object.Fields, req.Data)
	if err != nil {
		return nil, err
	}

	rec, err := s.repo.Create(ctx, object.ID, validated, actor.ID)
	if err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}
	rec.ObjectName = objectName
	return rec, nil
}

// Get enforces record-level visibility: non-admin roles may only fetch their
// own records (CONSTRAINTS.md: trainer cannot fetch another trainer's record).
func (s *Service) Get(ctx context.Context, objectName string, id int64, actor *auth.User) (*Record, error) {
	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return nil, err
	}

	rec, err := s.repo.GetByID(ctx, object.ID, id)
	if err != nil {
		return nil, fmt.Errorf("get record: %w", err)
	}
	if rec == nil {
		return nil, ErrNotFound
	}
	if !isAdmin(actor.Role) && rec.CreatedByID != actor.ID {
		return nil, ErrForbidden
	}
	rec.ObjectName = objectName
	return rec, nil
}

// List scopes non-admin roles to their own records; admin/super_admin see all.
func (s *Service) List(ctx context.Context, objectName string, actor *auth.User, limit, offset int) ([]Record, error) {
	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > MaxListLimit {
		limit = DefaultListLimit
	}
	if offset < 0 {
		offset = 0
	}

	var createdByFilter *int64
	if !isAdmin(actor.Role) {
		createdByFilter = &actor.ID
	}

	records, err := s.repo.List(ctx, object.ID, createdByFilter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list records: %w", err)
	}
	for i := range records {
		records[i].ObjectName = objectName
	}
	return records, nil
}

// Update and Delete are admin/super_admin-only retroactive edits (SPEC.md §7);
// route-level role gating handles that, so no ownership check here.

func (s *Service) Update(ctx context.Context, objectName string, id int64, req *UpsertRecordRequest) (*Record, error) {
	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return nil, err
	}

	validated, err := schema.ValidateData(object.Fields, req.Data)
	if err != nil {
		return nil, err
	}

	rec, err := s.repo.Update(ctx, object.ID, id, validated)
	if err != nil {
		return nil, fmt.Errorf("update record: %w", err)
	}
	if rec == nil {
		return nil, ErrNotFound
	}
	rec.ObjectName = objectName
	return rec, nil
}

func (s *Service) Delete(ctx context.Context, objectName string, id int64) error {
	object, err := s.schemaSvc.GetWithFields(ctx, objectName)
	if err != nil {
		return err
	}

	found, err := s.repo.SoftDelete(ctx, object.ID, id)
	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	if !found {
		return ErrNotFound
	}
	return nil
}
