package equipment

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

var ErrNotFound = errors.New("equipment not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func validate(req *UpsertEquipmentRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Type) == "" {
		return errors.New("type is required")
	}
	if req.Condition == "" {
		req.Condition = "good"
	}
	if !ValidConditions[req.Condition] {
		return fmt.Errorf("invalid condition: %s", req.Condition)
	}
	if req.Active == nil {
		active := true
		req.Active = &active
	}
	return nil
}

func (s *Service) Create(ctx context.Context, req *UpsertEquipmentRequest) (*Equipment, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	e, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create equipment: %w", err)
	}
	return e, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Equipment, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get equipment: %w", err)
	}
	if e == nil {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Equipment, error) {
	if limit <= 0 || limit > MaxListLimit {
		limit = DefaultListLimit
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list equipment: %w", err)
	}
	return items, nil
}

func (s *Service) ListByHorse(ctx context.Context, horseID int64) ([]Equipment, error) {
	items, err := s.repo.ListByHorse(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("list equipment by horse: %w", err)
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpsertEquipmentRequest) (*Equipment, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	e, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update equipment: %w", err)
	}
	if e == nil {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	found, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete equipment: %w", err)
	}
	if !found {
		return ErrNotFound
	}
	return nil
}
