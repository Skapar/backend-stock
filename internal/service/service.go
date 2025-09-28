package service

import (
	"context"

	"github.com/onec-tech/bot/config"
	"github.com/onec-tech/bot/internal/models/entities"
	"github.com/onec-tech/bot/internal/repository"
	"github.com/onec-tech/bot/pkg/cache"
	"github.com/onec-tech/bot/pkg/logger"
)

type service struct {
	pgRepository repository.PGRepository
	cache        cache.ICache
	log          logger.Logger
	config       *config.Config
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

func (s *service) CreateOrUpdateUser(ctx context.Context, user *entities.User) error {
	err := s.pgRepository.CreateOrUpdateUser(ctx, user)
	if err != nil {
		s.log.Errorf("failed to create user: %v", err)
		return err
	}
	return nil
}

func (s *service) GetUserByTGID(ctx context.Context, tgID int64) (*entities.User, error) {
	return s.pgRepository.GetUserByTGID(ctx, tgID)
}

func (s *service) CreateReceipt(ctx context.Context, userID int64, filePath string) error {
	receipt := &entities.Receipt{
		UserID:   userID,
		FilePath: filePath,
		Status:   entities.StatusPending,
	}

	return s.pgRepository.CreateReceipt(ctx, receipt)
}

func (s *service) ProcessApprovedReceipts(ctx context.Context) error {
	receipts, err := s.pgRepository.GetApprovedReceipts(ctx)
	if err != nil {
		s.log.Errorf("failed to get approved receipts: %v", err)
		return err
	}

	for _, r := range receipts {
		err := s.pgRepository.UpdateReceiptStatus(ctx, entities.StatusConfirmed, r.ID)
		if err != nil {
			s.log.Errorf("failed to confirm receipt %d: %v", r.ID, err)
			return err
		}

		sub, err := s.pgRepository.GetDefaultSubscription(ctx)
		if err != nil {
			s.log.Errorf("failed to get subscription for user %d: %v", r.UserID, err)
			return err
		}

		userSubID, err := s.pgRepository.CreateUserSubscription(ctx, r.UserID, sub)
		if err != nil {
			s.log.Errorf("failed to create user subscription for user %d: %v", r.UserID, err)
			return err
		}

		err = s.pgRepository.CreatePayment(ctx, r.UserID, userSubID, sub.Price)
		if err != nil {
			s.log.Errorf("failed to create payment for user %d: %v", r.UserID, err)
			return err
		}
	}

	return nil
}
