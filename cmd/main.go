package main

import (
	"context"
	"flight-routes-api/internal/config"
	"flight-routes-api/internal/handler"
	"flight-routes-api/internal/middleware"
	"flight-routes-api/internal/repository"
	"flight-routes-api/internal/service"
	"flight-routes-api/internal/sqlite"
	"fmt"
	"net/http"
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

	log.Info("sqlite client creation completed")

	if err = sqlite.InitDB(db); err != nil {
		log.Fatal("failed to init db", zap.Error(err))
	}

	log.Info("database initialization completed")

	airportRepository := repository.NewAirportRepository(db)

	airportService := service.NewAirportService(airportRepository)

	airportHandler := handler.NewAirportHandler(airportService)

	mux := http.NewServeMux()
	withErrors := func(h middleware.HandlerFunc) http.HandlerFunc {
		return middleware.ErrorHandler(log, h)
	}

	mux.HandleFunc("GET /airports", withErrors(airportHandler.GetAirports))
	mux.HandleFunc("GET /airport/{iataCode}", withErrors(airportHandler.GetAirportByIataCode))
	mux.HandleFunc("POST /airports", withErrors(airportHandler.CreateAirport))

	server := http.NewServeMux()
	server.Handle("/api/", http.StripPrefix("/api", mux))

	addr := fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port)
	log.Info("http server listening",
		zap.String("addr", addr))
	if err = http.ListenAndServe(addr, server); err != nil {
		log.Fatal("failed to start http server", zap.Error(err))
	}

	<-ctx.Done()

	log.Info("http server stopped")
}
