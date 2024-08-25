package hyperview

import (
	"net/http"

	"github.com/hypergopher/hyperview/response"
)

// JSONAdapter is an adapter for rendering JSON responses.
type JSONAdapter struct{}

// NewJSONViewAdapter creates a new JSON view adapter.
func NewJSONViewAdapter() *JSONAdapter {
	return &JSONAdapter{}
}

func (v *JSONAdapter) Init() error {
	return nil
}

func (v *JSONAdapter) Render(w http.ResponseWriter, r *http.Request, resp *response.Response) {
	if resp.StatusCode() == 0 {
		resp.Status(http.StatusOK)
	}

	if resp.StatusCode() > 299 {
		err := JSONFailure(w, resp.ViewData(r).Data(), "Failure", resp.StatusCode(), resp.HTTPHeader())
		if err != nil {
			v.RenderSystemError(w, r, err, resp)
		}
		return
	}

	err := JSONSuccessWithStatus(w, resp.StatusCode(), resp.ViewData(r).Data(), resp.HTTPHeader())
	if err != nil {
		v.RenderSystemError(w, r, err, resp)
	}
}

func (v *JSONAdapter) RenderForbidden(w http.ResponseWriter, _ *http.Request, _ *response.Response) {
	err := JSONFailure(w, nil, "Forbidden", http.StatusForbidden, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *JSONAdapter) RenderMaintenance(w http.ResponseWriter, _ *http.Request, _ *response.Response) {
	err := JSONFailure(w, nil, "Maintenance", http.StatusServiceUnavailable, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *JSONAdapter) RenderMethodNotAllowed(w http.ResponseWriter, _ *http.Request, _ *response.Response) {
	err := JSONFailure(w, nil, "Method not allowed", http.StatusMethodNotAllowed, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *JSONAdapter) RenderNotFound(w http.ResponseWriter, _ *http.Request, _ *response.Response) {
	err := JSONFailure(w, nil, "Not found", http.StatusNotFound, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *JSONAdapter) RenderSystemError(w http.ResponseWriter, _ *http.Request, err error, _ *response.Response) {
	e := JSONError(w, err.Error(), http.StatusInternalServerError, nil)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func (v *JSONAdapter) RenderUnauthorized(w http.ResponseWriter, _ *http.Request, _ *response.Response) {
	err := JSONFailure(w, nil, "Unauthorized", http.StatusUnauthorized, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
