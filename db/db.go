package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/LionJr/music-library/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(ctx context.Context, cfg *config.Postgres) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if err = runMigrations(cfg); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func runMigrations(cfg *config.Postgres) error {
	migrationDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	m, err := migrate.New("file://migrations", migrationDSN)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
