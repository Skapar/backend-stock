package service

import (
	"context"

	"github.com/onec-tech/bot/config"
	"github.com/onec-tech/bot/internal/bot"
	"github.com/onec-tech/bot/internal/repository"
	"github.com/onec-tech/bot/pkg/cache"
	"github.com/onec-tech/bot/pkg/logger"
)

type service struct {
	pgRepository repository.PGRepository
	cache        cache.ICache
	log          logger.Logger
	config       *config.Config
	bot          bot.Bot
}

type SConfig struct {
	PGRepository repository.PGRepository
	Cache        cache.ICache
	Log          logger.Logger
	Config       *config.Config
}

func NewService(cfg *SConfig) (Service, error) {
	return &service{
		pgRepository: cfg.PGRepository,
		cache:        cfg.Cache,
		log:          cfg.Log,
		config:       cfg.Config,
	}, nil
}

func (s *service) StartTelegramBot(ctx context.Context) error {
	b, err := bot.New(s.config.TelegramToken, s.log) // вот тут используем New
	if err != nil {
		return err
	}
	return b.Start(ctx)
}
