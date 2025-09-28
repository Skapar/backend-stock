package service

import (
	"context"

	"github.com/onec-tech/bot/internal/models/entities"
)

type Service interface {
	CreateOrUpdateUser(ctx context.Context, user *entities.User) error
	GetUserByTGID(ctx context.Context, tgID int64) (*entities.User, error)
	GetActiveSubscription(ctx context.Context, userID int64) (*entities.UserSubscription, error)

	CreateReceipt(ctx context.Context, userID int64, filePath string) error
	ProcessApprovedReceipts(ctx context.Context) error
}

type Notifier interface {
	Notify(ctx context.Context, tgID int64, message string) error
}
