package organization

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, organization *Organization) error {
	_, err := r.db.NewInsert().Model(organization).Exec(ctx)
	return err
}