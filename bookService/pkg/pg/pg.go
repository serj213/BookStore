package pg

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgDb struct {
	*pgxpool.Pool
}

func Deal(dsn string) (*PgDb, error) {

	if dsn == "" {
		return nil, fmt.Errorf("dsn empty")
	}

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, fmt.Errorf("failed connect db: %w", err)
	}

	return &PgDb{
		dbpool,
	}, nil

}
