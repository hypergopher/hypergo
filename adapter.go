package hypergo

import (
	"net/http"
)

// Adapter is an interface for rendering templates with various template engines.
//
//goland:noinspection GoNameStartsWithPackageName
type Adapter interface {
	// Init initializes the view adapter.
	Init() error
	// Render renders the specified template with the provided data.
	Render(w http.ResponseWriter, r *http.Request, opts *Response)
	// RenderForbidden renders the forbidden page.
	RenderForbidden(w http.ResponseWriter, r *http.Request, opts *Response)
	// RenderMaintenance renders the maintenance page.
	RenderMaintenance(w http.ResponseWriter, r *http.Request, opts *Response)
	// RenderMethodNotAllowed renders the method not allowed page.
	RenderMethodNotAllowed(w http.ResponseWriter, r *http.Request, opts *Response)
	// RenderNotFound renders the not found page.
	RenderNotFound(w http.ResponseWriter, r *http.Request, opts *Response)
	// RenderSystemError renders the system error page.
	RenderSystemError(w http.ResponseWriter, r *http.Request, err error, opts *Response)
	// RenderUnauthorized renders the unauthorized page.
	RenderUnauthorized(w http.ResponseWriter, r *http.Request, opts *Response)
}
