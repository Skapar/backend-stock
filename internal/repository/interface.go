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
	GetApprovedReceipts(ctx context.Context) ([]entities.Receipt, error)
	UpdateReceiptStatus(ctx context.Context, status entities.StatusType, receiptID int64) error
	GetDefaultSubscription(ctx context.Context) (*entities.Subscription, error)

	// UserSubcription
	CreateUserSubscription(ctx context.Context, userID int64, sub *entities.Subscription) (int64, error)

	// Payment
	CreatePayment(ctx context.Context, userID, userSubID int64, amount int64) error
}
