package enumeration

import (
	"reflect"
	"slices"
	"testing"

	"github.com/oesand/go-enumer/internal"
)

func cleanValues[T any]() {
	delete(enums, internal.Type[T]{})
}

type TestEnum struct {
	Key string
}

func (e *TestEnum) EnumerationKey() any {
	return e.Key
}

func TestDeclareDuplicateKey(t *testing.T) {
	cleanValues[*TestEnum]()

	val := Declare(&TestEnum{
		Key: "foo",
	})
	if val.Key != "foo" {
		t.Fatalf("expected foo, got %v", val.Key)
	}

	defer func() {
		if r := recover(); !reflect.DeepEqual(r, "duplicate enum value key: foo") {
			t.Fatalf("expected duplicate key panic, got %v", r)
		}
	}()
	Declare(&TestEnum{
		Key: "foo",
	})
}

func TestFromKey(t *testing.T) {
	cleanValues[*TestEnum]()

	Declare(&TestEnum{Key: "foo"})

	val, ok := FromKey[*TestEnum]("foo")
	if !ok || val.Key != "foo" {
		t.Fatalf("expected foo, got %v, ok=%v", val, ok)
	}

	_, ok = FromKey[*TestEnum]("123")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestFindValue(t *testing.T) {
	cleanValues[*TestEnum]()

	Declare(&TestEnum{Key: "alpha"})
	Declare(&TestEnum{Key: "beta"})

	val, ok := FindValue(func(v *TestEnum) bool {
		return v.Key == "beta"
	})
	if !ok || val.Key != "beta" {
		t.Fatalf("expected beta, got %v, ok=%v", val, ok)
	}
}

func TestValues(t *testing.T) {
	cleanValues[*TestEnum]()

	expected := []*TestEnum{
		{Key: "alpha"},
		{Key: "beta"},
	}

	for _, enum := range expected {
		Declare(enum)
	}

	values := slices.Collect(Values[*TestEnum]())
	if !slices.Equal(values, expected) {
		t.Fatalf("expected %v, got %v", expected, values)
	}
}
