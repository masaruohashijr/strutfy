package project

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProjectStatus string

const (
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusPaused     ProjectStatus = "paused"
	ProjectStatusDone       ProjectStatus = "done"
)

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID          uuid.UUID      `bun:"id,pk,type:uuid" json:"id"`
	OrgID       uuid.UUID      `bun:"org_id,type:uuid,notnull" json:"org_id"`
	Name        string         `bun:"name,notnull" json:"name"`
	Description string         `bun:"description" json:"description"`
	Status      ProjectStatus  `bun:"status,notnull" json:"status"`
	Props       map[string]any `bun:"props,type:jsonb" json:"props"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}
