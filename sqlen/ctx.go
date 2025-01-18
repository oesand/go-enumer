package sqlen

import (
	"context"
	"database/sql"
)

const defaultCtxExecutorKey = "go-enumer@key/exec"

func WithDB(ctx context.Context, db *sql.DB) context.Context {
	if db == nil {
		panic("db cannot be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	} else if val := ctx.Value(defaultCtxExecutorKey); val != nil {
		return ctx
	}
	return context.WithValue(ctx, defaultCtxExecutorKey, db)
}

func GetExecutor(ctx context.Context) ContextExecutor {
	if ctx == nil {
		panic("ctx cannot be nil")
	}
	if exec, ok := ctx.Value(defaultCtxExecutorKey).(ContextExecutor); ok && exec != nil {
		return exec
	}
	panic("executor not passed in ctx")
}
