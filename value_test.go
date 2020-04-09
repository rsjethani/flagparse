package flagparse

import (
	"fmt"
	"math"
	"testing"
)

const (
	// uintSize      = 32 << (^uint(0) >> 32 & 1)
	minUint uint = 0
	maxUint uint = ^minUint
	maxInt  int  = int(maxUint >> 1)
	minInt  int  = -maxInt - 1
)

type customValue struct{}

func (c *customValue) Set(values ...string) error {
	return nil
}
func (c *customValue) String() string {
	return ""
}

func (c *customValue) Get() interface{} {
	return nil
}

func TestNewValue_SupportedType(t *testing.T) {
	// Test value creation for types implementing Value interface
	supported := []interface{}{
		new(customValue),
		new(int),
		new([]int),
		new(bool),
		new([]bool),
		new(string),
		new([]string),
		new(float64),
		new([]float64),
	}
	for _, val := range supported {
		_, err := NewValue(val)
		if err != nil {
			t.Errorf("Expected: NewValue(%T) should succeed, Got: %s", val, err)
		}
	}
}

func TestNewValue_UnsupportedType(t *testing.T) {
	type unsupported struct{}
	var testVar unsupported
	_, err := NewValue(&testVar)
	if err == nil {
		t.Errorf("Expected: NewValue(%T) should reult in error, Got: no error", testVar)
	}
}

func TestStringType(t *testing.T) {
	var testVar string
	arg := NewString(&testVar)

	data := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"", ""},
	}

	// Test valid values
	for _, val := range data {
		if err := arg.Set(val.input); err != nil {
			t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, val.input)
		}
		if val.expected != testVar {
			t.Errorf("Expected: %v, Got: %v", val.expected, testVar)
		}
		if val.input != arg.String() {
			t.Errorf("Expected: %v, Got: %v", val.input, arg.String())
		}
	}
}

func TestBoolType(t *testing.T) {
	var testVar bool
	arg := NewBool(&testVar)

	// Test Set() with no arguments
	arg.Set()
	if testVar != true {
		t.Errorf("Expected: true, Got: %v", testVar)
	}

	data := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	// Test valid values
	for _, val := range data {
		if err := arg.Set(val.input); err != nil {
			t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, val.input)
		}
		if val.expected != testVar {
			t.Errorf("Expected: %v, Got: %v", val.expected, testVar)
		}
		if val.input != arg.String() {
			t.Errorf("Expected: %v, Got: %v", val.input, arg.String())
		}
	}
}

func TestStringListType(t *testing.T) {
	var testVar []string
	arg := NewStringList(&testVar)
	data := struct {
		input    []string
		expected []string
	}{
		input:    []string{"hello", ""},
		expected: []string{"hello", ""},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
	}
}

func TestBoolListType(t *testing.T) {
	var testVar []bool
	arg := NewBoolList(&testVar)
	data := struct {
		input    []string
		expected []bool
	}{
		input:    []string{"true", "false"},
		expected: []bool{true, false},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
	}

	// Test invalid values
	input := []string{"tRUe", "hello", "1.1"}
	if err := arg.Set(input...); err == nil {
		t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
	}

}

func TestIntType(t *testing.T) {
	var testVar int
	intVal := NewInt(&testVar)

	// Test valid values
	validValues := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"10", 10},
		{"-10", -10},
		{fmt.Sprint(maxInt), maxInt},
		{fmt.Sprint(minInt), minInt},
	}
	for _, val := range validValues {
		if err := intVal.Set(val.input); err != nil {
			t.Errorf("Expected: no error for Int.Set(%q); Got: error %q", val.input, err)
		}
		if testVar != val.expected {
			t.Errorf("Expected: Int's underlying variable should have the value %d; Got: %d", val.expected, testVar)
		}
		if intVal.Get() != val.expected {
			t.Errorf("Expected: Int.Get() should return the value %d; Got: %d", val.expected, intVal.Get())
		}
		if intVal.String() != val.input {
			t.Errorf("Expected: Int.String() should return the string %q, Got: %q", val.input, intVal.String())
		}
	}

	// Test invalid values
	invalidValues := []string{"hello", "1.1", "true", "666666666666666666666666"}
	for _, input := range invalidValues {
		if err := intVal.Set(input); err == nil {
			t.Errorf("Expected: Int.Set(%q) should result in error, Got: no error", input)
		}
	}
}

func TestIntListType(t *testing.T) {
	var testVar []int
	arg := NewIntList(&testVar)
	data := struct {
		input    []string
		expected []int
	}{
		input:    []string{"0", "10", "-10", fmt.Sprint(maxInt), fmt.Sprint(minInt)},
		expected: []int{0, 10, -10, maxInt, minInt},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
	}

	// Test invalid values
	input := []string{"hello", "100", "true", "666666666666666666666666"}
	if err := arg.Set(input...); err == nil {
		t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
	}
}

func TestFloat64Type(t *testing.T) {
	var testVar float64
	arg := NewFloat64(&testVar)

	data := []struct {
		input    string
		expected float64
	}{
		{"0", 0},
		{"100", 100.00},
		{"10.11", 10.11},
		{"-10.11", -10.11},
		{fmt.Sprint(math.MaxFloat64), math.MaxFloat64},
		{fmt.Sprint(math.SmallestNonzeroFloat64), math.SmallestNonzeroFloat64},
	}

	// Test valid values
	for _, val := range data {
		if err := arg.Set(val.input); err != nil {
			t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, val.input)
		}
		if val.expected != testVar {
			t.Errorf("Expected: %v, Got: %v", val.expected, testVar)
		}
		if val.input != arg.String() {
			t.Errorf("Expected: %v, Got: %v", val.input, arg.String())
		}
	}

	// Test invalid values
	for _, input := range []string{"hello", "1.1xx", "true", "100abcd"} {
		if err := arg.Set(input); err == nil {
			t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
		}
	}
}

func TestFloat64ListType(t *testing.T) {
	var testVar []float64
	arg := NewFloat64List(&testVar)
	data := struct {
		input    []string
		expected []float64
	}{
		input:    []string{"0", "100", "10.11", "-10.11", fmt.Sprint(math.MaxFloat64), fmt.Sprint(math.SmallestNonzeroFloat64)},
		expected: []float64{0, 100.00, 10.11, -10.11, math.MaxFloat64, math.SmallestNonzeroFloat64},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
	}

	// Test invalid values
	input := []string{"hello", "1.1", "true", "66666666666"}
	if err := arg.Set(input...); err == nil {
		t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
	}
}
