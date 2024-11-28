package song

import (
	"go.uber.org/zap"
	"music-library/config"
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
