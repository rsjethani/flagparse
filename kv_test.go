package flagparse

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_splitKV(t *testing.T) {
	data := make(map[string][]string)
	data[""] = []string{}
	data[fmt.Sprintf("%[1]c%[1]c%[1]c", kvPairSep)] = []string{}
	data[fmt.Sprintf("a%[1]cb%[1]cc%[1]cd", kvPairSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca%[1]cb%[1]cc%[1]cd%[1]c", kvPairSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca\\%[1]cb%[1]cc%[1]cd\\%[1]c", kvPairSep)] = []string{fmt.Sprintf("a%cb", kvPairSep), "c", fmt.Sprintf("d%c", kvPairSep)}
	for input, expected := range data {
		if got := splitKVs(input, kvPairSep); !reflect.DeepEqual(expected, got) {
			t.Errorf("Testing: splitKV(%q, %c); Expected: %q; Got: %q", input, kvPairSep, expected, got)
		}
	}
}

func Test_parseKVs_InvalidKeyValues(t *testing.T) {
	invalidKVs := []string{
		"hello",
		"usage=",
		"hello=hi",
		"name=flag_name",
		"nargs=1x",
	}

	for _, kv := range invalidKVs {
		if _, err := parseKVs(kv); err == nil {
			t.Errorf("Testing: parseTags(%q); Expected: error; Got: no error", kv)
		}
	}
}

func Test_parseKVs_ValidKeyValues(t *testing.T) {
	data := []struct {
		kvs      string
		expected map[string]string
	}{
		{
			"",
			map[string]string{},
		},
		{
			"positional,name=flag-name,usage=a usage message,nargs=10",
			map[string]string{
				nargsKey: "10",
				posKey:   "yes",
				nameKey:  "flag-name",
				usageKey: "a usage message",
			},
		},
		{
			"name=ArgName10,usage=abc,nargs=-10",
			map[string]string{
				nargsKey: "-10",
				nameKey:  "ArgName10",
				usageKey: "abc",
			},
		},
	}

	for _, input := range data {
		got, err := parseKVs(input.kvs)
		if err != nil {
			t.Errorf("Testing: parseTags(%q); Expected: no error; Got: %s", input.kvs, err)
		}

		if !reflect.DeepEqual(input.expected, got) {
			t.Errorf("Testing: parseTags(%q); Expected: %+v; Got: %+v", input.kvs, input.expected, got)
		}
	}
}

func Test_newFlagFromKVs_InvalidInput(t *testing.T) {
	testValue := newIntValue(new(int))
	data := []map[string]string{
		{posKey: "yes", nargsKey: "0"},
		{nargsKey: "123abc"},
		{nargsKey: "99999999999999999999999999999999999999"},
	}

	for _, input := range data {
		if _, err := newFlagFromKVs(testValue, input); err == nil {
			t.Errorf("Testing: newArgFromTags(%q); expected: error; got: no error", input)
		}
	}
}

func Test_newFlagFromKVs_ValidInput(t *testing.T) {
	x := 100
	testValue := newIntValue(&x)
	usageMsg := "usage message"
	data := []struct {
		keyValues map[string]string
		expected  Flag
	}{
		{
			keyValues: map[string]string{posKey: "yes", usageKey: usageMsg, nargsKey: "10"},
			expected:  Flag{value: testValue, defVal: testValue.String(), positional: true, nArgs: 10, usage: usageMsg},
		},
		{
			keyValues: map[string]string{usageKey: usageMsg, nargsKey: "10"},
			expected:  Flag{value: testValue, defVal: testValue.String(), positional: false, nArgs: 10, usage: usageMsg},
		},
	}

	for _, input := range data {
		fl, err := newFlagFromKVs(testValue, input.keyValues)
		if err != nil {
			t.Errorf("Testing: newFlagFromTags(%p,%q); Expected: no error; Got: error", testValue, input.keyValues)
		}
		if !reflect.DeepEqual(fl, &input.expected) {
			t.Errorf("Testing: newFlagFromTags(%p,%q); Expected: %+v; Got: %+v", testValue, input.keyValues, input.expected, fl)
		}
	}
}
