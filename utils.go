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
