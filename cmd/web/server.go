package main

import (
	"log"
	"net/http"
)

type Webserver http.Server

func NewWebserver(addr string, errorLog *log.Logger) *http.Server {

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  NewServeMux(),
	}

	return srv
}

func NewServeMux() http.Handler {
	mux := http.NewServeMux()
	staticFileHandler := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)
	mux.Handle("/static/", http.StripPrefix("/static", staticFileHandler))

	return mux
}
