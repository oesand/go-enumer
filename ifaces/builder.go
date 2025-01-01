package ifaces

type Builder[T any] interface {
	Build() T
}

type BuilderQuery[T any] interface {
	Builder[T]
	QueryZip() ([]string, []any)
}

type QueryableModel interface {
	GetFieldValue(fieldName string) any
	QueryZip() ([]string, []any)
}
