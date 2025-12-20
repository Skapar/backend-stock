package cqrs

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/service"
)

type cqrsImpl struct {
	svc service.Service
}

func NewCQRS(svc service.Service) (Command, Query) {
	c := &cqrsImpl{svc: svc}
	return c, c
}

// Commands
func (c *cqrsImpl) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	return c.svc.CreateUser(ctx, user)
}

func (c *cqrsImpl) UpdateUser(ctx context.Context, user *entities.User) error {
	return c.svc.UpdateUser(ctx, user)
}

func (c *cqrsImpl) DeleteUser(ctx context.Context, id int64) error {
	return c.svc.DeleteUser(ctx, id)
}

func (c *cqrsImpl) CreateStock(ctx context.Context, stock *entities.Stock) (int64, error) {
	return c.svc.CreateStock(ctx, stock)
}

func (c *cqrsImpl) UpdateStock(ctx context.Context, stock *entities.Stock) error {
	return c.svc.UpdateStock(ctx, stock)
}

func (c *cqrsImpl) DeleteStock(ctx context.Context, id int64) error {
	return c.svc.DeleteStock(ctx, id)
}

func (c *cqrsImpl) CreateOrder(ctx context.Context, order *entities.Order) (int64, error) {
	return c.svc.CreateOrder(ctx, order)
}

func (c *cqrsImpl) UpdateOrderStatus(ctx context.Context, orderID int64, status entities.OrderStatus) error {
	return c.svc.UpdateOrderStatus(ctx, orderID, status)
}

func (c *cqrsImpl) ExecuteOrder(ctx context.Context, order *entities.Order) error {
	return c.svc.ExecuteOrder(ctx, order)
}

func (c *cqrsImpl) CreateOrUpdatePortfolio(ctx context.Context, p *entities.Portfolio) error {
	return c.svc.CreateOrUpdatePortfolio(ctx, p)
}

func (c *cqrsImpl) AddHistoryRecord(ctx context.Context, h *entities.History) (int64, error) {
	return c.svc.AddHistoryRecord(ctx, h)
}

// Queries
func (c *cqrsImpl) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	return c.svc.GetUserByID(ctx, id)
}

func (c *cqrsImpl) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return c.svc.GetUserByEmail(ctx, email)
}

func (c *cqrsImpl) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	return c.svc.GetAllUsers(ctx)
}

func (c *cqrsImpl) GetStockByID(ctx context.Context, id int64) (*entities.Stock, error) {
	return c.svc.GetStockByID(ctx, id)
}

func (c *cqrsImpl) GetAllStocks(ctx context.Context) ([]*entities.Stock, error) {
	return c.svc.GetAllStocks(ctx)
}

func (c *cqrsImpl) GetOrdersByUserID(ctx context.Context, userID int64) ([]*entities.Order, error) {
	return c.svc.GetOrdersByUserID(ctx, userID)
}

func (c *cqrsImpl) GetOrderByID(ctx context.Context, orderID int64) (*entities.Order, error) {
	return c.svc.GetOrderByID(ctx, orderID)
}

func (c *cqrsImpl) GetPortfolio(ctx context.Context, userID, stockID int64) (*entities.Portfolio, error) {
	return c.svc.GetPortfolio(ctx, userID, stockID)
}

func (c *cqrsImpl) GetPortfoliosByUserID(ctx context.Context, userID int64) ([]*entities.Portfolio, error) {
	return c.svc.GetPortfoliosByUserID(ctx, userID)
}

func (c *cqrsImpl) GetHistoryByUserID(ctx context.Context, userID int64) ([]*entities.History, error) {
	return c.svc.GetHistoryByUserID(ctx, userID)
}
