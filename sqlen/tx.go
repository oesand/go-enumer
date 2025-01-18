package sqlen

import (
	"context"
	"database/sql"
)

var DefaultTxOptions = sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

func ExecuteTx(ctx context.Context, opts *sql.TxOptions, execFunc func(ctx context.Context) error) (err error) {
	if ctx == nil {
		panic("ctx cannot be nit")
	}

	var tx *sql.Tx
	var nested bool
	if tx, nested = ctx.Value(defaultCtxExecutorKey).(*sql.Tx); !nested {
		if opts == nil {
			opts = &DefaultTxOptions
		}
		if db, ok := ctx.Value(defaultCtxExecutorKey).(*sql.DB); ok {
			tx, err = db.BeginTx(ctx, opts)
			ctx = WithTx(ctx, tx)
		} else {
			panic("db not passed")
		}
	}

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err = execFunc(ctx); err != nil {
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
