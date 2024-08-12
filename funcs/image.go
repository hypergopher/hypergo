package funcs

import "strings"

// Srcset returns a srcset string for an image with the given src and widths.
// It takes a src string and a variadic list of widths. The widths can be
// either a number or a number followed by a "w" or "x" suffix. If the width
// does not have a suffix, it will default to "w". Use "x" for density descriptors.
// Example:
//
//	Srcset("/foo/example.webp", "100", "200", "300")
//	// => "/foo/example-100w.webp, /foo/example-200w.webp, /foo/example-300w.webp"
//	Srcset("/foo/example.webp", "100w", "200w", "300w")
//	// => "/foo/example-100w.webp, /foo/example-200w.webp, /foo/example-300w.webp"
//	Srcset("/example.jpg", "1x", "2x")
//	// => "/example-1x.jpg, /example-2x.jpg"
func Srcset(src string, widths ...string) string {
	srcset := ""

	if strings.TrimSpace(src) == "" {
		return srcset
	}

	for i, w := range widths {
		if i > 0 {
			srcset += ", "
		}
		// Extract the extension from the src based on the last period
		// and append the width and the extension
		ext := src[strings.LastIndex(src, "."):]

		// If the width does not have a "w" suffix, and it does not have an "x" suffix, then append "w"
		if !strings.HasSuffix(w, "w") && !strings.HasSuffix(w, "x") {
			w += "w"
		}

		// Append the src (minus the extension) with the width and extension
		srcset += src[:strings.LastIndex(src, ".")] + "-" + w + ext
	}
	return srcset
}
