package response

import "net/http"

// Renderer is the interface for a HyperGo response renderer
type Renderer interface {
	// Render renders the response to the given http.ResponseWriter
	Render(w http.ResponseWriter, r *http.Request, resp *Response)
}
