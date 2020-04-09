package flagparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	kvPairSep rune   = ','
	kvSep     rune   = '='
	posKey    string = "positional"
)

var validKVs = map[string]*regexp.Regexp{
	"name":  regexp.MustCompile(fmt.Sprintf(`^name%c([[:alnum:]-]+)$`, kvSep)),
	"help":  regexp.MustCompile(fmt.Sprintf(`^help%c(.+)$`, kvSep)),
	"nargs": regexp.MustCompile(fmt.Sprintf(`^nargs%c(-?[[:digit:]]+)$`, kvSep)),
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
	var newFlag *Flag
	if kvs[posKey] == "yes" {
		newFlag = NewPosFlag(value, kvs["help"])
	} else {
		newFlag = NewOptFlag(value, kvs["help"])
	}
	if kvs["nargs"] != "" {
		nargs, err := strconv.ParseInt(kvs["nargs"], 0, strconv.IntSize)
		if err != nil {
			return nil, formatParseError(kvs["nargs"], fmt.Sprintf("%T", int(1)), err)
		}
		err = newFlag.SetNArgs(int(nargs))
		if err != nil {
			return nil, err
		}
	}
	return newFlag, nil
}
