package utils

import (
	"io"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/labstack/echo/v4"
)

func JSONError(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}

// TemplateRenderer implements echo.Renderer interface
type TemplateRenderer struct {
	templates map[string]*template.Template
	mu        sync.Mutex
}

// NewTemplateRenderer initializes the renderer and template cache
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: make(map[string]*template.Template),
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	tmpl, ok := t.templates[name]
	if !ok {
		// Load and cache the template if not already
		parsedTemplate, err := template.ParseFiles(filepath.Join("templates", name))
		if err != nil {
			return err
		}
		t.templates[name] = parsedTemplate
		tmpl = parsedTemplate
	}

	return tmpl.Execute(w, data)
}
