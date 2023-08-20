package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/ahmadyogi543/snippetbox/internal/models"
	"github.com/ahmadyogi543/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	User            *models.User
}

var templateFunctions = template.FuncMap{
	"humanDate": formatHumanReadableDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.go.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.go.html",
			"html/partials/*.go.html",
			page,
		}

		htmlTemplate := template.New(name).Funcs(templateFunctions)
		ts, err := htmlTemplate.ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
