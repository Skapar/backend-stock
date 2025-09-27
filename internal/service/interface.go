package service

import (
	"context"
)

type Service interface {
	StartTelegramBot(ctx context.Context) error
}
