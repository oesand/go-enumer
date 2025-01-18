package sqlen

import (
	"context"
	"database/sql"
)

type Repo[T any] interface {
	Table() string
	PK() string
	Fields() []string
	Formatter() ParamFormatter

	Template() (*T, []any)
	Extract(*T) []any
}

type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
