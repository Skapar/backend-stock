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

	commands := []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Начать работу",
		},
		{
			Command:     "mysubscription",
			Description: "Моя подписка",
		},
	}

	setCmdCfg := tgbotapi.NewSetMyCommands(commands...)

	_, err = tg.Request(setCmdCfg)
	if err != nil {
		cfg.Log.Errorf("failed to set commands: %v", err)
	}

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
	sem := make(chan struct{}, b.config.MaxWorkers)

	for {
		select {
		case update := <-updates:
			sem <- struct{}{}

			go func(update tgbotapi.Update) {
				defer func() { <-sem }()
				b.handleUpdate(update)
			}(update)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *telegramBot) Notify(ctx context.Context, tgID int64, message string) error {
	msg := tgbotapi.NewMessage(tgID, message)
	_, err := b.tg.Send(msg)
	return err
}
