package ifaces

import "github.com/oesand/go-enumer/cases"

type Builder[T any] interface {
	Build() T
}

type BuilderQuery[T any] interface {
	Builder[T]
	QueryZip(caseType cases.CaseType) ([]string, []any)
}

type QueryableModel interface {
	GetFieldValue(fieldName string) any
	QueryZip(caseType cases.CaseType) ([]string, []any)
	QueryScan(rowScan RowScanner) error
}
