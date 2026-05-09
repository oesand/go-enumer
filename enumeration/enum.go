// Package enumeration provides a runtime enumeration registry for values
// that implement a custom enum key. It supports declaring values by type
// and looking them up by that key.
package enumeration

import (
	"fmt"
	"iter"
	"reflect"

	"github.com/oesand/go-enumer/internal"
)

// enums stores registered enum values keyed by enum type.
var enums = map[internal.ShadowType]map[any]any{}

// Enumeration is the interface used to register enum values.
// The EnumerationKey method provides a stable, comparable key for each value.
type Enumeration interface {
	EnumerationKey() any
}

// Declare registers a new enum value for the given enum type.
// It panics if a value with the same key has already been declared.
func Declare[Enum Enumeration](value Enum) Enum {
	var enum internal.Type[Enum]
	entries, ok := enums[enum]
	if !ok {
		entries = make(map[any]any)
		enums[enum] = entries
	}

	valueKey := value.EnumerationKey()
	if _, found := entries[valueKey]; found {
		panic(fmt.Sprintf("duplicate enum value key: %v", valueKey))
	}

	entries[valueKey] = value
	return value
}

func getEntries[Enum Enumeration]() map[any]any {
	var enum internal.Type[Enum]
	entries, ok := enums[enum]
	if !ok {
		panic(fmt.Sprintf("unknown enumeration: %s", reflect.TypeFor[Enum]()))
	}
	return entries
}

// FromKey looks up an enum value by its key and returns it if found.
// It returns the zero value and false when no entry exists for the provided key.
func FromKey[Enum Enumeration](key any) (Enum, bool) {
	entries := getEntries[Enum]()
	value, ok := entries[key]
	if !ok {
		var zero Enum
		return zero, false
	}
	return value.(Enum), true
}

// FindValue returns the first enum value matching the provided predicate.
// If no value matches the matcher, it returns the zero value and false.
func FindValue[Enum Enumeration](matcher func(Enum) bool) (Enum, bool) {
	entries := getEntries[Enum]()
	for _, value := range entries {
		val := value.(Enum)
		if matcher(val) {
			return val, true
		}
	}
	var zero Enum
	return zero, false
}

// Values returns an iterator over all declared enum values for type Enum.
func Values[Enum Enumeration]() iter.Seq[Enum] {
	entries := getEntries[Enum]()
	return func(yield func(Enum) bool) {
		for _, value := range entries {
			if !yield(value.(Enum)) {
				break
			}
		}
	}
}
