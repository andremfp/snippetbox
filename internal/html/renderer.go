package html

import (
	"io"
	"text/template"
)

func RenderTemplate(w io.Writer, htmlFiles []string) error {

	templateSet, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		return err
	}

	err = templateSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		return err
	}

	return nil
}
