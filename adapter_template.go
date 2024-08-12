package hypergo

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/hypergopher/hypergo/funcs"
)

// TemplateAdapter is a template adapter for the HyperGo framework that uses the Go html/template package.
type TemplateAdapter struct {
	extension     string
	fileSystemMap map[string]fs.FS
	logger        *slog.Logger
	funcMap       template.FuncMap
	templates     map[string]*template.Template
}

// TemplateViewAdapterOptions are the options for the TemplateAdapter.
type TemplateViewAdapterOptions struct {
	// Extension is the file extension for the templates. Default is ".gtml".
	Extension string
	// FileSystemMap is a map of file systems to use for the templates.
	FileSystemMap map[string]fs.FS
	// Funcs is a map of functions to add to the template.FuncMap.
	Funcs template.FuncMap
	// Logger is the logger to use for the adapter.
	Logger *slog.Logger
}

// NewTemplateViewAdapter creates a new TemplateAdapter.
func NewTemplateViewAdapter(opts TemplateViewAdapterOptions) *TemplateAdapter {
	// Merge the other functions into the base template functions
	for k, v := range opts.Funcs {
		funcs.FuncMap[k] = v
	}

	if opts.Extension == "" {
		opts.Extension = ".gtml"
	}

	return &TemplateAdapter{
		extension:     opts.Extension,
		fileSystemMap: opts.FileSystemMap,
		funcMap:       funcs.FuncMap,
		logger:        opts.Logger,
		templates:     make(map[string]*template.Template),
	}
}

func (a *TemplateAdapter) Init() error {
	// Reset the template cache
	a.templates = make(map[string]*template.Template)

	baseTemplate, err := a.loadPartials()
	if err != nil {
		return fmt.Errorf("error loading partials. %w", err)
	}

	// Function to recursively process directories from all FileSystemMap
	for fsID, fsys := range a.fileSystemMap {
		processDirectory := func(path string, dir fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !dir.IsDir() && filepath.Ext(path) == a.extension {
				relPath, err := filepath.Rel("", path)
				if err != nil {
					return err
				}
				pageName := strings.TrimSuffix(relPath, filepath.Ext(relPath))
				if fsID != RootFSID {
					pageName = fsID + ":" + pageName
				}

				// Clone the base template and parse the page template
				tmpl, err := template.Must(baseTemplate.Clone()).ParseFS(fsys, path)
				if err != nil {
					return err
				}
				a.templates[pageName] = tmpl
			}
			return nil
		}

		// If the "views" directory exists, parse it. Otherwise, parse the root directory
		if _, err := fsys.Open(ViewsDir); err == nil {
			if err := fs.WalkDir(fsys, ViewsDir, processDirectory); err != nil {
				return err
			}
		}
	}
	// Uncomment to view the template names found
	// a.printTemplateNames()

	return nil
}

func (a *TemplateAdapter) loadPartials() (*template.Template, error) {
	baseTemplate := template.New("base").Funcs(a.funcMap)

	for _, fsys := range a.fileSystemMap {
		processPartials := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(path) == a.extension {
				fullPath := path
				_, err := baseTemplate.ParseFS(fsys, fullPath)
				if err != nil {
					return err
				}
			}
			return nil
		}

		// If the "partials" directory exists, parse it
		if _, err := fsys.Open(PartialsDir); err == nil {
			if err := fs.WalkDir(fsys, PartialsDir, processPartials); err != nil {
				return nil, err
			}
		}
	}

	return baseTemplate, nil
}

func (a *TemplateAdapter) printTemplateNames() {
	for name, tmpl := range a.templates {
		fmt.Printf("Template: %s\n", name)
		associatedTemplates := tmpl.Templates()
		for _, tmpl := range associatedTemplates {
			fmt.Printf("\tPartial/Child: %s\n", tmpl.Name())
		}
	}
}
