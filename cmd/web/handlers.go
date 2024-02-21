package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andremfp/snippetbox/internal/html"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent / from being catch all
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	htmlFiles := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	err := html.RenderTemplate(w, htmlFiles)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
