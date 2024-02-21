package main

import (
	"log"
	"net/http"
)

type Webserver http.Server

func NewWebserver(addr string, errorLog *log.Logger, app *application) *http.Server {

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  NewServeMux(app),
	}

	return srv
}

func NewServeMux(app *application) http.Handler {
	mux := http.NewServeMux()
	staticFileHandler := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/snippet/view", app.snippetViewHandler)
	mux.HandleFunc("/snippet/create", app.snippetCreateHandler)
	mux.Handle("/static/", http.StripPrefix("/static", staticFileHandler))

	return mux
}
