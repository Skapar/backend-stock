package service

import (
	"context"

	"github.com/onec-tech/bot/config"
	"github.com/onec-tech/bot/internal/models/entities"
	"github.com/onec-tech/bot/internal/repository"
	"github.com/onec-tech/bot/pkg/cache"
	"github.com/onec-tech/bot/pkg/logger"
)

type service struct {
	pgRepository repository.PGRepository
	cache        cache.ICache
	log          logger.Logger
	config       *config.Config
}

type SConfig struct {
	PGRepository repository.PGRepository
	Cache        cache.ICache
	Log          logger.Logger
	Config       *config.Config
}

func NewService(cfg *SConfig) (Service, error) {
	return &service{
		pgRepository: cfg.PGRepository,
		cache:        cfg.Cache,
		log:          cfg.Log,
		config:       cfg.Config,
	}, nil
}

func (s *service) CreateOrUpdateUser(ctx context.Context, user *entities.User) error {
	err := s.pgRepository.CreateOrUpdateUser(ctx, user)
	if err != nil {
		s.log.Errorf("failed to create user: %v", err)
		return err
	}
	return nil
}
