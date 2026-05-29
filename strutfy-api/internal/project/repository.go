package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, project *Project) error {
	_, err := r.db.NewInsert().
		Model(project).
		Exec(ctx)

	return err
}

func (r *Repository) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]Project, error) {
	var projects []Project

	err := r.db.NewSelect().
		Model(&projects).
		Where("org_id = ?", orgID).
		Order("created_at DESC").
		Scan(ctx)

	return projects, err
}

func (r *Repository) FindByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Project, error) {
	project := new(Project)

	err := r.db.NewSelect().
		Model(project).
		Where("org_id = ?", orgID).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *Repository) Update(ctx context.Context, project *Project) error {
	_, err := r.db.NewUpdate().
		Model(project).
		WherePK().
		Where("org_id = ?", project.OrgID).
		Exec(ctx)

	return err
}

func (r *Repository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*Project)(nil)).
		Where("org_id = ?", orgID).
		Where("id = ?", id).
		Exec(ctx)

	return err
}
