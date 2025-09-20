package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/projectWithOrder/configs"
	"github.com/t1xelLl/projectWithOrder/internal/handler"
	"github.com/t1xelLl/projectWithOrder/internal/repository"
	"github.com/t1xelLl/projectWithOrder/internal/service"
	"github.com/t1xelLl/projectWithOrder/pkg/httpserver"
	"github.com/t1xelLl/projectWithOrder/pkg/logger"
	"github.com/t1xelLl/projectWithOrder/pkg/postgres"
	"log"
)

func main() {
	// DONE: logger: logrus
	logger.SetLogrus()

	// DONE: configs: viper
	cfg, err := configs.LoadConfig("./configs/config.yaml")
	if err != nil {
		logrus.Fatalf("Error loading config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env files: %s", err.Error())
	}

	// DONE: postgres: sqlx
	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("init postgres error: %s", err.Error())
	}

	// DONE: repository, service and handler
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(httpserver.Server)

	// DONE: init router : gin
	if err := srv.Run(cfg.Http.Port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("http server error: %s", err.Error())
	}
}

/*

// TODO: логику ручки

// TODO: init cache and restore from db: redis

// TODO: init Kafka consumer

// TODO: start Kafka consumer

// TODO: graceful shutdown

//TODO: env : godotenv
*/
