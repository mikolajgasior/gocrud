package builder

import "testing"

type structFieldValueFromString struct {
	Int64Field uint64
}

func TestStructFieldValueFromString(t *testing.T) {
	testObj := &structFieldValueFromString{}

	ok, value := StructFieldValueFromString(testObj, "Int64Field", "0")
	if !ok {
		t.Fatal("Failed to parse string value")
	}

	if value.(uint64) != 0 {
		t.Fatal("Parsed value is invalid")
	}
}

func TestSetStructFieldValueFromString(t *testing.T) {
	testObj := &structFieldValueFromString{}

	ok, value := SetStructFieldValueFromString(testObj, "Int64Field", "8")
	if !ok {
		t.Fatal("Failed to parse string value")
	}

	if value.(uint64) != 8 {
		t.Fatal("Parsed value is invalid")
	}
}
