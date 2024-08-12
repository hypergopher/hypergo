package funcs

import (
	"bytes"
	"strings"
	"unicode"
)

func Pluralize(count any, singular string, plural string) (string, error) {
	n, err := toInt64(count)
	if err != nil {
		return "", err
	}

	if n == 1 {
		return singular, nil
	}

	return plural, nil
}

// Humanize a string from snake case, slug, or camel case to a human readable string.
// Example: "hello_world" -> "Hello world"
func Humanize(s string) string {
	var buf bytes.Buffer

	for i, r := range s {
		switch {
		case i == 0:
			buf.WriteRune(unicode.ToUpper(r))
		case r == '_', r == '-':
			buf.WriteRune(' ')
		case unicode.IsUpper(r):
			buf.WriteRune(' ')
			buf.WriteRune(unicode.ToLower(r))
		default:
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

// Slugify converts a string to a slug naively. For a more robust solution, consider using a library like github.com/gosimple/slug
func Slugify(s string) string {
	var buf bytes.Buffer

	for _, r := range s {
		switch {
		case r > unicode.MaxASCII:
			continue
		case unicode.IsLetter(r):
			buf.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r), r == '_', r == '-':
			buf.WriteRune(r)
		case unicode.IsSpace(r):
			buf.WriteRune('-')
		}
	}

	return buf.String()
}

func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n] + "..."
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func NotBlank(s string) bool {
	return !IsBlank(s)
}
