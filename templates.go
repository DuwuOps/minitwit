package main

import (
    "html/template"
    "io"
    "path/filepath"

    "github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
    templates *template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
    tmpl := template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
    return &TemplateRenderer{
        templates: tmpl,
    }
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}
