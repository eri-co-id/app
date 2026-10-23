package user

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const columns = `id, email, name, phone, role, status, created_by_id, created_at, updated_at`

func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.Role, &u.Status, &u.CreatedByID, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) Create(ctx context.Context, req *CreateUserRequest, createdByID *int64) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, phone, role, status, created_by_id)
		VALUES ($1, $2, $3, $4, 'active', $5)
		RETURNING `+columns,
		req.Email, req.Name, req.Phone, req.Role, createdByID,
	)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return u, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+columns+` FROM users WHERE id = $1 AND deleted_at IS NULL`, id)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	return u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+columns+` FROM users WHERE email = $1 AND deleted_at IS NULL`, email)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	return u, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+columns+` FROM users WHERE deleted_at IS NULL
		ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE users SET name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING `+columns,
		req.Name, req.Phone, id,
	)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

func (r *Repository) UpdateRole(ctx context.Context, id int64, role string) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE users SET role = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING `+columns,
		role, id,
	)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("update user role: %w", err)
	}
	return u, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false, fmt.Errorf("soft-delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return rows > 0, nil
}

func (r *Repository) CountByRole(ctx context.Context, role string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*) FROM users WHERE role = $1 AND deleted_at IS NULL`, role).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users by role: %w", err)
	}
	return count, nil
}
