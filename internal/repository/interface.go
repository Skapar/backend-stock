package repository

import (
	"context"

	"github.com/onec-tech/bot/internal/models/entities"
)

type PGRepository interface {
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
}
