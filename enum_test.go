package enum

import (
	"reflect"
	"slices"
	"testing"

	"github.com/oesand/go-enumer/internal"
)

func cleanValues[T any]() {
	delete(enums, internal.Type[T]{})
}

type TestEnum int

type StringEnum string

func TestDeclareNamed(t *testing.T) {
	cleanValues[TestEnum]()

	// Test normal declaration
	val := DeclareNamed[TestEnum]("TestValue", 456)
	if val != 456 {
		t.Errorf("Expected 456, got %d", val)
	}

	// Test duplicate name (should panic)
	defer func() {
		if r := recover(); !reflect.DeepEqual(r, "duplicate enum name: TestValue") {
			t.Error("Expected panic for duplicate name")
		}
	}()
	DeclareNamed[TestEnum]("TestValue", 789)
}

func TestDeclareDuplicateValue(t *testing.T) {
	cleanValues[TestEnum]()

	val := Declare[TestEnum](123)
	if val != 123 {
		t.Errorf("Expected 123, got %d", val)
	}

	defer func() {
		if r := recover(); !reflect.DeepEqual(r, "duplicate enum value: 123") {
			t.Error("Expected panic for duplicate value")
		}
	}()
	Declare[TestEnum](123)
}

func TestIsValid(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("TestValue", 123)

	if !IsValid[TestEnum](123) {
		t.Error("Expected 123 to be valid")
	}
	if IsValid[TestEnum](999) {
		t.Error("Expected 999 to be invalid")
	}
}

func TestFromName(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("TestValue", 123)

	val, ok := FromName[TestEnum]("TestValue")
	if !ok || val != 123 {
		t.Errorf("Expected 123, got %d, ok=%v", val, ok)
	}

	val, ok = FromName[TestEnum]("NonExistent")
	if ok || val != 0 {
		t.Errorf("Expected zero value and false, got %d, ok=%v", val, ok)
	}
}

func TestNames(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("TestValue", 123)
	DeclareNamed[TestEnum]("AnotherValue", 321)

	names := slices.Collect(Names[TestEnum]())
	expected := []string{"AnotherValue", "TestValue"}
	slices.Sort(names)
	if !slices.Equal(names, expected) {
		t.Errorf("Expected %v, got %v", expected, names)
	}
}

func TestValues(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("TestValue", 123)
	DeclareNamed[TestEnum]("AnotherValue", 456)

	values := slices.Collect(Values[TestEnum]())
	expected := []TestEnum{123, 456}
	slices.Sort(values)
	if !slices.Equal(values, expected) {
		t.Errorf("Expected %v, got %v", expected, values)
	}
}

func TestStringEnum(t *testing.T) {
	cleanValues[StringEnum]()
	DeclareNamed[StringEnum]("StringA", "a")
	DeclareNamed[StringEnum]("StringB", "b")

	if !IsValid[StringEnum]("a") {
		t.Error("Expected 'a' to be valid")
	}
	if IsValid[StringEnum]("c") {
		t.Error("Expected 'c' to be invalid")
	}

	val, ok := FromName[StringEnum]("StringB")
	if !ok || val != "b" {
		t.Errorf("Expected 'b', got %s, ok=%v", val, ok)
	}

	names := slices.Collect(Names[StringEnum]())
	expectedNames := []string{"StringA", "StringB"}
	slices.Sort(names)
	if !slices.Equal(names, expectedNames) {
		t.Errorf("Expected %v, got %v", expectedNames, names)
	}
}
