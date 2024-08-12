package hypergo

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func (a *TemplateAdapter) Render(w http.ResponseWriter, r *http.Request, resp *Response) {
	tmpl, err := a.getTemplate(resp)
	if err != nil {
		a.handleError(w, r, err)
		return
	}

	a.execTemplate(w, r, resp, tmpl)
}

func (a *TemplateAdapter) RenderForbidden(w http.ResponseWriter, r *http.Request, resp *Response) {
	path := a.viewsPath(SystemDir, "403")
	if _, ok := a.templates[path]; ok {
		a.Render(w, r, resp.Path(path))
		return
	}
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func (a *TemplateAdapter) RenderMaintenance(w http.ResponseWriter, r *http.Request, resp *Response) {
	path := a.viewsPath(SystemDir, "503")
	if _, ok := a.templates[path]; ok {
		a.Render(w, r, resp.Path(path))
		return
	}
	http.Error(w, "Maintenance", http.StatusServiceUnavailable)
}

func (a *TemplateAdapter) RenderMethodNotAllowed(w http.ResponseWriter, r *http.Request, resp *Response) {
	path := a.viewsPath(SystemDir, "405")
	if _, ok := a.templates[path]; ok {
		a.Render(w, r, resp.Path(path))
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (a *TemplateAdapter) RenderNotFound(w http.ResponseWriter, r *http.Request, resp *Response) {
	path := a.viewsPath(SystemDir, "404")
	if _, ok := a.templates[path]; ok {
		a.Render(w, r, resp.Path(path))
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (a *TemplateAdapter) RenderSystemError(w http.ResponseWriter, r *http.Request, err error, resp *Response) {
	// Get the stack trace and output to the log
	a.logger.Error("Server error", slog.String("err", err.Error()))
	lineErrors := ""
	lines := strings.Split(string(debug.Stack()), "\n")
	for i, line := range lines {
		// replace \t with 4 spaces
		line = strings.ReplaceAll(line, "\t", "    ")
		lineErrors += fmt.Sprintf("--- traceLine%03d: %s\n", i, line)
		a.logger.Error("Stack trace", slog.String(fmt.Sprintf("--- traceLine%03d", i), line))
	}

	// If there is a template with the name "system/server_error" in the template cache, use it
	path := a.viewsPath(SystemDir, "500")
	if _, ok := a.templates[path]; ok {
		resp.Path(path).
			Errors(err.Error(), map[string]string{"LineErrors": lineErrors}).
			StatusError()
		a.Render(w, r, resp)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (a *TemplateAdapter) RenderUnauthorized(w http.ResponseWriter, r *http.Request, resp *Response) {
	path := a.viewsPath(SystemDir, "401")
	if _, ok := a.templates[path]; ok {
		a.Render(w, r, resp.Path(path))
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (a *TemplateAdapter) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *TemplateAdapter) getTemplate(resp *Response) (*template.Template, error) {
	// Retrieve the preloaded page template from the cache
	pageTmpl, ok := a.templates[resp.TemplatePath()]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", resp.TemplatePath())
	}

	// Clone the page template to avoid altering the original cached template
	tmpl, err := pageTmpl.Clone()
	if err != nil {
		return nil, fmt.Errorf("error cloning template: %w", err)
	}

	// Combine the page with its layout template from the embedded filesystem
	layoutPath := LayoutsDir + "/" + resp.TemplateLayout() + a.extension
	tmpl, err = tmpl.ParseFS(a.fileSystemMap[RootFSID], layoutPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing layout: %w", err)
	}

	return tmpl, nil
}

func (a *TemplateAdapter) execTemplate(w http.ResponseWriter, r *http.Request, resp *Response, tmpl *template.Template) {
	// Creating a buffer, so we can capture write errors before we write to the header
	// Note that layouts are always defined as "layout" in the templates
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "layout", resp.ViewData(r).Data())
	if err != nil {
		path := a.viewsPath(SystemDir, "server-error")
		if resp.TemplatePath() == path {
			http.Error(w, fmt.Errorf("error executing template: %w", err).Error(), http.StatusInternalServerError)
		} else {
			a.handleError(w, r, fmt.Errorf("error executing template: %w", err))
		}
		return
	}

	// Add any additional headers
	for key, value := range resp.Headers() {
		w.Header().Set(key, value)
	}

	// Set the status code
	w.WriteHeader(resp.StatusCode())

	// Write the buffer to the response
	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *TemplateAdapter) viewsPath(path ...string) string {
	// For each path, append to the ViewsDir, separated by a slash
	return fmt.Sprintf("%s/%s", ViewsDir, strings.Join(path, "/"))
}
