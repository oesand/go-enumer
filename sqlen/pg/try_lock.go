package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/oesand/go-enumer/sqlen"
)

func TryAcqTxAdvisoryLock[T any](repo sqlen.Repo[T], ctx context.Context, key int64) (bool, error) {
	return tryAcqTxAdvisoryLock(repo, ctx, key)
}

func TryAcqTxAdvisoryLock2[T any](repo sqlen.Repo[T], ctx context.Context, one int, two int) (bool, error) {
	return tryAcqTxAdvisoryLock(repo, ctx, one, two)
}

func TryAcqTxAdvisoryLockString[T any](repo sqlen.Repo[T], ctx context.Context, key string) (bool, error) {
	return tryAcqTxAdvisoryLock(repo, ctx, hashString(key))
}

func TryAcqTxAdvisoryLock2String[T any](repo sqlen.Repo[T], ctx context.Context, one string, two string) (bool, error) {
	return tryAcqTxAdvisoryLock(repo, ctx, hashString(one), hashString(two))
}

func tryAcqTxAdvisoryLock[T any](repo sqlen.Repo[T], ctx context.Context, keys ...any) (bool, error) {
	exec := sqlen.GetExecutor(ctx)
	if _, ok := ctx.Value(exec).(*sql.Tx); !ok {
		panic("acquire advisory lock should use only in transaction")
	}

	var query string
	switch len(keys) {
	case 1:
		query = fmt.Sprintf("SELECT pg_try_advisory_xact_lock(%s)",
			repo.Formatter().Format(1))
	case 2:
		query = fmt.Sprintf("SELECT pg_try_advisory_xact_lock(%s, %s)",
			repo.Formatter().Format(1), repo.Formatter().Format(2))
	default:
		panic("invalid count keys")
	}

	var result bool
	err := exec.QueryRowContext(ctx, query, keys...).Scan(&result)
	return result, err
}
