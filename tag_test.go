package flagparse

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_splitKV(t *testing.T) {
	data := make(map[string][]string)
	data[""] = []string{}
	data[fmt.Sprintf("%[1]c%[1]c%[1]c", tagSep)] = []string{}
	data[fmt.Sprintf("a%[1]cb%[1]cc%[1]cd", tagSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca%[1]cb%[1]cc%[1]cd%[1]c", tagSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca\\%[1]cb%[1]cc%[1]cd\\%[1]c", tagSep)] = []string{fmt.Sprintf("a%cb", tagSep), "c", fmt.Sprintf("d%c", tagSep)}
	for input, expected := range data {
		if got := splitKV(input, tagSep); !reflect.DeepEqual(expected, got) {
			t.Errorf("Testing: splitKV(%q, %c); Expected: %q; Got: %q", input, tagSep, expected, got)
		}
	}
}

func Test_parseTags_InvalidKeyValues(t *testing.T) {
	invalidKVs := []string{
		"hello",
		"help=",
		"hello=hi",
		"name=flag_name",
		"nargs=1x",
	}

	for _, kv := range invalidKVs {
		if _, err := parseTags(kv); err == nil {
			t.Errorf("Testing: parseTags(%q); Expected: error; Got: no error", kv)
		}
	}
}

func Test_parseTags_ValidKeyValues(t *testing.T) {
	data := []struct {
		validKVs string
		expected map[string]string
	}{
		{
			"",
			map[string]string{},
		},
		{
			"positional,name=flag-name,help=a help message,nargs=10",
			map[string]string{
				"nargs": "10",
				posKey:  "yes",
				"name":  "flag-name",
				"help":  "a help message",
			},
		},
		{
			"name=ArgName10,help=abc,nargs=-10",
			map[string]string{
				"nargs": "-10",
				"name":  "ArgName10",
				"help":  "abc",
			},
		},
	}

	for _, input := range data {
		got, err := parseTags(input.validKVs)
		if err != nil {
			t.Errorf("Testing: parseTags(%q); Expected: no error; Got: %s", input.validKVs, err)
		}

		if !reflect.DeepEqual(input.expected, got) {
			t.Errorf("Testing: parseTags(%q); Expected: %+v; Got: %+v", input.validKVs, input.expected, got)
		}
	}
}

func Test_newFlagFromTags_InvalidInput(t *testing.T) {
	testValue := NewInt(new(int))

	// Test that for invalid key/value syntax the error returned by parseTags is returned to the caller
	// Test that error is returned if nargs value has proper syntax but is outside valid range
	// Test that error is returned when positional key is given but nargs=0
	for _, input := range []string{"nargs=123abc", "nargs=9999999999999999999999999", "positional,nargs=0"} {
		if _, err := newFlagFromTags(testValue, "", input); err == nil {
			t.Errorf("Testing: newArgFromTags(%q); expected: error; got: no error", input)
		}
	}
}

func Test_newFlagFromTags_ValidInput(t *testing.T) {
	x := 100
	testValue := NewInt(&x)
	helpMsg := "help message"
	flagName := "flag-name"
	data := []struct {
		val       Value
		fName     string
		keyValues string
		expected  Flag
	}{
		{
			val:       testValue,
			fName:     "Field1",
			keyValues: "",
			expected:  Flag{name: "field1", value: testValue, defVal: "100", positional: false, nArgs: 1, help: ""},
		},
		{
			val:       testValue,
			fName:     "Field1",
			keyValues: fmt.Sprintf("name=%s,nargs=0,help=%s", flagName, helpMsg),
			expected:  Flag{name: flagName, value: testValue, defVal: "", positional: false, nArgs: 0, help: helpMsg},
		},
		{
			val:       testValue,
			fName:     "Field1",
			keyValues: fmt.Sprintf("name=%s,nargs=10,help=%s", flagName, helpMsg),
			expected:  Flag{name: flagName, value: testValue, defVal: "100", positional: false, nArgs: 10, help: helpMsg},
		},
		{
			val:       testValue,
			fName:     "Field1",
			keyValues: fmt.Sprintf("positional,name=%s,nargs=10,help=%s", flagName, helpMsg),
			expected:  Flag{name: flagName, value: testValue, defVal: "100", positional: true, nArgs: 10, help: helpMsg},
		},
	}

	for _, input := range data {
		fl, err := newFlagFromTags(input.val, input.fName, input.keyValues)
		if err != nil {
			t.Errorf("Testing: newFlagFromTags(%p,%q,%q); Expected: no error; Got: error", input.val, input.fName, input.keyValues)
		}
		if !reflect.DeepEqual(fl, &input.expected) {
			t.Errorf("Testing: newFlagFromTags(%p,%q,%q); Expected: %+v; Got: %+v", input.val, input.fName, input.keyValues, input.expected, fl)
		}
	}
}
