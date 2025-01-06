package types

type Enum interface {
	~string | ~int | ~int32 | ~int64
	IsValid() bool
	String() string
}

type CombinedEnum[T Enum] interface {
	Enum
	Has(en T) bool
}
