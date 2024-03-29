package html

import (
	"io"
	"text/template"

	"github.com/andremfp/snippetbox/internal/database"
)

func RenderTemplate(w io.Writer, htmlFiles []string, data *database.Snippet) error {

	templateSet, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		return err
	}

	err = templateSet.ExecuteTemplate(w, "base", data)
	if err != nil {
		return err
	}

	return nil
}
