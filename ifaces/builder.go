package ifaces

import (
	"github.com/oesand/go-enumer/cases"
)

type Builder[T any] interface {
	Build() T
}

type QueryModel interface {
	GetFieldValue(fieldName string) any
	QueryValues(caseType cases.CaseType, ptr bool) ([]string, []any)
}

type QueryBuilder[T any] interface {
	Builder[T]
	QueryValues(caseType cases.CaseType) ([]string, []any)
}
