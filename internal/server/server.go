package server

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/andremfp/snippetbox/internal/middleware"
	"github.com/andremfp/snippetbox/internal/templates"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Webserver http.Server

func NewWebserver(addr string, errorLog *log.Logger, app *Application) *http.Server {

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.NewServeMux(),
	}

	return srv
}

func (app *Application) NewServeMux() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	staticDir, err := fs.Sub(templates.Content, "ui/static")
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	staticFileHandler := http.FileServer(http.FS(staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", staticFileHandler))

	router.HandlerFunc(http.MethodGet, "/", app.HomeHandler)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetViewHandler)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateHandler)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePostHandler)

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, middleware.SecureHeaders)

	return standardMiddleware.Then(router)
}
