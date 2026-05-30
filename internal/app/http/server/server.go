package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"music-library/config"
	"music-library/internal/service/song"

	_ "music-library/docs"
)

type Server struct {
	cfg         *config.AppConfig
	logger      *zap.Logger
	songService *song.Service
	srv         *http.Server
}

func New(cfg *config.AppConfig, logger *zap.Logger, songService *song.Service) *Server {
	return &Server{
		cfg:         cfg,
		logger:      logger,
		songService: songService,
		srv: &http.Server{
			Handler: initHandlers(songService),
			Addr:    ":" + cfg.HTTP.Port,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info(
			"http server listening",
			zap.String("host", s.cfg.HTTP.Host),
			zap.String("port", s.cfg.HTTP.Port),
		)

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}
	s.logger.Info("http server stopped")
	return nil
}

func initHandlers(songService *song.Service) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := router.Group("/api")
	songsRouter := api.Group("/songs")

	songsRouter.GET("/", songService.GetSongs)
	songsRouter.GET("/:id/verses", songService.GetVerses)
	songsRouter.DELETE("/:id", songService.Delete)
	songsRouter.PATCH("/:id", songService.Edit)
	songsRouter.POST("/", songService.Add)

	return router
}
