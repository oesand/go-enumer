package sqlen

import (
	"context"
	"database/sql"
)

const (
	defaultCtxTransactionKey = "go-enumer@key/tx"
)

func WrapTx(ctx context.Context, tx *sql.Tx) context.Context {
	if ctx == nil || tx == nil {
		panic("parameters cannot be nit")
	}
	return context.WithValue(ctx, defaultCtxTransactionKey, tx)
}

func UnwrapTx(ctx context.Context) *sql.Tx {
	if ctx == nil {
		panic("ctx cannot be nit")
	}
	return ctx.Value(defaultCtxTransactionKey).(*sql.Tx)
}

func DefaultExecutor[T any](repo Repo[T], ctx context.Context) ContextExecutor {
	if tx := UnwrapTx(ctx); tx != nil {
		return tx
	}
	return repo.DB()
}
