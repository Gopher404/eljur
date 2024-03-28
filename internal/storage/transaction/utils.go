package transaction

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoTx = errors.New("tx is not in the context values")
)

type txKey struct{}

func CtxWidthTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func ExtractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func Rollback(ctx context.Context) error {
	tx := ExtractTx(ctx)
	if tx == nil {
		return ErrNoTx
	}
	return tx.Rollback()
}

func Commit(ctx context.Context) error {
	tx := ExtractTx(ctx)
	if tx == nil {
		return ErrNoTx
	}
	return tx.Commit()
}
