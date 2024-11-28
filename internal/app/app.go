package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"music-library/config"
	"music-library/internal/app/http/server"
	"music-library/internal/service/song"

	"go.uber.org/zap"
)

type Application struct {
	Config *config.AppConfig
	Logger *zap.Logger

	PostgresDb *sqlx.DB

	SongService *song.Service
}

func (app *Application) Run(ctx context.Context) {
	httpServerErrCh := server.NewServer(
		ctx,
		app.Logger,
		app.Config,
		app.SongService,
	)

	<-httpServerErrCh
}

func (app *Application) Shutdown() {
	app.Logger.Info("Shutdown database")
	_ = app.PostgresDb.Close()

	app.Logger.Info("Shutdown logger")
	_ = app.Logger.Sync()
}
