package types

import (
	"github.com/oesand/go-enumer/cases"
)

type Builder[T any] interface {
	Build() T
}

type QueryModel interface {
	GetFieldValue(fieldName string) any
	QueryValues(ptr bool) []any
}

type QueryBuilder interface {
	QueryValues(caseType cases.CaseType) ([]string, []any)
}
