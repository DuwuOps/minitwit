package template_rendering

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFiles(filepath.Join("templates", name)))
	return tmpl.ExecuteTemplate(w, name, data)
}

// Create and return a new instance of a TemplateRenderer.
func NewTemplateRenderer() *TemplateRenderer {
	funcMap := template.FuncMap{
		"gravatar":       gravatarUrl,
		"datetimeformat": formatDatetime,
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob(filepath.Join("templates", "*.html")))
	return &TemplateRenderer{
		templates: tmpl,
	}
}
