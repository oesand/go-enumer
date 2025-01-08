package sqlen

import (
	"context"
	"errors"
	"fmt"
	"github.com/oesand/go-enumer/types"
	"strings"
)

func ExecCreate[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, model *T) error {
	values := repo.Extract(model)
	fields := repo.Fields()
	if len(values) != len(fields) {
		panic(fmt.Sprintf("Count fields mismatch, values: %d, fields %d", len(values), len(fields)))
	}

	var valueString strings.Builder
	for i := 0; i < len(fields); i++ {
		if i > 0 {
			valueString.WriteString(", ")
		}
		valueString.WriteString(repo.Formatter().Format(i + 1))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		repo.Table(), strings.Join(fields, ", "), valueString.String())

	_, err := repo.DB().ExecContext(ctx, query, values...)
	return err
}

func ExecCreateNext[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, model *T) (TPK, error) {
	values := repo.Extract(model)
	fields := repo.Fields()
	if len(values) != len(fields) {
		panic(fmt.Sprintf("Count fields mismatch, values: %d, fields %d", len(values), len(fields)))
	}

	fields = fields[1:]
	values = values[1:]

	var valueString strings.Builder
	for i := 0; i < len(fields); i++ {
		if i > 0 {
			valueString.WriteString(", ")
		}
		valueString.WriteString(repo.Formatter().Format(i + 1))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		repo.Table(), strings.Join(fields, ", "), valueString.String(), repo.PK())

	var pk TPK
	err := repo.DB().QueryRowContext(ctx, query, values...).Scan(&pk)
	return pk, err
}

func ExecUpdate[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, pk TPK, builder types.QueryBuilder[T]) error {
	names, values := builder.QueryValues()
	if len(names) == 0 {
		return errors.New("set values cannot be empty")
	}
	if len(names) != len(values) {
		panic(fmt.Sprintf("Count fields mismatch, names: %d, values %d", len(names), len(values)))
	}

	var setString strings.Builder
	for i, name := range names {
		if i > 0 {
			setString.WriteString(", ")
		}
		setString.WriteString(fmt.Sprintf("%s = %s", name, repo.Formatter().Format(i+1)))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s",
		repo.Table(), setString.String(), repo.PK(), repo.Formatter().Format(len(names)+1))

	values = append(values, pk)
	_, err := repo.DB().ExecContext(ctx, query, values...)
	return err
}

func ExecDelete[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, pk TPK) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		repo.Table(), repo.PK(), repo.Formatter().Format(1))

	_, err := repo.DB().ExecContext(ctx, query, pk)
	return err
}

func QuerySelectByPK[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, pk TPK) (*T, error) {
	return QuerySelectSingle[T, TPK](repo, ctx, fmt.Sprintf("%s = %s", repo.PK(), repo.Formatter().Format(1)), pk)
}

func QuerySelectSingle[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, whereStatement string, values ...any) (*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(repo.Fields(), ", "), repo.Table(), whereStatement)

	model, pointers := repo.Template()
	err := repo.DB().QueryRowContext(ctx, query, values...).Scan(pointers...)
	return model, err
}

func QuerySelectMany[T any, TPK comparable](repo Repo[T, TPK], ctx context.Context, whereStatement string, values ...any) ([]*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(repo.Fields(), ", "), repo.Table(), whereStatement)

	rows, err := repo.DB().QueryContext(ctx, query, values...)
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
