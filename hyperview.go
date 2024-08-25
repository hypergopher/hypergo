package hyperview

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/hypergopher/hyperview/constants"
	"github.com/hypergopher/hyperview/htmx"
	"github.com/hypergopher/hyperview/request"
	"github.com/hypergopher/hyperview/response"
)

// Option is a function that can be used to configure the HyperView struct.
type Option func(*HyperView) error

// HyperView provides a service to render views from different template adapters.
type HyperView struct {
	adapters      map[string]Adapter // map of view adapters
	baseLayout    string             // default layout to use if none is specified
	systemLayout  string             // layout to use for system pages
	filesystemMap map[string]fs.FS   // map of file systems to use for the view adapters
	funcMap       template.FuncMap   // map of html/template functions to pass to the view
	logger        *slog.Logger       // logger to use for the view service
	mu            sync.RWMutex       // protects the adapters map
}

// NewHyperView creates a new view service. It accepts a list of options to configure the view service.
//
// Available options:
//
//   - WithLayouts: sets the base and system layouts for the view service.
//   - WithFuncMap: sets an initial function map to use for the template engine.
//   - WithBaseTemplateFS: sets an initial template and assets filesystem to use for the template engine.
//   - WithLogger: sets an initial logger to use for the HyperView instance. If not set, a default logger is created when the HyperView instance is created.
//   - WithViewAdapter: sets a view adapter to use for the view service. If no view adapters are set, the default adapters are used. Default adapters
//     use html/template for html templates and json for json templates.
func NewHyperView(options ...Option) (*HyperView, error) {
	hgo := &HyperView{
		adapters:      make(map[string]Adapter),
		baseLayout:    "base",
		systemLayout:  "base",
		filesystemMap: nil,
		funcMap:       nil,
		logger:        nil,
	}

	// Apply options
	for _, opt := range options {
		if err := opt(hgo); err != nil {
			return nil, fmt.Errorf("error applying option: %w", err)
		}
	}

	// If no func map is set, create an empty func map
	if hgo.funcMap == nil {
		hgo.funcMap = make(template.FuncMap)
	}

	// If no logger is set, create a default logger
	if hgo.logger == nil {
		hgo.logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelDebug,
		}))
	}

	// If there is no adapter for html, use the default html adapter
	if _, ok := hgo.adapters["html"]; !ok {
		tempAdapter := NewTemplateViewAdapter(TemplateViewAdapterOptions{
			Extension:     ".tmpl",
			FileSystemMap: hgo.filesystemMap,
			Funcs:         hgo.funcMap,
			Logger:        hgo.logger,
		})

		if err := hgo.RegisterAdapter("html", tempAdapter); err != nil {
			return nil, fmt.Errorf("error registering default HTML adapter: %w", err)
		}
	}

	return hgo, nil
}

// WithLayouts sets the base and system layouts for the view service.
func WithLayouts(base, system string) Option {
	return func(hgo *HyperView) error {
		hgo.baseLayout = base
		hgo.systemLayout = system
		return nil
	}
}

// WithFuncMap sets an initial function map to use for the template engine.
// Additional functions can be added later via Plugin options.
func WithFuncMap(funcs template.FuncMap) Option {
	return func(hgo *HyperView) error {
		hgo.funcMap = funcs
		return nil
	}
}

// WithBaseTemplateFS sets an initial template and assets filesystem to use for the template engine.
func WithBaseTemplateFS(efs *embed.FS) Option {
	return func(hgo *HyperView) error {
		hgo.filesystemMap = map[string]fs.FS{constants.RootFSID: efs}
		return nil
	}
}

// WithLogger sets an initial logger to use for the HyperView instance. If not set, a default logger is created when the HyperView instance is created.
func WithLogger(logger *slog.Logger) Option {
	return func(hgo *HyperView) error {
		hgo.logger = logger
		return nil
	}
}

