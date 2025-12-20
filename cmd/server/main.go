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

	"github.com/Skapar/backend/internal/handler"
	"github.com/Skapar/backend/internal/middleware"
	"github.com/Skapar/backend/internal/repository"
	"github.com/Skapar/backend/internal/service"
	"github.com/Skapar/backend/internal/worker"
	"github.com/Skapar/backend/pkg/cache"
	"github.com/Skapar/backend/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Skapar/backend/config"
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

	var cacheR *cache.Cache

	if cfg.RedisAddr != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:         cfg.RedisAddr,
			DialTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
		})

		cacheR = &cache.Cache{}
		cacheR.SetCacheImplementation(rdb)
	}

	// Подключение к БД
	db, err := database.New(cacheR, log, &database.Config{
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
		Cache:        cacheR,
		Log:          log,
		Config:       cfg,
	})
	if err != nil {
		log.Fatalf("failed to init service: %v", err)
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

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		log.Infow("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", duration.String(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := handler.NewAuthHandler(srv, cfg)
	userHandler := handler.NewUserHandler(srv)
	stockHandler := handler.NewStockHandler(srv, log)
	orderHandler := handler.NewOrderHandler(srv)
	portfolioHandler := handler.NewPortfolioHandler(srv)
	historyHandler := handler.NewHistoryHandler(srv)

	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		users := api.Group("/users")

		users.Use(middleware.AuthMiddleware(cfg))
		{
			users.GET("/me", userHandler.GetMe)
		}

		admin := api.Group("/users")
		admin.Use(middleware.AuthMiddleware(cfg, "ADMIN"))
		{
			admin.GET("/all", func(c *gin.Context) {
				users, err := srv.GetAllUsers(c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, users)
			})

			admin.GET("/:id", userHandler.GetUserByID)
			admin.PUT("/:id", userHandler.UpdateUser)
			admin.DELETE("/:id", userHandler.DeleteUser)
		}

		stocks := api.Group("/stocks")
		stocks.Use(middleware.AuthMiddleware(cfg))
		{
			stocks.GET("/", stockHandler.GetAllStocks)
			stocks.GET("/:id", stockHandler.GetStockByID)
		}

		adminStocks := api.Group("/stocks")
		adminStocks.Use(middleware.AuthMiddleware(cfg, "ADMIN"))
		{
			adminStocks.POST("/", stockHandler.CreateStock)
			adminStocks.PUT("/:id", stockHandler.UpdateStock)
			adminStocks.DELETE("/:id", stockHandler.DeleteStock)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg))
		{
			orders.POST("/", orderHandler.CreateOrder)
			orders.GET("/user/:user_id", orderHandler.GetOrdersByUser)
			orders.GET("/me", orderHandler.GetOrdersByUser)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}

		portfolio := api.Group("/portfolio")
		portfolio.Use(middleware.AuthMiddleware(cfg))
		{
			portfolio.GET("/:user_id/:stock_id", portfolioHandler.GetPortfolio)
			portfolio.POST("/", portfolioHandler.CreateOrUpdatePortfolio)
		}

		history := api.Group("/history")
		history.Use(middleware.AuthMiddleware(cfg))
		{
			history.POST("/", historyHandler.AddHistory)
			history.GET("/user/:user_id", historyHandler.GetHistoryByUser)
			history.GET("/me", historyHandler.GetHistoryByUser)
		}
	}

	//// HTTP server
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
