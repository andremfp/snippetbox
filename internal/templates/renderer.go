package templates

import (
	"embed"
	"io"
	"text/template"

	"github.com/andremfp/snippetbox/internal/database"
)

//go:embed ui
var Content embed.FS

func RenderTemplate(w io.Writer, htmlFiles []string, data *database.Snippet) error {

	templateSet, err := template.ParseFS(Content, htmlFiles...)
	if err != nil {
		return err
	}

	err = templateSet.ExecuteTemplate(w, "base", data)
	if err != nil {
		return err
	}

	return nil
}
