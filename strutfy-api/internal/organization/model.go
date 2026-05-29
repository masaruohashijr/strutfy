package organization

import (
	"time"

	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:organizations"`

	ID        string    `bun:"id,pk"`
	Name      string    `bun:"name,notnull"`
	Document  string    `bun:"document"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}