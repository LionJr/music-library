package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"music-library/config"
	"music-library/db"
	"music-library/internal/app"
	"music-library/internal/repository/postgres"
	"music-library/internal/service/song"
)

// @title           Swagger Example API
// @version         1.0
// @description     Music library example.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("read environment variables error: %s", err.Error())
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("create logger error: %s", err.Error())
	}

	songDb, err := db.NewDB(&cfg.Postgres)
	if err != nil {
		logger.Fatal("database connect error: %s", zap.Error(err))
	}

	application := &app.Application{
		Config:      cfg,
		Logger:      logger,
		SongService: song.NewService(cfg, logger, postgres.NewSongRepository(songDb)),
	}

	application.Run(ctx)
	application.Shutdown()
}
