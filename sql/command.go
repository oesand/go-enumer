package sqlen

import (
	"context"
	"fmt"
	"github.com/oesand/go-enumer/types"
	"strings"
)

func ExecCreate[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, model T) (TPK, error) {
	return execCreate(repo, ctx, model, true)
}

func ExecCreateNext[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, model T) (TPK, error) {
	return execCreate(repo, ctx, model, false)
}

func execCreate[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, model T, includePk bool) (TPK, error) {
	values := model.QueryValues(false)
	if len(values) != len(repo.fields) {
		panic(fmt.Sprintf("Count fields mismatch, values: %d, fields %d", len(values), len(repo.fields)))
	}

	var fields []string
	if includePk {
		fields = repo.fields
	} else {
		fields = repo.fields[1:]
		values = values[1:]
	}

	paramsString := repeatParamPlaceholders(repo.formatter, 0, len(fields))
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		repo.table, strings.Join(fields, ", "), paramsString, repo.PK())

	fmt.Printf("Query: %s, Values: %v\n", query, values)

	var pk TPK
	err := repo.db.QueryRowContext(ctx, query, values...).Scan(&pk)
	return pk, err
}

func ExecUpdate[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, pk TPK, builder types.QueryBuilder) error {
	names, values := builder.QueryValues(repo.caseType)
	if len(names) != len(values) {
		panic(fmt.Sprintf("Count fields mismatch, names: %d, values %d", len(names), len(values)))
	}

	var setString strings.Builder
	for i, name := range names {
		if i > 0 {
			setString.WriteString(", ")
		}
		setString.WriteString(fmt.Sprintf("%s = ?", name))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		repo.table, setString.String(), repo.PK())

	values = append(values, pk)
	_, err := repo.db.ExecContext(ctx, query, values...)
	return err
}

func ExecDelete[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, pk TPK) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", repo.table, repo.PK())
	_, err := repo.db.ExecContext(ctx, query, pk)
	return err
}

func QuerySelectByPK[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, pk TPK) (T, error) {
	return QuerySelectSingle(repo, ctx, fmt.Sprintf("%s = ?", repo.PK()), pk)
}

func QuerySelectSingle[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, whereStatement string, values ...any) (T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(repo.fields, ", "), repo.table, whereStatement)

	var model T
	row := repo.db.QueryRowContext(ctx, query, values...)
	pointers := model.QueryValues(true)
	err := row.Scan(pointers...)
	return model, err
}

func QuerySelectMany[T types.QueryModel, TPK comparable](repo *Repo[T, TPK], ctx context.Context, whereStatement string, values ...any) ([]T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(repo.fields, ", "), repo.table, whereStatement)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []T
	for rows.Next() {
		var model T
		pointers := model.QueryValues(true)
		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}
