package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/projectWithOrder/configs"
	"github.com/t1xelLl/projectWithOrder/internal/handler"
	"github.com/t1xelLl/projectWithOrder/internal/repository"
	"github.com/t1xelLl/projectWithOrder/internal/service"
	cache2 "github.com/t1xelLl/projectWithOrder/internal/service/cache"
	"github.com/t1xelLl/projectWithOrder/pkg/httpserver"
	"github.com/t1xelLl/projectWithOrder/pkg/logger"
	"github.com/t1xelLl/projectWithOrder/pkg/postgres"
	"github.com/t1xelLl/projectWithOrder/pkg/redis"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// DONE: logger: logrus
	logger.SetLogrus()

	// DONE: configs: viper
	cfg, err := configs.LoadConfig("./configs/config.yaml")
	if err != nil {
		logrus.Fatalf("Error loading config: %s", err.Error())
	}
	// DONE: .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env files: %s", err.Error())
	}

	// DONE: postgres: sqlx
	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("init postgres error: %s", err.Error())
	}
	defer db.Close()

	// DONE: init cache
	redisClient, err := redis.NewClientRedis(cfg.Redis)
	if err != nil {
		logrus.Fatalf("init redis client error: %s", err.Error())
	}
	defer redisClient.Close()

	cache := cache2.NewCache(redisClient)

	// DONE: repository, service
	repos := repository.NewRepository(db)
	services := service.NewService(repos, cache)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// DONE restore from db
	if err := services.PreloadCache(ctx); err != nil {
		logrus.Warnf("Failed to restore cache: %v", err)
	}
	// DONE: handler
	handlers := handler.NewHandler(services)
	srv := new(httpserver.Server)

	// DONE: init router : gin
	go func() {
		logrus.Infof("Starting server on port %s", cfg.Http.Port)
		if err := srv.Run(cfg.Http.Port, handlers.InitRoutes()); err != nil {
			logrus.Errorf("http server error: %s", err.Error())
		}
	}()
	// DONE: graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.ShutDown(shutdownCtx); err != nil {
		logrus.Errorf("Error shutting down server: %s", err.Error())
	}

	logrus.Info("Server exited gracefully")
}

/*

// TODO: init Kafka consumer

// TODO: start Kafka consumer

// TODO: Web

*/
