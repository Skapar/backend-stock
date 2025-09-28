package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/onec-tech/bot/config"
	"github.com/onec-tech/bot/internal/service"
	"github.com/onec-tech/bot/pkg/logger"
)

type telegramBot struct {
	tg      *tgbotapi.BotAPI
	log     logger.Logger
	service service.Service
	config  *config.Config
}

type BotConfig struct {
	Service service.Service
	Log     logger.Logger
	Config  *config.Config
}

func NewBot(cfg *BotConfig) (Bot, error) {
	tg, err := tgbotapi.NewBotAPI(cfg.Config.TelegramToken)
	if err != nil {
		return nil, err
	}

	tg.Debug = cfg.Config.TelegramDebug

	return &telegramBot{
		tg:      tg,
		service: cfg.Service,
		log:     cfg.Log,
		config:  cfg.Config,
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
