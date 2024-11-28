package server

import (
	"context"
	"github.com/gin-gonic/gin"
	errch "github.com/proxeter/errors-channel"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"music-library/config"
	"music-library/internal/service/song"
	"net/http"

	_ "music-library/docs"
)

type Server struct {
	Logger *zap.Logger
	Config *config.AppConfig

	SongService *song.Service
}

func NewServer(
	ctx context.Context,
	logger *zap.Logger,
	config *config.AppConfig,
	musicService *song.Service,
) <-chan error {
	return errch.Register(func() error {
		return (&Server{
			Logger:      logger,
			Config:      config,
			SongService: musicService,
		}).Start(ctx)
	})
}

func (s *Server) Start(ctx context.Context) error {
	h := s.initHandlers()

	server := http.Server{
		Handler: h,
		Addr:    ":" + s.Config.HTTP.Port,
	}

	s.Logger.Info(
		"Server running",
		zap.String("host", s.Config.HTTP.Host),
		zap.String("port", s.Config.HTTP.Port),
	)

	select {
	case err := <-errch.Register(server.ListenAndServe):
		s.Logger.Info("Shutdown music_library server", zap.String("by", "error"), zap.Error(err))
		return server.Shutdown(ctx)
	case <-ctx.Done():
		s.Logger.Info("Shutdown music_library server", zap.String("by", "context.Done"))
		return server.Shutdown(ctx)
	}
}

func (s *Server) initHandlers() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := router.Group("/api")
	songsRouter := api.Group("/songs")

	songsRouter.GET("/", s.SongService.GetSongs)
	songsRouter.GET("/:id/verses", s.SongService.GetVerses)
	songsRouter.DELETE("/:id", s.SongService.Delete)
	songsRouter.PATCH("/:id", s.SongService.Edit)
	songsRouter.POST("/", s.SongService.Add)

	return router
}
