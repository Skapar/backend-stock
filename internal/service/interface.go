package service

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
)

type Service interface {
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
}
