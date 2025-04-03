package main

import (
	"log/slog"
	"os"

	"github.com/serj213/bookService/internal/app"
	"github.com/serj213/bookService/internal/config"
)

const (
	local = "local"
	dev   = "dev"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	log := setupLogger(cfg.Env)

	log = log.With(slog.String("service", "bookService"))

	log.Info("logger enabled")

	application := app.New(log, cfg.Dsn, cfg.MigrationPath, cfg.Grpc.Port)

	application.GRPCServer.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
