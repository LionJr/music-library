package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"music-library/config"
	"music-library/db"
	"music-library/internal/app/http/server"
	"music-library/internal/repository/postgres"
	"music-library/internal/service/song"
)

const shutdownTimeout = 15 * time.Second

type Application struct {
	cfg    *config.AppConfig
	logger *zap.Logger
	db     *sqlx.DB
	http   *server.Server
}

func New(ctx context.Context) (*Application, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	postgresDB, err := db.NewDB(ctx, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	songRepo := postgres.NewSongRepository(postgresDB)
	songService := song.NewService(cfg, logger, songRepo)

	return &Application{
		cfg:    cfg,
		logger: logger,
		db:     postgresDB,
		http:   server.New(cfg, logger, songService),
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	a.logger.Info("application started")
	return a.http.Run(ctx)
}

func (a *Application) Shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if a.http != nil {
		if err := a.http.Shutdown(shutdownCtx); err != nil {
			a.logger.Error("http shutdown failed", zap.Error(err))
		}
	}

	if a.db != nil {
		a.logger.Info("closing database connection")
		if err := a.db.Close(); err != nil {
			a.logger.Error("database close failed", zap.Error(err))
		}
	}

	if a.logger != nil {
		a.logger.Info("application stopped")
		if err := a.logger.Sync(); err != nil {
			a.logger.Error("logger sync failed", zap.Error(err))
		}
	}
}
