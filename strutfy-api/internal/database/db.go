package database

import (
	"context"
	"database/sql"
    "log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func New(ctx context.Context, databaseURL string) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseURL)))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
    log.Println("Database connected")
    if err := RunMigrations(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}