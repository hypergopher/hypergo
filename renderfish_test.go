package renderfish_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hypergopher/renderfish"
	"github.com/hypergopher/renderfish/response"
)

type mockViewAdapter struct {
	renderCalled bool
}

func (ma *mockViewAdapter) Init() error { return nil }
func (ma *mockViewAdapter) Render(w http.ResponseWriter, r *http.Request, resp *response.Response) {
	ma.renderCalled = true
}
func (ma *mockViewAdapter) RenderForbidden(w http.ResponseWriter, r *http.Request, resp *response.Response) {
}
func (ma *mockViewAdapter) RenderMaintenance(w http.ResponseWriter, r *http.Request, resp *response.Response) {
}
func (ma *mockViewAdapter) RenderMethodNotAllowed(w http.ResponseWriter, r *http.Request, resp *response.Response) {
}
func (ma *mockViewAdapter) RenderNotFound(w http.ResponseWriter, r *http.Request, resp *response.Response) {
}
func (ma *mockViewAdapter) RenderSystemError(w http.ResponseWriter, r *http.Request, err error, resp *response.Response) {
}
func (ma *mockViewAdapter) RenderUnauthorized(w http.ResponseWriter, r *http.Request, resp *response.Response) {
}

func TestViewService_RegisterAdapter(t *testing.T) {
	hgo, err := renderfish.NewRenderFish()
	if err != nil {
		t.Fatalf("error creating RenderFish: %v", err)
	}

	adapter1 := &mockViewAdapter{}
	adapter2 := &mockViewAdapter{}
	adapter3 := &mockViewAdapter{}

	tests := []struct {
		name    string
		key     string
		adapter renderfish.Adapter
	}{
		{
			name:    "Test Case 1: Single Registration",
			key:     "A1",
			adapter: adapter1,
		},
		{
			name:    "Test Case 2: Update Existing Registration",
			key:     "A1",
			adapter: adapter2,
		},
		{
			name:    "Test Case 3: RegisterAdapter Another Adapter",
			key:     "B1",
			adapter: adapter3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// RegisterAdapter the adapter
			_ = hgo.RegisterAdapter(test.key, test.adapter)

			actual, ok := hgo.Adapter(test.key)

			if !ok {
				t.Errorf("Adapter '%s' was not correctly registered", test.key)
			}

			if actual != test.adapter {
				t.Error("Registered and retrieved adapters do not match")
			}
		})
	}
}

func TestViewService_Redirect(t *testing.T) {
	hgo, err := renderfish.NewRenderFish()
	if err != nil {
		t.Fatalf("error creating RenderFish: %v", err)
	}

	var tests = []struct {
		name       string
		request    *http.Request
		url        string
		code       int
		wantStatus int
		wantBody   string
		wantHeader http.Header
	}{
		{
			name:       "HX request",
			request:    httptest.NewRequest("GET", "/", nil),
			url:        "https://example.com",
			wantStatus: http.StatusSeeOther,
			wantBody:   "redirecting...",
			wantHeader: http.Header{
				"HX-Redirect": []string{"https://example.com"},
			},
		},
		{
			name:       "XMLHttpRequest",
			request:    httptest.NewRequest("GET", "/", nil),
			url:        "https://example.com",
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"redirecting...","status":"redirect","url":"https://example.com"}`,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		{
			name:       "default request",
			request:    httptest.NewRequest("GET", "/", nil),
			url:        "https://example.com",
			wantStatus: http.StatusFound,
			wantHeader: http.Header{
				"Location": []string{"https://example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "HX request" {
				tt.request.Header.Set("HX-Request", "true")
			} else if tt.name == "XMLHttpRequest" {
				tt.request.Header.Set("X-Requested-With", "XMLHttpRequest")
			}

			rr := httptest.NewRecorder()
			hgo.Redirect(rr, tt.request, tt.url)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if rr.Body.String() != tt.wantBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), tt.wantBody)
				}
			}

			for key, value := range tt.wantHeader {
				if !strings.EqualFold(rr.Header().Get(key), value[0]) {
					t.Errorf("handler returned wrong header %s: got %v want %v",
						key, rr.Header().Get(key), value[0])
				}
			}
		})
	}
}

func TestViewService_Render(t *testing.T) {
	hgo, err := renderfish.NewRenderFish()
	if err != nil {
		t.Fatalf("error creating RenderFish: %v", err)
	}

	mockedAdapter := &mockViewAdapter{}
	mockedJSONAdapter := &mockViewAdapter{}
	_ = hgo.RegisterAdapter("html", mockedAdapter)
	_ = hgo.RegisterAdapter("json", mockedJSONAdapter)

	r, err := http.NewRequest("GET", "/test.html", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	// create a test ResponseWriter via httptest
	w := httptest.NewRecorder()

	tests := []struct {
		name        string
		resp        *response.Response
		wantErr     bool
		wantAdapter renderfish.Adapter
	}{
		{
			name:        "HTMLRenderWithExtension",
			resp:        response.NewResponse().Path("sample.html"),
			wantErr:     false,
			wantAdapter: mockedAdapter,
		},
		{
			name:        "HTMLRenderWithoutExtension",
			resp:        response.NewResponse().Path("sample"),
			wantErr:     false,
			wantAdapter: mockedAdapter,
		},
		{
			name:        "JSONRender",
			resp:        response.NewResponse().Path("sample.json"),
			wantErr:     false,
			wantAdapter: mockedJSONAdapter,
		},
		{
			name:        "JSONResponseWithHeader",
			resp:        response.NewResponse().Path("sample").Header("Content-Type", "application/json"),
			wantErr:     false,
			wantAdapter: mockedJSONAdapter,
		},
		{
			name:        "NonHTMLRender",
			resp:        response.NewResponse().Path("sample.pdf"),
			wantErr:     true,
			wantAdapter: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedAdapter.renderCalled = false
			mockedJSONAdapter.renderCalled = false
			hgo.Render(w, r, tt.resp)

			if tt.wantAdapter != nil {
				if tt.wantAdapter == mockedAdapter {
					if !mockedAdapter.renderCalled {
						t.Error("Render() did not call the HTML adapter's Render method")
					}
				} else if tt.wantAdapter == mockedJSONAdapter {
					if !mockedJSONAdapter.renderCalled {
						t.Error("Render() did not call the JSON adapter's Render method")
					}
				}
			}

			if tt.wantErr {
				resp := w.Result()
				body := w.Body.String()

				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("Render() expected %v, got %v", http.StatusInternalServerError, resp.StatusCode)
				}

				if !strings.Contains(body, "Adapter not found") {
					t.Errorf("Render() expected 'Adapter not found' in response body")
				}
			}
		})
	}
}
