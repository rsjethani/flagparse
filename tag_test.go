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
			t.Errorf("testing: splitKV(%q, %c); expected: %q; Got: %q", input, tagSep, expected, got)
		}
	}
}

func Test_parseTags_InvalidKeyValues(t *testing.T) {
	invalidKVs := []string{
		"hello",
		"help=",
		"hello=hi",
		"name=flag_name",
		"type=OPT",
		"nargs=1x",
	}

	for _, kv := range invalidKVs {
		if _, err := parseTags(kv); err == nil {
			t.Errorf("testing: parseTags(%#v); expected: error; got: no error", kv)
		}
	}
}

func Test_parseTags_ValidKeyValues(t *testing.T) {
	data := []struct {
		validKVs string
		expected map[string]string
	}{
		{
			"name=flag-name,type=pos,help=a help message,nargs=10",
			map[string]string{
				"nargs": "10",
				"type":  "pos",
				"name":  "flag-name",
				"help":  "a help message",
			},
		},
		{
			"name=ArgName10,type=opt,help=a,nargs=-10",
			map[string]string{
				"nargs": "-10",
				"type":  "opt",
				"name":  "ArgName10",
				"help":  "a",
			},
		},
		{
			"type=switch,nargs=-10",
			map[string]string{
				"nargs": "-10",
				"type":  "switch",
			},
		},
	}

	for _, input := range data {
		got, err := parseTags(input.validKVs)
		if err != nil {
			t.Errorf("testing: parseTags(%#v); expected: no error; got: %s", input.validKVs, err)
		}

		if !reflect.DeepEqual(input.expected, got) {
			t.Errorf("testing: parseTags(%#v); expected: %+v; got: %+v", input.validKVs, input.expected, got)
		}
	}
}

func Test_newArgFromTags_InvalidInput(t *testing.T) {
	testValue := NewInt(new(int))
	testKVs := "nargs=123abc"
	if arg, err := newFlagFromTags(testValue, "", testKVs); arg != nil || err == nil {
		t.Errorf("testing: newArgFromTags(%#v); expected: non-nil error since key-value parsing should fail for invalid key/value; got: %#v, %#v ", testKVs, arg, err)
	}

	testKVs = "type=pos,nargs=0"
	if arg, err := newFlagFromTags(testValue, "", testKVs); arg != nil || err == nil {
		t.Errorf("testing: newArgFromTags(%#v); expected: non-nil error since nargs cannot be 0 for type=pos; got: %#v, %#v ", testKVs, arg, err)
	}

	testKVs = "type=switch,nargs=10"
	if arg, err := newFlagFromTags(testValue, "", testKVs); arg != nil || err == nil {
		t.Errorf("testing: newArgFromTags(%#v); expected: non-nil error since nargs can only be 0 for type=switch; got: %#v, %#v ", testKVs, arg, err)
	}

	testKVs = "nargs=9999999999999999999999999"
	if arg, err := newFlagFromTags(testValue, "", testKVs); arg != nil || err == nil {
		t.Errorf("testing: newArgFromTags(%#v); expected: non-nil error since nargs value overflows int size; got: %#v, %#v ", testKVs, arg, err)
	}
}

func Test_newArgFromTags_ValidInput(t *testing.T) {
	testValue := NewInt(new(int))

	testKVs := "help=help message"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(%s); expected: non error since empty key-values is a valid input; got: %#v, %#v", testKVs, arg, err)
	} else {
		if arg.name != "field1" {
			t.Errorf("testing: newArgFromTags(%s); expected: name==field1; got: %s", testKVs, arg.name)
		}
		if arg.positional {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.positional==false; got: %v", testKVs, arg.positional)
		}
		if arg.nArgs != 1 {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.nargs==1; got: %v", testKVs, arg.nArgs)
		}
		if arg.help != "help message" {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.Help==\"help message\"; got: %v", testKVs, arg.help)
		}
		if arg.value != testValue {
			t.Errorf("testing: newArgFromTags(testValue,\"Field1\",%s); expected: arg.Value==testValue; got: unequal", testKVs)
		}
	}

	testKVs = "name=hello,help=help message"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(%s); expected: non error since empty key-values is a valid input; got: %#v, %#v", testKVs, arg, err)
	} else {
		if arg.name != "hello" {
			t.Errorf("testing: newArgFromTags(%s); expected: name==hello; got: %s", testKVs, arg.name)
		}
	}

	testKVs = "type=switch,help=help message"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(nil,\"Field1\",%s); expected: non error; got: %#v, %#v", testKVs, arg, err)
	} else {
		if !arg.isSwitch() {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.positional==false; got: %v", testKVs, arg.positional)
		}
		if arg.help != "help message" {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.Help==\"help message\"; got: %v", testKVs, arg.help)
		}
		if arg.value != testValue {
			t.Errorf("testing: newArgFromTags(testValue,\"Field1\",%s); expected: arg.Value==testValue; got: unequal", testKVs)
		}
	}

	testKVs = "type=switch,help=help message,nargs=123"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg != nil || err == nil {
		t.Errorf("testing: newArgFromTags(nil,\"Field1\",%s); expected: error; got: %#v, %#v", testKVs, arg, err)
	}

	// Test explicit opt type
	testKVs = "type=opt,help=help message"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(nil,\"Field1\",%s); expected: non error; got: %#v, %#v", testKVs, arg, err)
	} else {
		if arg.positional {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.positional==false; got: %v", testKVs, arg.positional)
		}
		if arg.nArgs != 1 {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.nargs==1; got: %v", testKVs, arg.nArgs)
		}
		if arg.help != "help message" {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.Help==\"help message\"; got: %v", testKVs, arg.help)
		}
		if arg.value != testValue {
			t.Errorf("testing: newArgFromTags(testValue,\"Field1\",%s); expected: arg.Value==testValue; got: unequal", testKVs)
		}
	}

	// Test explicit opt type with nargs
	testKVs = "type=opt,help=help message,nargs=123"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(nil,\"Field1\",%s); expected: non error; got: %#v, %#v", testKVs, arg, err)
	} else {
		if arg.nArgs != 123 {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.nargs==123; got: %v", testKVs, arg.nArgs)
		}
	}

	// Test pos type
	testKVs = "type=pos,help=help message"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(nil,\"Field1\",%s); expected: non error; got: %#v, %#v", testKVs, arg, err)
	} else {
		if !arg.positional {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.positional==false; got: %v", testKVs, arg.positional)
		}
		if arg.nArgs != 1 {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.nargs==1; got: %v", testKVs, arg.nArgs)
		}
		if arg.help != "help message" {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.Help==\"help message\"; got: %v", testKVs, arg.help)
		}
		if arg.value != testValue {
			t.Errorf("testing: newArgFromTags(testValue,\"Field1\",%s); expected: arg.Value==testValue; got: unequal", testKVs)
		}
	}

	// Test pos type with nargs
	testKVs = "type=pos,help=help message,nargs=123"
	if arg, err := newFlagFromTags(testValue, "Field1", testKVs); arg == nil || err != nil {
		t.Errorf("testing: newArgFromTags(nil,\"Field1\",%s); expected: non error; got: %#v, %#v", testKVs, arg, err)
	} else {
		if arg.nArgs != 123 {
			t.Errorf("testing: newArgFromTags(%s); expected: arg.nargs==123; got: %v", testKVs, arg.nArgs)
		}
	}
}
