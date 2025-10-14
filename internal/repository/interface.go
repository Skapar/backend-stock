package repository

import (
	"context"

	"github.com/Skapar/backend-go/internal/models/entities"
)

type PGRepository interface {
	// User
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
}
