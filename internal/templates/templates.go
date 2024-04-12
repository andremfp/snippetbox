package templates

import (
	"embed"
	"html/template"
	"path/filepath"

	"github.com/andremfp/snippetbox/internal/database"
)

//go:embed ui
var Content embed.FS

type TemplateData struct {
	Snippet  *database.Snippet
	Snippets []*database.Snippet
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages := []string{
		"ui/html/pages/home.html",
		"ui/html/pages/view.html",
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"ui/html/base.html",
			"ui/html/partials/nav.html",
			page,
		}
		ts, err := template.ParseFS(Content, files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}