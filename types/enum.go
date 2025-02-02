package types

import (
	"fmt"
	"reflect"
)

type Enum interface {
	~string | ~int | ~int32 | ~int64
	IsValid() bool
	String() string
}

type CombinedEnum[T Enum] interface {
	Enum
	Has(en T) bool
}

func EnsureValid[T Enum](en T) T {
	if !en.IsValid() {
		panic(fmt.Sprintf("invalid value '%s' for %s enum", en, reflect.TypeFor[T]().Name()))
	}
	return en
}
