package sqlen

import (
	"context"
	"database/sql"
)

var DefaultTxOptions = sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

func ExecuteTx[T any](repo Repo[T], ctx context.Context, opts *sql.TxOptions, execFunc func(repo Repo[T], ctx context.Context) error) (err error) {
	if ctx == nil {
		panic("ctx cannot be nit")
	}

	var tx *sql.Tx
	var nested bool
	if tx, nested = ctx.Value(defaultCtxTransactionKey).(*sql.Tx); !nested {
		if opts == nil {
			opts = &DefaultTxOptions
		}
		tx, err = repo.DB().BeginTx(ctx, opts)
		ctx = WrapTx(ctx, tx)
	}

	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := execFunc(repo, ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	select {
	case <-ctx.Done():
		_ = tx.Rollback()
		return ctx.Err()
	default:
	}

	if !nested {
		return tx.Commit()
	}
	return
}
