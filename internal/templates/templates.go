package templates

import (
	"embed"
	"html/template"
	"path/filepath"
	"time"

	"github.com/andremfp/snippetbox/internal/database"
)

//go:embed ui
var Content embed.FS

type TemplateData struct {
	CurrentYear int
	Snippet     *database.Snippet
	Snippets    []*database.Snippet
	Form        any
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages := []string{
		"ui/html/pages/home.html",
		"ui/html/pages/view.html",
		"ui/html/pages/create.html",
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"ui/html/base.html",
			"ui/html/partials/nav.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(Content, files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
