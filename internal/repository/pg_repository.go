package repository

import (
	"context"
	"time"

	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/pkg/database"
	"github.com/Skapar/backend/pkg/logger"
	"github.com/pkg/errors"
)

type pgRepository struct {
	DB  database.IDatabase
	log logger.Logger
}

func NewPGRepository(db database.IDatabase, log logger.Logger) PGRepository {
	return &pgRepository{
		DB:  db,
		log: log,
	}
}

func (r *pgRepository) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	q := `
		INSERT INTO stock_user (email, password, role, balance)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int64
	if err := r.DB.Insert(ctx, &id, q, user.Email, user.Password, user.Role, user.Balance); err != nil {
		return 0, errors.Wrap(err, "CreateUser: failed to create user")
	}
	return id, nil
}

func (r *pgRepository) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	q := `
		SELECT id, email, password, role, balance::float8, created_at
		FROM stock_user
		WHERE id = $1;
	`

	var user entities.User
	if err := r.DB.GetOne(ctx, &user, q, id); err != nil {
		r.log.Errorf("GetUserByID error: %v (id=%d)", err, id)
		return nil, errors.Wrap(err, "GetUserByID: failed to get user")
	}
	return &user, nil
}

func (r *pgRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	q := `
        SELECT id, email, password, role, balance, created_at
        FROM stock_user
        WHERE email = $1;
    `
	var user entities.User
	if err := r.DB.GetOne(ctx, &user, q, email); err != nil {
		r.log.Errorf("GetUserByEmail error: %v (email=%s)", err, email)
		return nil, errors.Wrap(err, "GetUserByEmail: failed to get user")
	}
	return &user, nil
}

func (r *pgRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	q := `
		UPDATE stock_user
		SET email = $1,
			password = $2,
			role = $3,
			balance = $4
		WHERE id = $5
		RETURNING id;
	`

	var updatedID int64
	if err := r.DB.Update(ctx, &updatedID, q, user.Email, user.Password, user.Role, user.Balance, user.ID); err != nil {
		return errors.Wrap(err, "UpdateUser: failed to update user")
	}
	return nil
}

func (r *pgRepository) DeleteUser(ctx context.Context, id int64) error {
	q := `DELETE FROM stock_user WHERE id = $1;`
	if err := r.DB.Delete(ctx, nil, q, id); err != nil {
		return errors.Wrapf(err, "DeleteUser: failed to delete user")
	}
	return nil
}

func (r *pgRepository) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	q := `
		SELECT id, email, password, role, balance, created_at
		FROM stock_user
		ORDER BY id DESC;
	`

	var users []*entities.User
	if err := r.DB.Get(ctx, &users, q); err != nil {
		return nil, errors.Wrap(err, "GetAllUsers: failed to get all users")
	}
	return users, nil
}

func (r *pgRepository) CreateStock(ctx context.Context, stock *entities.Stock) (int64, error) {
	q := `
		INSERT INTO stock_stock (symbol, name, price, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int64
	if err := r.DB.Insert(ctx, &id, q, stock.Symbol, stock.Name, stock.Price, time.Now()); err != nil {
		return 0, errors.Wrap(err, "CreateStock: failed to create stock")
	}
	return id, nil
}

func (r *pgRepository) GetStockByID(ctx context.Context, id int64) (*entities.Stock, error) {
	q := `
		SELECT id, symbol, name, price, updated_at
		FROM stock_stock
		WHERE id = $1;
	`

	var stock entities.Stock
	if err := r.DB.GetOne(ctx, &stock, q, id); err != nil {
		return nil, errors.Wrap(err, "GetStockByID: failed to get stock")
	}
	return &stock, nil
}

func (r *pgRepository) GetAllStocks(ctx context.Context) ([]*entities.Stock, error) {
	q := `
		SELECT id, symbol, name, price, updated_at
		FROM stock_stock
		ORDER BY id DESC;
	`

	var stocks []*entities.Stock
	if err := r.DB.Get(ctx, &stocks, q); err != nil {
		return nil, errors.Wrap(err, "GetAllStocks: failed to get all stocks")
	}
	return stocks, nil
}

func (r *pgRepository) UpdateStock(ctx context.Context, stock *entities.Stock) error {
	q := `
		UPDATE stock_stock
		SET symbol = $1,
			name = $2,
			price = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id;
	`

	var updatedID int64
	if err := r.DB.Update(ctx, &updatedID, q, stock.Symbol, stock.Name, stock.Price, time.Now(), stock.ID); err != nil {
		return errors.Wrap(err, "UpdateStock: failed to update stock")
	}
	return nil
}

func (r *pgRepository) DeleteStock(ctx context.Context, id int64) error {
	q := `DELETE FROM stock_stock WHERE id = $1;`
	if err := r.DB.Delete(ctx, nil, q, id); err != nil {
		return errors.Wrapf(err, "DeleteStock: failed to delete stock")
	}
	return nil
}
