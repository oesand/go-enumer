package sqlen

import (
	"database/sql"
	"errors"
	"github.com/oesand/go-enumer/cases"
	"github.com/oesand/go-enumer/types"
)

func New[T types.QueryModel, TPK comparable](driverName, dataSource, table string, caseType cases.CaseType, fields []string) (*Repo[T, TPK], error) {
	if !caseType.IsValid() {
		return nil, errors.New("invalid case type")
	}

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	if caseType != cases.NoChange {
		for i, field := range fields {
			fields[i] = caseType.From(field)
		}
	}
	formatter := QuestionFormatter
	if frm, has := defaultFormatters[driverName]; has {
		formatter = frm
	}
	return &Repo[T, TPK]{
		db:        db,
		table:     table,
		fields:    fields,
		caseType:  caseType,
		formatter: formatter,
	}, nil
}

type Repo[T any, TPK comparable] struct {
	db        *sql.DB
	table     string
	fields    []string
	caseType  cases.CaseType
	formatter ParamFormatter
}

func (r *Repo[T, TPK]) DB() *sql.DB {
	return r.db
}

func (r *Repo[T, TPK]) Table() string {
	return r.table
}

func (r *Repo[T, TPK]) PK() string {
	return r.fields[0]
}

func (r *Repo[T, TPK]) Fields() []string {
	temp := make([]string, len(r.fields))
	copy(temp, r.fields)
	return temp
}

func (r *Repo[T, TPK]) ParamFormatter(formatter ParamFormatter) *Repo[T, TPK] {
	if formatter == nil {
		panic("formatter cannot be nil")
	}
	r.formatter = formatter
	return r
}
