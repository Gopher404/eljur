package transaction

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoTx = errors.New("tx is not in the context values")
)

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db}
}

type TxManager struct {
	db *sql.DB
}

func (m *TxManager) Begin(ctx context.Context) (context.Context, error) {
	/* if tx := ExtractTx(ctx); tx != nil {
		return ctx, nil
	} */

	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}

	return CtxWidthTx(ctx, tx), nil
}

func (*TxManager) Rollback(ctx context.Context) error {
	tx := ExtractTx(ctx)
	if tx == nil {
		return ErrNoTx
	}
	return tx.Rollback()
}

func (*TxManager) Commit(ctx context.Context) error {
	tx := ExtractTx(ctx)
	if tx == nil {
		return ErrNoTx
	}
	return tx.Commit()
}
