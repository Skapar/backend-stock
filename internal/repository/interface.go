package repository

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
)

type PGRepository interface {
	// User
	CreateUser(ctx context.Context, user *entities.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]*entities.User, error)

	// --- Stock ---
	CreateStock(ctx context.Context, stock *entities.Stock) (int64, error)
	GetStockByID(ctx context.Context, id int64) (*entities.Stock, error)
	GetAllStocks(ctx context.Context) ([]*entities.Stock, error)
	UpdateStock(ctx context.Context, stock *entities.Stock) error
	DeleteStock(ctx context.Context, id int64) error

	// --- Orders ---
	CreateOrder(ctx context.Context, order *entities.Order) (int64, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status entities.OrderStatus) error
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*entities.Order, error)
	GetOrderByID(ctx context.Context, orderID int64) (*entities.Order, error)

	// --- Portfolio ---
	GetPortfolio(ctx context.Context, userID, stockID int64) (*entities.Portfolio, error)
	CreateOrUpdatePortfolio(ctx context.Context, portfolio *entities.Portfolio) error
	GetPortfoliosByUserID(ctx context.Context, userID int64) ([]*entities.Portfolio, error)

	// --- History ---
	AddHistoryRecord(ctx context.Context, h *entities.History) (int64, error)
	GetHistoryByUserID(ctx context.Context, userID int64) ([]*entities.History, error)
}
