package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PgTransactionFn func(pgTX pgx.Tx) error

func HandlePgTransaction(ctx context.Context, pgTxFn PgTransactionFn, db *PgDb) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed begin transaction: %w", err)
	}

	errFn := pgTxFn(tx)

	if errFn != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("failed rollback transaction: %w", err)
		}

		return fmt.Errorf("failed executing transaction: %w", errFn)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed commit transaction: %w", err)
	}

	return nil

}
