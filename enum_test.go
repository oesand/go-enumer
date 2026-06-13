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

func TestEntries(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("First", 10)
	DeclareNamed[TestEnum]("Second", 20)
	DeclareNamed[TestEnum]("Third", 30)

	// Create a map for easier verification
	entryMap := make(map[string]TestEnum)
	count := 0
	for name, value := range Entries[TestEnum]() {
		entryMap[name] = value
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 entries, got %d", count)
	}

	expectedEntries := map[string]TestEnum{
		"First":  10,
		"Second": 20,
		"Third":  30,
	}

	if !reflect.DeepEqual(entryMap, expectedEntries) {
		t.Errorf("expected %v, got %v", expectedEntries, entryMap)
	}
}

func TestEntriesStringEnum(t *testing.T) {
	cleanValues[StringEnum]()
	DeclareNamed[StringEnum]("Red", "red")
	DeclareNamed[StringEnum]("Green", "green")

	entryMap := make(map[string]StringEnum)
	count := 0
	for name, value := range Entries[StringEnum]() {
		entryMap[name] = value
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 entries, got %d", count)
	}

	expectedEntries := map[string]StringEnum{
		"Red":   "red",
		"Green": "green",
	}

	if !reflect.DeepEqual(entryMap, expectedEntries) {
		t.Errorf("expected %v, got %v", expectedEntries, entryMap)
	}
}

func TestEntriesEmpty(t *testing.T) {
	// Test with an enum that has never been declared
	type UnusedEnum int
	cleanValues[UnusedEnum]()
	DeclareNamed[UnusedEnum]("OnlyOne", 1)

	count := 0
	for name, value := range Entries[UnusedEnum]() {
		if name != "OnlyOne" || value != 1 {
			t.Errorf("Expected OnlyOne=1, got %s=%v", name, value)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 entry, got %d", count)
	}
}

func TestEntriesR(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("Alpha", 100)
	DeclareNamed[TestEnum]("Beta", 200)
	DeclareNamed[TestEnum]("Gamma", 300)

	enumType := reflect.TypeOf(TestEnum(0))
	entryMap := make(map[string]TestEnum)
	count := 0
	for name, value := range EntriesR(enumType) {
		entryMap[name] = value.(TestEnum)
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 entries, got %d", count)
	}

	expectedEntries := map[string]TestEnum{
		"Alpha": 100,
		"Beta":  200,
		"Gamma": 300,
	}

	if !reflect.DeepEqual(entryMap, expectedEntries) {
		t.Errorf("expected %v, got %v", expectedEntries, entryMap)
	}
}

func TestEntriesRStringEnum(t *testing.T) {
	cleanValues[StringEnum]()
	DeclareNamed[StringEnum]("ColorRed", "red")
	DeclareNamed[StringEnum]("ColorBlue", "blue")
	DeclareNamed[StringEnum]("ColorGreen", "green")

	enumType := reflect.TypeOf(StringEnum(""))
	entryMap := make(map[string]StringEnum)
	count := 0
	for name, value := range EntriesR(enumType) {
		entryMap[name] = value.(StringEnum)
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 entries, got %d", count)
	}

	expectedEntries := map[string]StringEnum{
		"ColorRed":   "red",
		"ColorBlue":  "blue",
		"ColorGreen": "green",
	}

	if !reflect.DeepEqual(entryMap, expectedEntries) {
		t.Errorf("expected %v, got %v", expectedEntries, entryMap)
	}
}

func TestEntriesRUnknownType(t *testing.T) {
	// Test with a type that was never declared as enum
	type NeverDeclared int
	unknownType := reflect.TypeOf(NeverDeclared(0))

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for unknown enum type")
		}
	}()

	EntriesR(unknownType)
}

func TestEntriesIteratorEarlyExit(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("One", 1)
	DeclareNamed[TestEnum]("Two", 2)
	DeclareNamed[TestEnum]("Three", 3)

	count := 0
	for range Entries[TestEnum]() {
		count++
		if count >= 2 {
			// Break early to ensure iterator handles it correctly
			break
		}
	}

	if count != 2 {
		t.Errorf("Expected iterator to stop after 2 iterations, got %d", count)
	}
}

func TestEntriesRIteratorEarlyExit(t *testing.T) {
	cleanValues[TestEnum]()
	DeclareNamed[TestEnum]("X", 10)
	DeclareNamed[TestEnum]("Y", 20)
	DeclareNamed[TestEnum]("Z", 30)

	enumType := reflect.TypeOf(TestEnum(0))
	count := 0
	for range EntriesR(enumType) {
		count++
		if count >= 2 {
			// Break early to ensure iterator handles it correctly
			break
		}
	}

	if count != 2 {
		t.Errorf("Expected iterator to stop after 2 iterations, got %d", count)
	}
}