// WithViewAdapter sets a view adapter to use for the view service. If no view adapters are set, the default adapters are used.
func WithViewAdapter(name string, adapter Adapter) Option {
	return func(hgo *HyperView) error {
		return hgo.RegisterAdapter(name, adapter)
	}
}

// Logger provides access to the logger, so that plugins can use it.
func (s *HyperView) Logger() *slog.Logger {
	return s.logger
}

// MaybeRegisterDefaultAdapters registers the built-in adapters for
// using html/template for html templates and json for json templates, but only
// if they are not already registered.
func (s *HyperView) MaybeRegisterDefaultAdapters() error {
	// Check if the html adapter is already registered
	if _, ok := s.adapters["html"]; !ok {
		tempAdapter := NewTemplateViewAdapter(TemplateViewAdapterOptions{
			FileSystemMap: s.filesystemMap,
			Funcs:         s.funcMap,
			Logger:        s.logger,
		})

		if err := s.RegisterAdapter("html", tempAdapter); err != nil {
			return fmt.Errorf("error registering default HTML adapter: %w", err)
		}
	}

	// Check if the json adapter is already registered
	if _, ok := s.adapters["json"]; !ok {
		jsonAdapter := NewJSONViewAdapter()
		if err := s.RegisterAdapter("json", jsonAdapter); err != nil {
			return fmt.Errorf("error registering default JSON adapter: %w", err)
		}
	}
	return nil
}

// RegisterAdapter registers a new view adapter with the view service
func (s *HyperView) RegisterAdapter(name string, adapter Adapter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adapters[name] = adapter
	return s.adapters[name].Init()
}

// Reinit reinitialize the view service adapters. This is useful for reloading templates after they have changed.
func (s *HyperView) Reinit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, adapter := range s.adapters {
		//s.logger.Debug("Reinitializing view adapter", slog.String("adapter", fmt.Sprintf("%T", adapter)))
		if err := adapter.Init(); err != nil {
			return err
		}
	}
	return nil
}

// Adapter returns the view adapter with the specified name
func (s *HyperView) Adapter(name string) (Adapter, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	adapter, ok := s.adapters[name]
	return adapter, ok
}

// Render renders the specified opts with the provided adapter key
func (s *HyperView) Render(w http.ResponseWriter, r *http.Request, resp *response.Response) {
	// First, find an extension if there is one
	ext := ""
	if idx := strings.LastIndex(resp.TemplatePath(), "."); idx != -1 {
		ext = resp.TemplatePath()[idx:]
		resp.Path(resp.TemplatePath()[:idx])
	}

	// If the resp has a content-type header of application/json, use the json adapter
	if resp.HTTPHeader().Get("Content-Type") == "application/json" {
		s.RenderAs(w, r, "json", resp)
		return
	}

	// If the extension is empty or .html, use the html adapter
	if ext == "" || ext == ".html" {
		s.RenderAs(w, r, "html", resp)
		return
	}

	// Otherwise, use the specified extension
	s.RenderAs(w, r, ext[1:], resp)
}

// RenderAs renders the specified opts with the provided adapter key
func (s *HyperView) RenderAs(w http.ResponseWriter, r *http.Request, adapterKey string, resp *response.Response) {
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		// If there is no layout set, set the base layout
		if resp.TemplateLayout() == "" {
			resp.Layout(s.baseLayout)
		}
		adapter.Render(w, r, resp)
	}
}

// RenderNotFound renders a 404 not found page
func (s *HyperView) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	s.RenderNotFoundAs(w, r, "html")
}

// RenderNotFoundAs renders a 404 not found page as the specified adapter
func (s *HyperView) RenderNotFoundAs(w http.ResponseWriter, r *http.Request, adapterKey string) {
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		adapter.RenderNotFound(w, r, s.NewSystemResponse().StatusNotFound())
	}
}

// RenderSystemError renders a system error page
func (s *HyperView) RenderSystemError(w http.ResponseWriter, r *http.Request, err error) {
	s.RenderSystemErrorAs(w, r, "html", err)
}

