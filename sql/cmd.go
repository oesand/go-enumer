package sqlen

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/oesand/go-enumer/types"
	"strings"
)

func ExecCreate[T any](repo Repo[T], ctx context.Context, model *T) error {
	values := repo.Extract(model)
	fields := repo.Fields()
	if len(values) != len(fields) {
		panic(fmt.Sprintf("count fields mismatch, values: %d, fields %d", len(values), len(fields)))
	}

	var valueString strings.Builder
	for i, value := range values {
		if i > 0 {
			valueString.WriteString(", ")
		}
		valueString.WriteString(repo.Formatter().Format(i + 1))

		if repo.Formatter().Type() == NamedFormatterType {
			if _, ok := values[i].(sql.NamedArg); !ok {
				values[i] = namedParam(i+1, value)
			}
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		repo.Table(), strings.Join(fields, ", "), valueString.String())

	_, err := DefaultExecutor[T](repo, ctx).ExecContext(ctx, query, values...)
	return err
}

func ExecCreateNext[T any](repo Repo[T], ctx context.Context, model *T) (int64, error) {
	values := repo.Extract(model)
	fields := repo.Fields()
	if len(values) != len(fields) {
		panic(fmt.Sprintf("count fields mismatch, values: %d, fields %d", len(values), len(fields)))
	}

	fields = fields[1:]
	values = values[1:]

	var valueString strings.Builder
	for i, value := range values {
		if i > 0 {
			valueString.WriteString(", ")
		}
		valueString.WriteString(repo.Formatter().Format(i + 1))

		if repo.Formatter().Type() == NamedFormatterType {
			if _, ok := values[i].(sql.NamedArg); !ok {
				values[i] = namedParam(i+1, value)
			}
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		repo.Table(), strings.Join(fields, ", "), valueString.String(), repo.PK())

	var pk int64
	err := DefaultExecutor[T](repo, ctx).QueryRowContext(ctx, query, values...).Scan(&pk)
	return pk, err
}

func ExecUpdate[T any, TPK comparable](repo Repo[T], ctx context.Context, pk TPK, builder types.QueryBuilder[T]) error {
	names, values := builder.QueryValues()
	if len(names) == 0 {
		return errors.New("set values cannot be empty")
	}
	if len(names) != len(values) {
		panic(fmt.Sprintf("count fields mismatch, names: %d, values %d", len(names), len(values)))
	}

	var setString strings.Builder
	for i, name := range names {
		if i > 0 {
			setString.WriteString(", ")
		}
		setString.WriteString(fmt.Sprintf("%s = %s", name, repo.Formatter().Format(i+1)))

		if repo.Formatter().Type() == NamedFormatterType {
			if _, ok := values[i].(sql.NamedArg); !ok {
				values[i] = namedParam(i+1, values[i])
			}
		}
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s",
		repo.Table(), setString.String(), repo.PK(), repo.Formatter().Format(len(names)+1))

	if repo.Formatter().Type() == NamedFormatterType {
		values = append(values, namedParam(len(values)+1, pk))
	} else {
		values = append(values, pk)
	}
	_, err := DefaultExecutor[T](repo, ctx).ExecContext(ctx, query, values...)
	return err
}

func ExecDelete[T any, TPK comparable](repo Repo[T], ctx context.Context, pk TPK) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		repo.Table(), repo.PK(), repo.Formatter().Format(1))

	var value any
	if repo.Formatter().Type() == NamedFormatterType {
		value = namedParam(1, value)
	} else {
		value = pk
	}

	_, err := DefaultExecutor[T](repo, ctx).ExecContext(ctx, query, value)
	return err
}

func QuerySelectByPK[T any, TPK comparable](repo Repo[T], ctx context.Context, pk TPK) (*T, error) {
	return QuerySelectSingle[T](repo, ctx, fmt.Sprintf("%s = %s", repo.PK(), repo.Formatter().Format(1)), pk)
}

func QuerySelectSingle[T any](repo Repo[T], ctx context.Context, whereStatement string, values ...any) (*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(repo.Fields(), ", "), repo.Table(), whereStatement)

	if repo.Formatter().Type() == NamedFormatterType {
		for i, value := range values {
			if _, ok := values[i].(sql.NamedArg); !ok {
				values[i] = namedParam(i+1, value)
			}
		}
	}

	model, pointers := repo.Template()
	err := DefaultExecutor[T](repo, ctx).QueryRowContext(ctx, query, values...).Scan(pointers...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model, err
}

func QuerySelectMany[T any](repo Repo[T], ctx context.Context, whereStatement string, values ...any) ([]*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(repo.Fields(), ", "), repo.Table(), whereStatement)

	if repo.Formatter().Type() == NamedFormatterType {
		for i, value := range values {
			if _, ok := values[i].(sql.NamedArg); !ok {
				values[i] = namedParam(i+1, value)
			}
		}
	}

	rows, err := DefaultExecutor[T](repo, ctx).QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*T
	for rows.Next() {
		model, pointers := repo.Template()
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func QueryExists[T any](repo Repo[T], ctx context.Context, whereStatement string, values ...any) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", repo.Table(), whereStatement)

	if repo.Formatter().Type() == NamedFormatterType {
		for i, value := range values {
			if _, ok := values[i].(sql.NamedArg); !ok {
				values[i] = namedParam(i+1, value)
			}
		}
	}

	var exists bool
	err := DefaultExecutor[T](repo, ctx).QueryRowContext(ctx, query, values...).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
