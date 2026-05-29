package project

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateProjectInput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      ProjectStatus  `json:"status"`
	Props       map[string]any `json:"props"`
}

func (s *Service) Create(ctx context.Context, orgID uuid.UUID, input CreateProjectInput) (*Project, error) {
	status := input.Status
	if status == "" {
		status = ProjectStatusPlanning
	}

	now := time.Now().UTC()

	project := &Project{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        input.Name,
		Description: input.Description,
		Status:      status,
		Props:       input.Props,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if project.Props == nil {
		project.Props = map[string]any{}
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID) ([]Project, error) {
	return s.repo.ListByOrg(ctx, orgID)
}

func (s *Service) Find(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Project, error) {
	return s.repo.FindByID(ctx, orgID, id)
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(ctx, orgID, id)
}
