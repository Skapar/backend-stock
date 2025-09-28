package worker

import (
	"context"
	"time"

	"github.com/onec-tech/bot/pkg/logger"

	"github.com/go-co-op/gocron"
	"github.com/onec-tech/bot/internal/service"
)

type worker struct {
	service   service.Service
	log       logger.Logger
	scheduler *gocron.Scheduler
}

type WorkerConfig struct {
	Service service.Service
	Log     logger.Logger
}

func NewWorker(cfg *WorkerConfig) Worker {
	return &worker{
		service:   cfg.Service,
		log:       cfg.Log,
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (w *worker) Start() {
	w.startProcessApprovedReceipts()
	w.scheduler.StartAsync()
}

func (w *worker) Stop() {
	w.scheduler.Stop()
	w.log.Info("Scheduler stopping...")
}

func (w *worker) startProcessApprovedReceipts() {
	_, err := w.scheduler.Every(5).Seconds().Do(func() {
		w.log.Info("ProcessApprovedReceipts task started")
		err := w.service.ProcessApprovedReceipts(context.Background())
		if err != nil {
			w.log.Error("ProcessApprovedReceipts", err)
		}
	})
	if err != nil {
		w.log.Error("ProcessApprovedReceipts ", err)
	}
}
