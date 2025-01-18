package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/oesand/go-enumer/sqlen"
)

func AcqTxAdvisoryLock[T any](repo sqlen.Repo[T], ctx context.Context, key int64) error {
	return acqTxAdvisoryLock(repo, ctx, key)
}

func AcqTxAdvisoryLock2[T any](repo sqlen.Repo[T], ctx context.Context, one int, two int) error {
	return acqTxAdvisoryLock(repo, ctx, one, two)
}

func AcqTxAdvisoryLockString[T any](repo sqlen.Repo[T], ctx context.Context, key string) error {
	return acqTxAdvisoryLock(repo, ctx, hashString(key))
}

func AcqTxAdvisoryLock2String[T any](repo sqlen.Repo[T], ctx context.Context, one string, two string) error {
	return acqTxAdvisoryLock(repo, ctx, hashString(one), hashString(two))
}

func acqTxAdvisoryLock[T any](repo sqlen.Repo[T], ctx context.Context, keys ...any) error {
	exec := sqlen.GetExecutor(ctx)
	if _, ok := ctx.Value(exec).(*sql.Tx); !ok {
		panic("acquire advisory lock should use only in transaction")
	}

	var query string
	switch len(keys) {
	case 1:
		query = fmt.Sprintf("SELECT pg_advisory_xact_lock(%s)",
			repo.Formatter().Format(1))
	case 2:
		query = fmt.Sprintf("SELECT pg_advisory_xact_lock(%s, %s)",
			repo.Formatter().Format(1), repo.Formatter().Format(2))
	default:
		panic("invalid count keys")
	}

	_, err := exec.ExecContext(ctx, query, keys...)
	return err
}
