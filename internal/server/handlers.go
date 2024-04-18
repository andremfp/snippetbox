package server

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/validator"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	SnippetStore  database.Store
	TemplateCache map[string]*template.Template
	FormDecoder   *form.Decoder
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.Render(w, http.StatusOK, "create.html", data)
}

func (app *Application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusSeeOther, "create.html", data)
		return
	}

	id, err := app.SnippetStore.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
