package app

import (
	"fmt"
	"log/slog"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/serj213/bookService/internal/app/grpc"
	gprcapp "github.com/serj213/bookService/internal/app/grpc"
	pgrepo "github.com/serj213/bookService/internal/repository/pgRepo"
	"github.com/serj213/bookService/internal/services/book"
	"github.com/serj213/bookService/pkg/pg"
)

type App struct {
	GRPCServer *grpc.App
}

func New(
	log *slog.Logger,
	dsn string,
	migrationPath string,
	port int,
) *App {
	pgDb, err := pg.Deal(dsn)
	if err != nil {
		log.Error(fmt.Sprintf("failed to connect to postgres: %v", err))
		panic(err)
	}

	log.Info("postgres connect succesfully")

	trManager := manager.Must(trmpgx.NewDefaultFactory(pgDb))

	err = migrations(migrationPath, dsn)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		panic(err)
	}

	log.Info("migrations successfuly")

	bookRepo := pgrepo.NewBookRepo(pgDb)
	bookService := book.NewBookService(log, bookRepo, trManager)

	grpcApp := gprcapp.New(log, bookService, port)

	return &App{
		GRPCServer: grpcApp,
	}

}

func migrations(migrationsPath string, dsn string) error {

	if migrationsPath == "" {
		return fmt.Errorf("migrations path empty")
	}

	if dsn == "" {
		return fmt.Errorf("dsn empty")
	}

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("failed migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed connect migration: %w", err)
	}

	return nil
}
