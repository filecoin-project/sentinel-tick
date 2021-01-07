package main

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func setupDB(ctx context.Context, url string, poolSize int, appName string) (*pg.DB, error) {
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	opt.PoolSize = poolSize
	opt.ApplicationName = appName

	db := pg.Connect(opt)
	db = db.WithContext(ctx)
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("no ping from db: %w", err)
	}

	if _, err = db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS "+schema+";"); err != nil {
		teardownDB(ctx, db)
		return nil, err
	}

	return db, nil
}

func teardownDB(ctx context.Context, db *pg.DB) error {
	return db.Close()
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	return db.Model(&Quote{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
}
