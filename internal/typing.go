// Package internal provides type utilities for the enum package.
package internal

import "reflect"

// Type is a generic type wrapper for reflection.
type Type[T any] struct{}

// Type returns the reflect.Type for T.
func (Type[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

// Zero returns the zero value of T as any.
func (Type[T]) Zero() any {
	var zero T
	return zero
}

// ShadowType is an interface for type introspection.
type ShadowType interface {
	Type() reflect.Type
	Zero() any
}
