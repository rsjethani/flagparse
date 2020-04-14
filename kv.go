package flagparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	kvSep     rune   = '='
	kvPairSep        = ','
	nameKey   string = "name"
	usageKey         = "usage"
	nargsKey         = "nargs"
	posKey           = "positional"
)

var validKVs = map[string]*regexp.Regexp{
	usageKey: regexp.MustCompile(fmt.Sprintf(`^%s%c(.+)$`, usageKey, kvSep)),
	nameKey:  regexp.MustCompile(fmt.Sprintf(`^%s%c([[:alnum:]-]+)$`, nameKey, kvSep)),
	nargsKey: regexp.MustCompile(fmt.Sprintf(`^%s%c(-?[[:digit:]]+)$`, nargsKey, kvSep)),
}

func splitKVs(src string, sep rune) []string {
	backSlash := '\\'
	parts := make([]string, 0)
	b := &strings.Builder{}
	var prevRune rune
	for _, curRune := range src {
		switch {
		// rune is a backSlash simply skip it
		case curRune == backSlash:
		// rune is a sep but it is not escaped by backSlash
		case curRune == sep && prevRune != backSlash:
			if b.Len() != 0 {
				parts = append(parts, b.String())
				b.Reset()
			}
		// rune is either not a sep/backslash or if it is a sep then it is escaped by backskash
		default:
			b.WriteRune(curRune)
		}
		prevRune = curRune
	}
	// append any remaining runes between last sep and end of src
	if b.Len() != 0 {
		parts = append(parts, b.String())
	}
	return parts
}

func parseKVs(structTag string) (map[string]string, error) {
	kvs := make(map[string]string)
	for _, key := range splitKVs(structTag, kvPairSep) {
		if key == posKey {
			kvs[posKey] = "yes"
			continue
		}
		invalid := true
		for name, regex := range validKVs {
			res := regex.FindStringSubmatch(key)
			if len(res) == 2 {
				kvs[name] = res[1]
				invalid = false
			}
		}
		if invalid {
			return nil, fmt.Errorf("unknown key and/or invalid value: %s", key)
		}
	}
	return kvs, nil
}

func newFlagFromKVs(value Value, kvs map[string]string) (*Flag, error) {
	var fl *Flag
	if kvs[posKey] == "yes" {
		fl = NewFlag(value, true, kvs[usageKey])
	} else {
		fl = NewFlag(value, false, kvs[usageKey])
	}
	if kvs[nargsKey] != "" {
		nargs, err := strconv.ParseInt(kvs[nargsKey], 0, strconv.IntSize)
		if err != nil {
			return nil, formatParseError(kvs[nargsKey], fmt.Sprintf("%T", int(1)), err)
		}
		err = fl.SetNArgs(int(nargs))
		if err != nil {
			return nil, err
		}
	}
	return fl, nil
}
