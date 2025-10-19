package service

import (
	"context"

	"github.com/Skapar/backend/config"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/repository"
	"github.com/Skapar/backend/pkg/cache"
	"github.com/Skapar/backend/pkg/logger"
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

func (s *service) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	id, err := s.pgRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	return s.pgRepository.GetUserByID(ctx, id)
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.pgRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Errorf("Service.GetUserByEmail failed: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, user *entities.User) error {
	return s.pgRepository.UpdateUser(ctx, user)
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	return s.pgRepository.DeleteUser(ctx, id)
}

func (s *service) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	return s.pgRepository.GetAllUsers(ctx)
}
