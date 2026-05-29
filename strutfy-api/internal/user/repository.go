package user

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user := new(User)

	err := r.db.NewSelect().
		Model(user).
		Column("id", "organization_id", "name", "email", "password_hash", "role").
		Where("email = ?", email).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
	user := new(User)

	err := r.db.NewSelect().
		Model(user).
		Column("id", "organization_id", "name", "email", "password_hash", "role").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}
