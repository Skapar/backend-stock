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
}
