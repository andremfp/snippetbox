package server

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/julienschmidt/httprouter"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	SnippetStore  database.Store
	TemplateCache map[string]*template.Template
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.SnippetStore.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.Render(w, http.StatusOK, "home.html", data)
}

func (app *Application) snippetViewHandler(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(r)

	data.Snippet = snippet

	app.Render(w, http.StatusOK, "view.html", data)

}

func (app *Application) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *Application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request) {

	// Dummy data for now
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.SnippetStore.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
