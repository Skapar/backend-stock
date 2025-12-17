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

func (s *service) CreateStock(ctx context.Context, stock *entities.Stock) (int64, error) {
	id, err := s.pgRepository.CreateStock(ctx, stock)
	if err != nil {
		s.log.Errorf("Service.CreateStock failed: %v", err)
		return 0, err
	}
	return id, nil
}

func (s *service) GetStockByID(ctx context.Context, id int64) (*entities.Stock, error) {
	return s.pgRepository.GetStockByID(ctx, id)
}

func (s *service) GetAllStocks(ctx context.Context) ([]*entities.Stock, error) {
	return s.pgRepository.GetAllStocks(ctx)
}

func (s *service) UpdateStock(ctx context.Context, stock *entities.Stock) error {
	return s.pgRepository.UpdateStock(ctx, stock)
}

func (s *service) DeleteStock(ctx context.Context, id int64) error {
	return s.pgRepository.DeleteStock(ctx, id)
}

// Order
func (s *service) CreateOrder(ctx context.Context, order *entities.Order) (int64, error) {
	return s.pgRepository.CreateOrder(ctx, order)
}

func (s *service) UpdateOrderStatus(ctx context.Context, orderID int64, status entities.OrderStatus) error {
	return s.pgRepository.UpdateOrderStatus(ctx, orderID, status)
}

func (s *service) GetOrdersByUserID(ctx context.Context, userID int64) ([]*entities.Order, error) {
	return s.pgRepository.GetOrdersByUserID(ctx, userID)
}

func (s *service) GetOrderByID(ctx context.Context, orderID int64) (*entities.Order, error) {
	return s.pgRepository.GetOrderByID(ctx, orderID)
}

// Portfolio
func (s *service) GetPortfolio(ctx context.Context, userID, stockID int64) (*entities.Portfolio, error) {
	return s.pgRepository.GetPortfolio(ctx, userID, stockID)
}

func (s *service) CreateOrUpdatePortfolio(ctx context.Context, p *entities.Portfolio) error {
	return s.pgRepository.CreateOrUpdatePortfolio(ctx, p)
}

// History
func (s *service) AddHistoryRecord(ctx context.Context, h *entities.History) (int64, error) {
	return s.pgRepository.AddHistoryRecord(ctx, h)
}

func (s *service) GetHistoryByUserID(ctx context.Context, userID int64) ([]*entities.History, error) {
	return s.pgRepository.GetHistoryByUserID(ctx, userID)
}
