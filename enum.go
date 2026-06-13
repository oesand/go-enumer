// Package enum provides a runtime enumeration system for Go.
// It allows declaring named constants of supported types and provides
// utilities to validate, convert, and iterate over enum values.
package enum

import (
	"fmt"
	"iter"
	"reflect"

	"github.com/oesand/go-enumer/internal"
)

var enums = map[internal.ShadowType]*enumEntries{}
var enumsReflect = map[reflect.Type]*enumEntries{}

type enumEntries struct {
	byName  map[string]any
	byValue map[any]string
}

// SupportedTypes defines the types that can be used for enum values.
type SupportedTypes interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Declare registers a new enum value without name.
// It panics if the value is already declared for this enum type.
// Returns the value for convenience.
func Declare[T SupportedTypes](value T) T {
	return DeclareNamed("", value)
}

// DeclareNamed registers a new enum value with the given name.
// It panics if the name or value is already declared for this enum type.
// Returns the value for convenience.
func DeclareNamed[Enum SupportedTypes](name string, value Enum) Enum {
	var enumKey internal.Type[Enum]
	entries, ok := enums[enumKey]
	if !ok {
		entries = &enumEntries{
			byName:  make(map[string]any),
			byValue: make(map[any]string),
		}
		enums[enumKey] = entries
		enumsReflect[enumKey.Type()] = entries
	} else {
		if name != "" {
			if _, found := entries.byName[name]; found {
				panic(fmt.Sprintf("duplicate enum name: %s", name))
			}
		}

		if _, found := entries.byValue[value]; found {
			panic(fmt.Sprintf("duplicate enum value: %v", value))
		}
	}

	entries.byName[name] = value
	entries.byValue[value] = name

	return value
}

func getEntries[Enum SupportedTypes]() *enumEntries {
	var enum internal.Type[Enum]
	entries, ok := enums[enum]
	if !ok {
		panic(fmt.Sprintf("enumer: unknown enum: %s", reflect.TypeFor[Enum]()))
	}
	return entries
}

func getEntriesR(t reflect.Type) *enumEntries {
	entries, ok := enumsReflect[t]
	if !ok {
		panic(fmt.Sprintf("enumer: unknown enum: %s", t.String()))
	}
	return entries
}

// IsValid checks if the given value is a declared enum value for type T.
func IsValid[Enum SupportedTypes](value Enum) bool {
	entries := getEntries[Enum]()
	_, ok := entries.byValue[value]
	return ok
}

// IsEnum checks if the given type is a declared enum type.
func IsEnum[Enum SupportedTypes]() bool {
	var enum internal.Type[Enum]
	_, ok := enums[enum]
	return ok
}

// IsEnumR checks if the given type is a declared enum type.
func IsEnumR(t reflect.Type) bool {
	_, ok := enumsReflect[t]
	return ok
}

// FromName returns the enum value for the given name, or the zero value and false if not found.
func FromName[Enum SupportedTypes](name string) (Enum, bool) {
	entries := getEntries[Enum]()
	value, ok := entries.byName[name]
	if !ok {
		var zero Enum
		return zero, false
	}
	return value.(Enum), true
}

// Names returns an iterator over all declared enum names for type T.
func Names[Enum SupportedTypes]() iter.Seq[string] {
	entries := getEntries[Enum]()
	return func(yield func(string) bool) {
		for name := range entries.byName {
			if !yield(name) {
				break
			}
		}
	}
}

// Values returns an iterator over all declared enum values for type T.
func Values[Enum SupportedTypes]() iter.Seq[Enum] {
	entries := getEntries[Enum]()
	return func(yield func(Enum) bool) {
		for value := range entries.byValue {
			if !yield(value.(Enum)) {
				break
			}
		}
	}
}

// Entries returns an iterator over all enum entries(name, value) for type T.
func Entries[Enum SupportedTypes]() iter.Seq2[string, Enum] {
	entries := getEntries[Enum]()
	return func(yield func(string, Enum) bool) {
		for value, name := range entries.byValue {
			if !yield(name, value.(Enum)) {
				break
			}
		}
	}
}

// EntriesR returns an iterator over all enum entries(name, value) for type T.
func EntriesR(t reflect.Type) iter.Seq2[string, any] {
	entries := getEntriesR(t)
	return func(yield func(string, any) bool) {
		for value, name := range entries.byValue {
			if !yield(name, value) {
				break
			}
		}
	}
}
