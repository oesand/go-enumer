package sqlen

import "database/sql"

type Repo[T any, TPK comparable] interface {
	DB() *sql.DB
	Table() string
	PK() string
	Fields() []string

	Formatter() ParamFormatter
	Template() (*T, []any)
	Extract(*T) []any
}
