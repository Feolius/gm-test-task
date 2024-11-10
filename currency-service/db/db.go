package db

import (
	"context"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations(ctx context.Context, dbstring string) error {
	db, err := goose.OpenDBWithDriver("mysql", dbstring)
	if err != nil {
		return fmt.Errorf("migrations: failed to open DB: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			// @todo may need to return a specific error and handle it separately
			log.Printf("goose: failed to close DB: %v", err)
		}
	}()
	if err = goose.RunContext(ctx, "up", db, "db/migrations"); err != nil {
		return fmt.Errorf("migrations: failed to run migrations: %w", err)
	}
	return nil
}
