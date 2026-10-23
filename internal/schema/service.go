package schema

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("object not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]ObjectDefinition, error) {
	objects, err := s.repo.ListObjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("list objects: %w", err)
	}
	return objects, nil
}

func (s *Service) GetWithFields(ctx context.Context, name string) (*ObjectWithFields, error) {
	object, err := s.repo.GetObjectByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	if object == nil {
		return nil, ErrNotFound
	}

	fields, err := s.repo.ListFieldsByObjectID(ctx, object.ID)
	if err != nil {
		return nil, fmt.Errorf("list fields: %w", err)
	}

	return &ObjectWithFields{ObjectDefinition: *object, Fields: fields}, nil
}
