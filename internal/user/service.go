package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	DefaultListLimit = 50
	MaxListLimit     = 200
)

var ErrNotFound = errors.New("user not found")

// Mailer sends the invite email after a user is created. See internal/mailer
// for the real SMTP implementation and the LogMailer dev-mode fallback.
type Mailer interface {
	SendInviteEmail(ctx context.Context, email, name string) error
}

type Service struct {
	repo   *Repository
	mailer Mailer
}

func NewService(repo *Repository, mailer Mailer) *Service {
	return &Service{repo: repo, mailer: mailer}
}

func validateCreate(req *CreateUserRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if !ValidRoles[req.Role] {
		return fmt.Errorf("invalid role: %s", req.Role)
	}
	return nil
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest, createdByID *int64) (*User, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	u, err := s.repo.Create(ctx, req, createdByID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := s.mailer.SendInviteEmail(ctx, u.Email, u.Name); err != nil {
		log.Printf("send invite email for %s: %v", u.Email, err)
	}

	return u, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if u == nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]User, error) {
	if limit <= 0 || limit > MaxListLimit {
		limit = DefaultListLimit
	}
	if offset < 0 {
		offset = 0
	}
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	u, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if u == nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (s *Service) UpdateRole(ctx context.Context, id int64, role string) (*User, error) {
	if !ValidRoles[role] {
		return nil, fmt.Errorf("invalid role: %s", role)
	}
	u, err := s.repo.UpdateRole(ctx, id, role)
	if err != nil {
		return nil, fmt.Errorf("update user role: %w", err)
	}
	if u == nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	found, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if !found {
		return ErrNotFound
	}
	return nil
}
