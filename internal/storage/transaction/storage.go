package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

type TxStorage struct {
	DB *sql.DB
}

func (s *TxStorage) Begin(ctx context.Context) (context.Context, error) {
	/*if tx := storage.ExtractTx(ctx); tx != nil {
		return ctx, nil
	}*/

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	return CtxWidthTx(ctx, tx), nil
}

func (s *TxStorage) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	tx := ExtractTx(ctx)
	if tx != nil {
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			errRB := tx.Rollback()
			return nil, fmt.Errorf("%e, rollback error: %e", err, errRB)
		}
		return rows, nil
	}
	return s.DB.QueryContext(ctx, query, args...)
}

func (s *TxStorage) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	tx := ExtractTx(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return s.DB.QueryRowContext(ctx, query, args...)
}

func (s *TxStorage) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx := ExtractTx(ctx)
	if tx != nil {
		res, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			errRB := tx.Rollback()
			return nil, fmt.Errorf("%e, rollback error: %e", err, errRB)
		}
		return res, nil
	}
	return s.DB.ExecContext(ctx, query, args...)
}
