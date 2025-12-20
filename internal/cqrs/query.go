package cqrs

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
)

type Query interface {
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]*entities.User, error)

	GetStockByID(ctx context.Context, id int64) (*entities.Stock, error)
	GetAllStocks(ctx context.Context) ([]*entities.Stock, error)

	GetOrdersByUserID(ctx context.Context, userID int64) ([]*entities.Order, error)
	GetOrderByID(ctx context.Context, orderID int64) (*entities.Order, error)

	GetPortfolio(ctx context.Context, userID, stockID int64) (*entities.Portfolio, error)
	GetPortfoliosByUserID(ctx context.Context, userID int64) ([]*entities.Portfolio, error)

	GetHistoryByUserID(ctx context.Context, userID int64) ([]*entities.History, error)
}
