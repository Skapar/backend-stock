package cqrs

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
)

type Command interface {
	CreateUser(ctx context.Context, user *entities.User) (int64, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id int64) error

	CreateStock(ctx context.Context, stock *entities.Stock) (int64, error)
	UpdateStock(ctx context.Context, stock *entities.Stock) error
	DeleteStock(ctx context.Context, id int64) error

	CreateOrder(ctx context.Context, order *entities.Order) (int64, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status entities.OrderStatus) error
	ExecuteOrder(ctx context.Context, order *entities.Order) error

	CreateOrUpdatePortfolio(ctx context.Context, p *entities.Portfolio) error

	AddHistoryRecord(ctx context.Context, h *entities.History) (int64, error)
}
