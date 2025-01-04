package ifaces

import "github.com/oesand/go-enumer/cases"

type Builder[T any] interface {
	Build() T
}

type BuilderQuery[T any] interface {
	Builder[T]
	QueryValues(caseType cases.CaseType) ([]string, []any)
}

type QueryableModel interface {
	GetFieldValue(fieldName string) any
	QueryValues(caseType cases.CaseType, ptr bool) ([]string, []any)
}
