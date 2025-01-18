package sqlen

import (
	"context"
	"database/sql"
)

const defaultCtxExecutorKey = "go-enumer@key/exec"

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	if tx == nil {
		panic("tx cannot be nil")
	}
	return withExecutor(ctx, tx)
}

func WithDB(ctx context.Context, db *sql.DB) context.Context {
	if db == nil {
		panic("db cannot be nil")
	}
	return withExecutor(ctx, db)
}

func withExecutor(ctx context.Context, exec ContextExecutor) context.Context {
	if ctx == nil {
		ctx = context.Background()
	} else if val := ctx.Value(defaultCtxExecutorKey); val != nil {
		return ctx
	}
	return context.WithValue(ctx, defaultCtxExecutorKey, exec)
}

func GetExecutor(ctx context.Context) ContextExecutor {
	if ctx == nil {
		panic("ctx cannot be nil")
	}
	if exec, ok := ctx.Value(defaultCtxExecutorKey).(ContextExecutor); ok && exec != nil {
		return exec
	}
	panic("executor not passed")
}
