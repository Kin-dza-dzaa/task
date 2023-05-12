package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Kin-dza-dzaa/task/config"
	_ "github.com/Kin-dza-dzaa/task/docs"
	"github.com/Kin-dza-dzaa/task/internal/repository/currencyapi"
	"github.com/Kin-dza-dzaa/task/internal/repository/postgresql"
	"github.com/Kin-dza-dzaa/task/internal/service"
	"github.com/Kin-dza-dzaa/task/internal/transport/http/server"
	"github.com/Kin-dza-dzaa/task/internal/transport/http/v1/rest"
	"github.com/Kin-dza-dzaa/task/pkg/logger"
	"github.com/Kin-dza-dzaa/task/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"golang.org/x/exp/slog"
)

func main() {
	ctx := context.Background()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("main - config.ParseConfig: %w", err))
	}
	log.Fatal(run(ctx, cfg))
}

func run(ctx context.Context, cfg config.Config) error {
	l := logger.New(slog.Level(cfg.Logger.Level))
	client, err := postgres.NewClient(ctx, cfg.PG.Address, cfg.PG.MaxPoolSize, pgx.LogLevel(cfg.PG.LogLevel), logger.NewPGXLogger(l))
	if err != nil {
		return fmt.Errorf("main - run - run - postgres.NewClient: %w", err)
	}

	// Layers.
	// Adapters.
	currencyAPI := currencyapi.New(cfg.CurrencyAPI.Address)
	currencyRepo := postgresql.New(client)
	// Services.
	s := service.New(currencyRepo, currencyAPI, l)
	// Ports.
	h := rest.New(l, s)

	// HTTP server.
	router := chi.NewMux()
	h.Register(router, cfg)

	srv := server.New(cfg, l, router)
	<-srv.Start(ctx)

	return nil
}
