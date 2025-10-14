package service

import (
	"context"

	"github.com/Skapar/backend-go/internal/models/entities"
)

type Service interface {
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
}
