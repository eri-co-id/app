package horse

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

var ErrNotFound = errors.New("horse not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func validate(req *UpsertHorseRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if req.HealthStatus == "" {
		req.HealthStatus = "active"
	}
	if !ValidHealthStatuses[req.HealthStatus] {
		return fmt.Errorf("invalid health_status: %s", req.HealthStatus)
	}
	return nil
}

func (s *Service) Create(ctx context.Context, req *UpsertHorseRequest) (*Horse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	horse, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create horse: %w", err)
	}
	return horse, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Horse, error) {
	horse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get horse: %w", err)
	}
	if horse == nil {
		return nil, ErrNotFound
	}
	return horse, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Horse, error) {
	if limit <= 0 || limit > MaxListLimit {
		limit = DefaultListLimit
	}
	if offset < 0 {
		offset = 0
	}
	horses, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list horses: %w", err)
	}
	return horses, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpsertHorseRequest) (*Horse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	horse, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update horse: %w", err)
	}
	if horse == nil {
		return nil, ErrNotFound
	}
	return horse, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	found, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete horse: %w", err)
	}
	if !found {
		return ErrNotFound
	}
	return nil
}
