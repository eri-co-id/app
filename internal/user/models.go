package user

import "time"

var ValidRoles = map[string]bool{
	"super_admin": true,
	"admin":       true,
	"trainer":     true,
	"groom":       true,
	"arit":        true,
	"cleaner":     true,
}

type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Phone       *string   `json:"phone,omitempty"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	CreatedByID *int64    `json:"created_by_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email string  `json:"email"`
	Name  string  `json:"name"`
	Phone *string `json:"phone"`
	Role  string  `json:"role"`
}

type UpdateUserRequest struct {
	Name  string  `json:"name"`
	Phone *string `json:"phone"`
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

type ListResponse struct {
	Users []User `json:"users"`
}
