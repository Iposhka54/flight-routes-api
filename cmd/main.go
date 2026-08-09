package main

import (
	"context"
	"flight-routes-api/internal/config"
	"flight-routes-api/internal/sqlite"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to create logger", zap.Error(err))
	}

	if err = godotenv.Load(); err != nil {
		log.Fatal("failed to load env variables from .env file", zap.Error(err))
	}

	cfg, err := config.New()

	if err != nil {
		log.Fatal("failed to init config", zap.Error(err))
	}

	db, err := sqlite.New(cfg.DatabaseConfig)
	if err != nil {
		log.Fatal("failed to create sqlite client", zap.Error(err))
	}
	defer db.Close()
}
