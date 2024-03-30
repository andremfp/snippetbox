package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/templates"
)

type Application struct {
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	SnippetStore database.Store
}

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent / from being catch all
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.SnippetStore.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := templates.TemplateData{
		Snippets: snippets,
	}

	htmlFiles := []string{
		"ui/html/base.html",
		"ui/html/partials/nav.html",
		"ui/html/pages/home.html",
	}

	err = templates.RenderTemplate(w, htmlFiles, data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *Application) snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.SnippetStore.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := templates.TemplateData{
		Snippet: snippet,
	}

	htmlFiles := []string{
		"ui/html/base.html",
		"ui/html/partials/nav.html",
		"ui/html/pages/view.html",
	}

	err = templates.RenderTemplate(w, htmlFiles, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

func (app *Application) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy data for now
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.SnippetStore.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
