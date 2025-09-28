package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/onec-tech/bot/internal/bot"
	"github.com/onec-tech/bot/internal/repository"
	"github.com/onec-tech/bot/internal/service"
	"github.com/onec-tech/bot/pkg/cache"
	"github.com/onec-tech/bot/pkg/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/onec-tech/bot/config"
)

func main() {
	// Загружаем .env
	_ = godotenv.Load()

	// Логгер
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "trace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}
	zlogger := zap.New(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), os.Stdout, zap.DebugLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	log := zlogger.Sugar()

	// Конфиг
	cfg := config.New()
	cfg.Init()

	cacheImpl := new(cache.Cache)

	// Подключение к БД
	db, err := database.New(cacheImpl, log, &database.Config{
		PostgresMasterAddr: cfg.PostgresAddr,
		PostgresSlaveAddr:  cfg.PostgresAddr,
	})
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * repository layer
	 */

	pgRepository := repository.NewPGRepository(db, log)

	/*
	 * service layer
	 */
	// Service
	srv, err := service.NewService(&service.SConfig{
		PGRepository: pgRepository,
		Cache:        cacheImpl,
		Log:          log,
		Config:       cfg,
	})
	if err != nil {
		log.Fatalf("failed to init service: %v", err)
	}

	// Telegram Bot init
	tgBot, err := bot.NewBot(&bot.BotConfig{
		Service: srv,
		Log:     log,
		Config:  cfg,
	})
	if err != nil {
		log.Fatalf("failed to init telegram bot: %v", err)
	}

	// Start Telegram Bot
	go func() {
		if err := tgBot.Start(context.Background()); err != nil {
			log.Errorf("telegram bot error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Даем немного времени на завершение горутин
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer shutdownCancel()

	<-shutdownCtx.Done()
}
