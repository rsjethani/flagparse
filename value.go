package flagparse

import (
	"fmt"
	"strconv"
)

// The Value interface specifies desired behavior that a type must have in order to be used with
// this package. Please see the implementation of Bool, Int etc. types in this pacakge as examples.
type Value interface {
	// Set takes a variable number of arguments and returns error if any of the arguments cannot be
	// parsed/converted correctly into the underlying type. For 'switch' types Set() will be called
	// with no arguments. Types that require only a single argument (Int, Float64 etc. for example)
	// would care only about the 0th argument and ignore the rest. Types that implement some kind of
	// list/slice/collection (IntList, Float64List for example) would normally want to parse all
	// given arguments.
	Set(...string) error

	// Get should return the value of underlying variable. The returned value's type should be the
	// same as underlying type
	Get() interface{}

	// String returns the current value of underlying variable as a string. This is useful for showing
	// default values in the help message.
	String() string
}

// NewValue takes address of a variable and returns a compatible Value type so that it can be used
// with this package. It returns error if there is no compatible type for the variable.
func NewValue(v interface{}) (Value, error) {
	switch addr := v.(type) {
	case Value: // the type itself implements Value interface hence simply return addr
		return addr, nil
	case *bool:
		return NewBool(addr), nil
	case *[]bool:
		return NewBoolList(addr), nil
	case *string:
		return NewString(addr), nil
	case *[]string:
		return NewStringList(addr), nil
	case *int:
		return NewInt(addr), nil
	case *[]int:
		return NewIntList(addr), nil
	case *float64:
		return NewFloat64(addr), nil
	case *[]float64:
		return NewFloat64List(addr), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", addr)
	}
}

// Bool type represents a bool value and also implements Value interface
type Bool bool

func NewBool(p *bool) *Bool {
	return (*Bool)(p)
}

func (b *Bool) Set(values ...string) error {
	if len(values) == 0 {
		// since Bool is a switch type calling Set() without args should set it to true
		values = append(values, "true")
	}
	v, err := strconv.ParseBool(values[0])
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", true), err)
	}
	*b = Bool(v)
	return nil
}

func (b *Bool) Get() interface{} { return bool(*b) }

func (b *Bool) String() string { return fmt.Sprint(*b) }

// Bool type represents a bool value and also implements Value interface
type BoolList []bool

func NewBoolList(p *[]bool) *BoolList {
	return (*BoolList)(p)
}

func (bl *BoolList) Set(values ...string) error {
	*bl = make([]bool, len(values))
	for i, val := range values {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return formatParseError(val, fmt.Sprintf("%T", true), err)
		}
		(*bl)[i] = v

	}
	return nil
}

func (bl *BoolList) Get() interface{} { return []bool(*bl) }

func (bl *BoolList) String() string { return fmt.Sprint(*bl) }

// Int type represents an int value
type Int int

func NewInt(p *int) *Int {
	return (*Int)(p)
}

// implement set like Bool does...do not change pointed value to zero if
// we get error while converting cmd arg string
func (i *Int) Set(values ...string) error {
	if len(values) == 0 {
		return nil
	}
	v, err := strconv.ParseInt(values[0], 0, strconv.IntSize)
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", int(1)), err)
	}
	*i = Int(v)
	return nil
}

func (i *Int) Get() interface{} { return int(*i) }

func (i *Int) String() string { return strconv.Itoa(int(*i)) }

// IntList type representing a list of integer values
type IntList []int

func NewIntList(p *[]int) *IntList {
	return (*IntList)(p)
}

func (il *IntList) Set(values ...string) error {
	*il = make([]int, len(values))
	for i, val := range values {
		n, err := strconv.ParseInt(val, 0, strconv.IntSize)
		if err != nil {
			return formatParseError(val, fmt.Sprintf("%T", int(1)), err)
		}
		(*il)[i] = int(n)
	}
	return nil
}

func (il *IntList) Get() interface{} { return []int(*il) }

func (il *IntList) String() string { return fmt.Sprint(*il) }

// String type represents a string value and implements Value interface
type String string

func NewString(p *string) *String {
	return (*String)(p)
}

func (s *String) Set(values ...string) error {
	if len(values) == 0 {
		return nil
	}
	*s = String(values[0])
	return nil
}

func (s *String) Get() interface{} { return string(*s) }

func (s *String) String() string { return fmt.Sprint(*s) }

// StringList type represents a list string value and implements Value interface
type StringList []string

func NewStringList(p *[]string) *StringList {
	return (*StringList)(p)
}

func (sl *StringList) Set(values ...string) error {
	*sl = make([]string, len(values))
	for i, val := range values {
		(*sl)[i] = val
	}
	return nil
}

func (sl *StringList) Get() interface{} { return []string(*sl) }

func (sl *StringList) String() string { return fmt.Sprint(*sl) }

// Float64 represents a float64 value and also implements Value interface
type Float64 float64

func NewFloat64(p *float64) *Float64 {
	return (*Float64)(p)
}

func (f *Float64) Set(values ...string) error {
	if len(values) == 0 {
		return nil
	}
	v, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", float64(1)), err)
	}
	*f = Float64(v)
	return nil
}

func (f *Float64) Get() interface{} { return float64(*f) }

func (f *Float64) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// Float64List type representing a list of float64 values and implements Value interface
type Float64List []float64

func NewFloat64List(p *[]float64) *Float64List {
	return (*Float64List)(p)
}

func (fl *Float64List) Set(values ...string) error {
	*fl = make([]float64, len(values))

	for i, val := range values {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return formatParseError(val, fmt.Sprintf("%T", float64(1)), err)
		}
		(*fl)[i] = f
	}
	return nil
}

func (fl *Float64List) Get() interface{} { return []float64(*fl) }

func (fl *Float64List) String() string { return fmt.Sprint(*fl) }
