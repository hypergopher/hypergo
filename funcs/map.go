package funcs

import (
	"fmt"
	"strings"
)

// ClassMap takes a pair of classes and boolean expressions, and returns a single string with the active classes.
func ClassMap(classMap ...any) (string, error) {
	if len(classMap)%2 != 0 {
		return "", fmt.Errorf("ClassMap expects an even number of arguments")
	}

	var classes []string
	for i := 0; i < len(classMap); i += 2 {
		// Assert that the first in the pair is a string (class name)
		// and the second is a boolean (condition)
		className, ok := classMap[i].(string)
		if !ok {
			return "", fmt.Errorf("expected first argument to be a string")
		}

		condition, ok := classMap[i+1].(bool)
		if !ok {
			return "", fmt.Errorf("expected second argument to be a boolean")
		}

		if condition {
			classes = append(classes, className)
		}
	}

	return strings.Join(classes, " "), nil
}
