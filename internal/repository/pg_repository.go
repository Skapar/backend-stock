package repository

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/pkg/database"
	"github.com/Skapar/backend/pkg/logger"
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
