package service

import (
	"context"

	"github.com/Skapar/backend/config"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/repository"
	"github.com/Skapar/backend/pkg/cache"
	"github.com/Skapar/backend/pkg/logger"
	stock "github.com/Skapar/backend/proto"
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

func (s *service) CreateUser(ctx context.Context, req *stock.CreateUserRequest) (int64, error) {
	user := &entities.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     entities.RoleTrader,
		Balance:  0,
	}

	id, err := s.pgRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	user, err := s.pgRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, user *entities.User) error {
	if err := s.pgRepository.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	if err := s.pgRepository.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	users, err := s.pgRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