// RenderSystemErrorAs renders a system error page as the specified adapter
func (s *HyperView) RenderSystemErrorAs(w http.ResponseWriter, r *http.Request, adapterKey string, err error) {
	s.logger.Error("Server error", slog.String("err", err.Error()))
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		adapter.RenderSystemError(w, r, err, s.NewSystemResponse().StatusError())
	}
}

// RenderMaintenance renders a maintenance page
func (s *HyperView) RenderMaintenance(w http.ResponseWriter, r *http.Request) {
	s.RenderMaintenanceAs(w, r, "html")
}

// RenderMaintenanceAs renders a maintenance page as the specified adapter
func (s *HyperView) RenderMaintenanceAs(w http.ResponseWriter, r *http.Request, adapterKey string) {
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		adapter.RenderMaintenance(w, r, s.NewSystemResponse().Status(http.StatusServiceUnavailable))
	}
}

// RenderForbidden renders a forbidden page
func (s *HyperView) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	s.RenderForbiddenAs(w, r, "html")
}

// RenderForbiddenAs renders a forbidden page as the specified adapter
func (s *HyperView) RenderForbiddenAs(w http.ResponseWriter, r *http.Request, adapterKey string) {
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		adapter.RenderForbidden(w, r, s.NewSystemResponse().StatusForbidden())
	}
}

// RenderMethodNotAllowed renders a method not allowed page
func (s *HyperView) RenderMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	s.RenderMethodNotAllowedAs(w, r, "html")
}

// RenderMethodNotAllowedAs renders a method not allowed page as the specified adapter
func (s *HyperView) RenderMethodNotAllowedAs(w http.ResponseWriter, r *http.Request, adapterKey string) {
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		adapter.RenderMethodNotAllowed(w, r, s.NewSystemResponse().Status(http.StatusMethodNotAllowed))
	}
}

// RenderUnauthorized renders an unauthorized page
func (s *HyperView) RenderUnauthorized(w http.ResponseWriter, r *http.Request) {
	s.RenderUnauthorizedAs(w, r, "html")
}

// RenderUnauthorizedAs renders an unauthorized page as the specified adapter
func (s *HyperView) RenderUnauthorizedAs(w http.ResponseWriter, r *http.Request, adapterKey string) {
	if adapter, ok := s.adapterFor(w, adapterKey); ok {
		adapter.RenderUnauthorized(w, r, s.NewSystemResponse().StatusUnauthorized())
	}
}

// HxRedirect sends an HX-Redirect header to the client
func (s *HyperView) HxRedirect(w http.ResponseWriter, url string) {
	w.Header().Set(htmx.HXRedirect, url)
	w.WriteHeader(http.StatusSeeOther)
	_, _ = w.Write([]byte("redirecting..."))
	return
}

// Redirect sends a redirect response to the client
func (s *HyperView) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	if htmx.IsHtmxRequest(r) {
		s.HxRedirect(w, url)
		return
	} else if request.IsXMLHttpRequest(r) {
		// Create a JSON response with a redirect
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := map[string]string{
			"status":  "redirect",
			"message": "redirecting...",
			"url":     url,
		}

		jsonBytes, _ := json.Marshal(data)
		_, _ = w.Write(jsonBytes)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// NewResponse creates a new response with the given layout
func (s *HyperView) NewResponse(layout string) *response.Response {
	return response.NewResponse().Layout(layout)
}

// NewSystemResponse creates a new response with the system layout
func (s *HyperView) NewSystemResponse() *response.Response {
	return response.NewResponse().Layout(s.systemLayout)
}

// adapterFor returns the adapter for the specified key
func (s *HyperView) adapterFor(w http.ResponseWriter, key string) (Adapter, bool) {
	if key == "" {
		key = "html"
	}

	adapter, ok := s.Adapter(key)
	if !ok {
		http.Error(w, "Adapter not found", http.StatusInternalServerError)
		return nil, false
	}
	return adapter, true
}
