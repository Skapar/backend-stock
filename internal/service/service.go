package service

import (
	"context"

	"github.com/Skapar/backend-go/config"
	"github.com/Skapar/backend-go/internal/models/entities"
	"github.com/Skapar/backend-go/internal/repository"
	"github.com/Skapar/backend-go/pkg/cache"
	"github.com/Skapar/backend-go/pkg/logger"
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
