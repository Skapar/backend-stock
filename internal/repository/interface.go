package repository

import (
	"context"

	"github.com/onec-tech/bot/internal/models/entities"
)

type PGRepository interface {
	// User
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
	GetUserByTGID(ctx context.Context, tgID int64) (*entities.User, error)

	// Receipt
	CreateReceipt(ctx context.Context, receipt *entities.Receipt) error
}
