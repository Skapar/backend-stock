package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Skapar/backend-go/internal/repository"
	"github.com/Skapar/backend-go/internal/service"
	"github.com/Skapar/backend-go/internal/worker"
	"github.com/Skapar/backend-go/pkg/cache"
	"github.com/Skapar/backend-go/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Skapar/backend-go/config"
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

	srv, err = service.NewService(&service.SConfig{
		PGRepository: pgRepository,
		Cache:        cacheImpl,
		Log:          log,
		Config:       cfg,
	})
	if err != nil {
		log.Fatalf("failed to re-init service with notifier: %v", err)
	}

	wrk := worker.NewWorker(&worker.WorkerConfig{
		Service: srv,
		Log:     log,
	})

	wrk.Start()

	// Gin
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health"}}))
	router.Use(gin.Recovery())

	// CORS
	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"http://localhost:3001",
			"https://localhost:3001",
			"http://127.0.0.1:8080",
			"http://localhost:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"X-Content-Type, Content-Length", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Static files (замена http.FileServer)
	router.Static("/files", "./receipts")

	// HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.ListenHttpPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server
	go func() {
		log.Infof("HTTP server started on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server forced to shutdown: %v", err)
	}

	wrk.Stop()
	log.Info("Server exited properly")
}
