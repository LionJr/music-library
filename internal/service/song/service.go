package song

import (
	"github.com/LionJr/music-library/config"
	"go.uber.org/zap"
)

type Service struct {
	config *config.AppConfig
	Logger *zap.Logger

	Repo Repo
}

func NewService(cfg *config.AppConfig, logger *zap.Logger, repo Repo) *Service {
	return &Service{
		config: cfg,
		Logger: logger,

		Repo: repo,
	}
}
