package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hypergopher/renderfish/htmx"
	"github.com/hypergopher/renderfish/request"
)

func assertEqual(t *testing.T, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func assertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Error("want true, got false")
	}
}

func assertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Error("want false, got true")
	}
}

func assertBool(t *testing.T, want, got bool) {
	t.Helper()
	if want != got {
		t.Errorf("want %t, got %t", want, got)
	}
}

func TestRequestInfoMethods(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "forwarded.example.com")
	req.Header.Set("X-Forwarded-Port", "8080")
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Boosted", "true")
	req.Header.Set("User-Agent", "testUserAgent")
	req.Header.Set("Referer", "testReferer")

	req.RemoteAddr = "10.0.0.2:8080"

	tests := []struct {
		name   string
		method func(*http.Request) string
		want   string
	}{
		{name: "Scheme", method: request.Scheme, want: "https"},
		{name: "Host", method: request.Host, want: "forwarded.example.com"},
		{name: "Port", method: request.Port, want: "8080"},
		{name: "Method", method: request.Method, want: "GET"},
		{name: "URLPath", method: request.URLPath, want: ""},
		{name: "Referer", method: request.Referer, want: "testReferer"},
		{name: "RemoteAddr", method: request.RemoteAddr, want: "10.0.0.1"},
		{name: "UserAgent", method: request.UserAgent, want: "testUserAgent"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertEqual(t, tt.want, tt.method(req))
		})
	}

	scheme, host, port := request.SchemeHostPort(req)
	assertEqual(t, "https", scheme)
	assertEqual(t, "forwarded.example.com", host)
	assertEqual(t, "8080", port)

	assertEqual(t, "https://forwarded.example.com:8080", request.BaseURL(req))

	assertTrue(t, request.IsSecure(req))
}

func TestHXHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		method  func(*http.Request) bool
		want    bool
	}{
		{name: "HX-Request", headers: map[string]string{"HX-Request": "true"}, method: htmx.IsHtmxRequest, want: true},
		{name: "HX-Boosted", headers: map[string]string{"HX-Boosted": "true"}, method: htmx.IsBoostedRequest, want: true},
		{name: "HX-Request and HX-Boosted", headers: map[string]string{"HX-Request": "true", "HX-Boosted": "true"}, method: htmx.IsAnyHtmxRequest, want: true},
		{name: "HX-Request and HX-Boosted means IsHtmxRequest is false", headers: map[string]string{"HX-Request": "true", "HX-Boosted": "true"}, method: htmx.IsHtmxRequest, want: false},
		{name: "No HX-Request", headers: map[string]string{}, method: htmx.IsHtmxRequest, want: false},
		{name: "No HX-Boosted", headers: map[string]string{}, method: htmx.IsBoostedRequest, want: false},
		{name: "No HX-Request and HX-Boosted", headers: map[string]string{}, method: htmx.IsAnyHtmxRequest, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			assertBool(t, tt.want, tt.method(req))
		})
	}
}

func TestRequestInfoMethodsEmptyHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com:1234", nil)
	req.RemoteAddr = "10.0.0.2:8080"

	assertEqual(t, "http", request.Scheme(req))
	assertEqual(t, "example.com", request.Host(req))
	assertEqual(t, "1234", request.Port(req))
	assertEqual(t, "http://example.com:1234", request.BaseURL(req))
	assertFalse(t, request.IsSecure(req))
	assertEqual(t, "10.0.0.2:8080", request.RemoteAddr(req))
}

func TestInPath(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com:1234/foo/bar/baz", nil)
	tests := []struct {
		name   string
		path   string
		option string
		req    *http.Request
		want   bool
	}{
		{name: "InPath exact match", option: "exact", path: "/foo/bar/baz", want: true},
		{name: "InPath no exact match", option: "exact", path: "/foo/bar", want: false},
		{name: "InPath contains", option: "contains", path: "/bar/", want: true},
		{name: "InPath no contains", option: "contains", path: "/bizzle", want: false},
		{name: "InPath suffix", option: "suffix", path: "/baz", want: true},
		{name: "InPath no suffix", option: "suffix", path: "/foo", want: false},
		{name: "InPath prefix", option: "prefix", path: "/foo", want: true},
		{name: "InPath no prefix", option: "prefix", path: "/bar", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBool(t, tt.want, request.InPath(req, tt.path, tt.option))
		})
	}
}
