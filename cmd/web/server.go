package main

import (
	"log"
	"net/http"

	"github.com/andremfp/snippetbox/internal/html/config"
)

type Webserver http.Server

func NewWebserver(addr string, errorLog *log.Logger, app *config.Application) *http.Server {

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  NewServeMux(app),
	}

	return srv
}

func NewServeMux(app *config.Application) http.Handler {
	mux := http.NewServeMux()
	staticFileHandler := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r, app)
	})
	mux.HandleFunc("/snippet/view", func(w http.ResponseWriter, r *http.Request) {
		snippetViewHandler(w, r, app)
	})
	mux.HandleFunc("/snippet/create", func(w http.ResponseWriter, r *http.Request) {
		snippetCreateHandler(w, r, app)
	})
	mux.Handle("/static/", http.StripPrefix("/static", staticFileHandler))

	return mux
}
