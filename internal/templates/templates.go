package templates

import "github.com/andremfp/snippetbox/internal/database"

type TemplateData struct {
	Snippet  *database.Snippet
	Snippets []*database.Snippet
}
