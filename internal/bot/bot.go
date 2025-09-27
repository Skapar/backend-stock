package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/onec-tech/bot/internal/service"
	"github.com/onec-tech/bot/pkg/logger"
)

type telegramBot struct {
	tg      *tgbotapi.BotAPI
	log     logger.Logger
	service service.Service
}

func NewBot(token string, log logger.Logger, service service.Service) (Bot, error) {
	tg, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	tg.Debug = false

	return &telegramBot{
		tg:      tg,
		log:     log,
		service: service,
	}, nil
}

func (b *telegramBot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tg.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(update)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
