package database

import (
	"embed"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/major75/online-subscriptions/pkg/logger"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(db *DB, log logger.Logger) error {
	goose.SetBaseFS(embedMigrations)

	conn := stdlib.OpenDBFromPool(db.pool)
	defer conn.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(conn, "migrations"); err != nil {
		return err
	}

	log.Info("Migrations applied successfully")
	return nil
}
