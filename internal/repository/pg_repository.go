package repository

import (
	"context"

	"github.com/onec-tech/bot/internal/models/entities"
	"github.com/onec-tech/bot/pkg/database"
	"github.com/onec-tech/bot/pkg/logger"
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

func (r *pgRepository) CreateOrUpdateUser(ctx context.Context, user *entities.User) error {
	q := `
	INSERT INTO onec_user (tg_id, nickname, name, phone)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (tg_id) 
	DO UPDATE SET
		nickname = EXCLUDED.nickname,
		name = EXCLUDED.name,
		phone = EXCLUDED.phone
	RETURNING id;
	`

	return r.DB.Insert(ctx, user, q, user.TGID, user.Nickname, user.Name, user.Phone)
}

func (r *pgRepository) GetUserByTGID(ctx context.Context, tgID int64) (*entities.User, error) {
	q := `
	SELECT id, tg_id, nickname, name, phone
	FROM onec_user
	WHERE tg_id = $1;
	`
	var user entities.User
	err := r.DB.GetOne(ctx, &user, q, tgID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgRepository) CreateReceipt(ctx context.Context, receipt *entities.Receipt) error {
	q := `
	INSERT INTO onec_receipt (user_id, file_path, status)
	VALUES ($1, $2, $3)
	RETURNING id;
	`
	return r.DB.Insert(ctx, receipt, q, receipt.UserID, receipt.FilePath, receipt.Status)
}
