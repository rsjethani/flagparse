package flagparse

import (
	"fmt"
	"strconv"
)

func formatParseError(val string, typeName string, err error) error {
	var reason string
	if ne, ok := err.(*strconv.NumError); ok {
		reason = ne.Err.Error()
	} else {
		reason = err.Error()
	}
	return fmt.Errorf("cannot parse '%s' as type '%s': %s", val, typeName, reason)
}

// Value interface specifies the desired behavior that a type must have in order to be used for
// creating flags. Please see the implementation of boolValue, intValue etc. types in this package
// as examples.
type Value interface {
	// Set is called once for each flag with the number of arguments equal to the flag's nargs
	// value. It returns error if any of the arguments cannot be parsed into the underlying type.
	// Since the number of arguments is variable each implementation can handle these arguments
	// differently. As a general guideline types representing some kind of list/collection for e.g.
	// []int would iterate over all the given arguments whereas types representing single value for
	// e.g int, uint etc. would only care about the first argument. Types represnting a switch
	// (nargs=0), would not be passed any arguments. Hence they should take appropriate action
	// without any arguments being passed to them. For e.g. the boolValue type in this package
	// simply sets the underlying bool variable to true when Set() is called on it.
	Set(...string) error

	// Get simply returns the value of underlying variable.
	Get() interface{}

	// String returns the string representation of the underlying value. This is useful for
	// showing default values in the usage message.
	String() string
}

// newValue takes address of a variable and returns a compatible Value type so that it can be used
// with this package. It returns error if there is no compatible type for the variable.
func newValue(v interface{}) (Value, error) {
	switch addr := v.(type) {
	case Value: // the type itself implements Value interface hence simply return addr
		return addr, nil
	case *bool:
		return newBoolValue(addr), nil
	case *[]bool:
		return newBoolListValue(addr), nil
	case *string:
		return newStringValue(addr), nil
	case *[]string:
		return newStringListValue(addr), nil
	case *int:
		return newIntValue(addr), nil
	case *[]int:
		return newIntListValue(addr), nil
	case *float64:
		return newFloat64Value(addr), nil
	case *[]float64:
		return newFloat64ListValue(addr), nil
	default:
		return nil, fmt.Errorf("type '%T' does not implement the Value interface", addr)
	}
}

// boolValue wraps the built-in bool type and implements the Value interface
type boolValue bool

func newBoolValue(p *bool) *boolValue {
	return (*boolValue)(p)
}

func (b *boolValue) Set(values ...string) error {
	// allow boolValue to be used as a switch flag: set underlying variable to true when Set()
	// called without arguments
	if len(values) == 0 {
		values = append(values, "true")
	}
	v, err := strconv.ParseBool(values[0])
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", true), err)
	}
	*b = boolValue(v)
	return nil
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return fmt.Sprint(*b) }

// boolListValue wraps the built-in []bool type and implements the Value interface
type boolListValue []bool

func newBoolListValue(p *[]bool) *boolListValue {
	return (*boolListValue)(p)
}

func (bl *boolListValue) Set(values ...string) error {
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

func (bl *boolListValue) Get() interface{} { return []bool(*bl) }

func (bl *boolListValue) String() string { return fmt.Sprint(*bl) }

// stringValue wraps the built-in string type and implements the Value interface
type stringValue string

func newStringValue(p *string) *stringValue {
	return (*stringValue)(p)
}

func (s *stringValue) Set(values ...string) error {
	if len(values) == 0 {
		return nil
	}
	*s = stringValue(values[0])
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return fmt.Sprint(*s) }

// stringListValue wraps the built-in []string type and implements the Value interface
type stringListValue []string

func newStringListValue(p *[]string) *stringListValue {
	return (*stringListValue)(p)
}

func (sl *stringListValue) Set(values ...string) error {
	*sl = make([]string, len(values))
	for i, val := range values {
		(*sl)[i] = val
	}
	return nil
}

func (sl *stringListValue) Get() interface{} { return []string(*sl) }

func (sl *stringListValue) String() string { return fmt.Sprint(*sl) }

// intValue wraps the built-in int type and implements the Value interface
type intValue int

func newIntValue(p *int) *intValue {
	return (*intValue)(p)
}

// implement set like Bool does...do not change pointed value to zero if
// we get error while converting cmd arg string
func (i *intValue) Set(values ...string) error {
	if len(values) == 0 {
		return nil
	}
	v, err := strconv.ParseInt(values[0], 0, strconv.IntSize)
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", int(1)), err)
	}
	*i = intValue(v)
	return nil
}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

// intListValue wraps the built-in []int type and implements the Value interface
type intListValue []int

func newIntListValue(p *[]int) *intListValue {
	return (*intListValue)(p)
}

func (il *intListValue) Set(values ...string) error {
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

func (il *intListValue) Get() interface{} { return []int(*il) }

func (il *intListValue) String() string { return fmt.Sprint(*il) }

// float64Value wraps the built-in float64 type and implements the Value interface
type float64Value float64

func newFloat64Value(p *float64) *float64Value {
	return (*float64Value)(p)
}

func (f *float64Value) Set(values ...string) error {
	if len(values) == 0 {
		return nil
	}
	v, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", float64(1)), err)
	}
	*f = float64Value(v)
	return nil
}

func (f *float64Value) Get() interface{} { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// float64ListValue wraps the built-in []float64 type and implements the Value interface
type float64ListValue []float64

func newFloat64ListValue(p *[]float64) *float64ListValue {
	return (*float64ListValue)(p)
}

func (fl *float64ListValue) Set(values ...string) error {
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

func (fl *float64ListValue) Get() interface{} { return []float64(*fl) }

func (fl *float64ListValue) String() string { return fmt.Sprint(*fl) }
