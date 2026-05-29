package user

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID             string `bun:"id,pk"`
	OrganizationID string `bun:"organization_id"`

	Name         string `bun:"name"`
	Email        string `bun:"email"`
	PasswordHash string `bun:"password_hash"`
	Role         string `bun:"role"`

	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}