package service

import (
	"context"

	"github.com/onec-tech/bot/internal/models/entities"
)

type Service interface {
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
}
