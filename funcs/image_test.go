package funcs_test

import (
	"strings"
	"testing"

	"github.com/hypergopher/renderfish/funcs"
)

func TestSrcset(t *testing.T) {
	var tests = []struct {
		name   string
		src    string
		widths []string
		want   string
	}{
		{"widths (no w)", "/foo/example.webp", []string{"100", "200", "300"}, "/foo/example-100w.webp, /foo/example-200w.webp, /foo/example-300w.webp"},
		{"widths (with w)", "/foo/example.webp", []string{"100w", "200w", "300w"}, "/foo/example-100w.webp, /foo/example-200w.webp, /foo/example-300w.webp"},
		{"densities", "/example.jpg", []string{"1x", "2x"}, "/example-1x.jpg, /example-2x.jpg"},
		{"url", "https://cdn.example.com/image-54.png", []string{"1000", "1500"}, "https://cdn.example.com/image-54-1000w.png, https://cdn.example.com/image-54-1500w.png"},
		{"no sizes", "example.tiff", []string{}, ""},
		{"no image", "", []string{"100", "200"}, ""},
	}

	for _, tt := range tests {
		testname := tt.src + strings.Join(tt.widths, ", ")
		t.Run(testname, func(t *testing.T) {
			ans := funcs.Srcset(tt.src, tt.widths...)
			if ans != tt.want {
				t.Errorf("%s: got %s, want %s", tt.name, ans, tt.want)
			}
		})
	}
}
