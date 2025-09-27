package repository

import (
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
