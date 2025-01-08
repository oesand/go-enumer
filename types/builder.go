package types

type Builder[T any] interface {
	Build() T
}

type QueryBuilder[T any] interface {
	Builder[T]
	QueryValues() ([]string, []any)
}
