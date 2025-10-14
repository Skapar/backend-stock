package worker

import (
	"time"

	"github.com/Skapar/backend/pkg/logger"

	"github.com/Skapar/backend/internal/service"
	"github.com/go-co-op/gocron"
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
	w.scheduler.StartAsync()
}

func (w *worker) Stop() {
	w.scheduler.Stop()
	w.log.Info("Scheduler stopping...")
}
