package server

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/andremfp/snippetbox/internal/middleware"
	"github.com/andremfp/snippetbox/internal/templates"
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
	mux := http.NewServeMux()

	staticDir, err := fs.Sub(templates.Content, "ui/static")
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	staticFileHandler := http.FileServer(http.FS(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", staticFileHandler))

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/snippet/view", app.snippetViewHandler)
	mux.HandleFunc("/snippet/create", app.snippetCreateHandler)

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, middleware.SecureHeaders)

	return standardMiddleware.Then(mux)
}
