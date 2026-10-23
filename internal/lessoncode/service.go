package lessoncode

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultListLimit = 50
	MaxListLimit     = 200
)

var ErrNotFound = errors.New("lesson code not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func validate(req *UpsertLessonCodeRequest) error {
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("code is required")
	}
	if strings.TrimSpace(req.Type) == "" {
		return errors.New("type is required")
	}
	if req.Active == nil {
		active := true
		req.Active = &active
	}
	return nil
}

func (s *Service) Create(ctx context.Context, req *UpsertLessonCodeRequest) (*LessonCode, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	lc, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create lesson code: %w", err)
	}
	return lc, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*LessonCode, error) {
	lc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get lesson code: %w", err)
	}
	if lc == nil {
		return nil, ErrNotFound
	}
	return lc, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]LessonCode, error) {
	if limit <= 0 || limit > MaxListLimit {
		limit = DefaultListLimit
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list lesson codes: %w", err)
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpsertLessonCodeRequest) (*LessonCode, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	lc, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update lesson code: %w", err)
	}
	if lc == nil {
		return nil, ErrNotFound
	}
	return lc, nil
}

func (s *Service) Deactivate(ctx context.Context, id int64) error {
	found, err := s.repo.Deactivate(ctx, id)
	if err != nil {
		return fmt.Errorf("deactivate lesson code: %w", err)
	}
	if !found {
		return ErrNotFound
	}
	return nil
}
