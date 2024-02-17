package main

import "net/http"

type Webserver struct{}

// server.go
func (p *Webserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mux := http.NewServeMux()
	staticFileHandler := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)
	mux.Handle("/static/", http.StripPrefix("/static", staticFileHandler))

	mux.ServeHTTP(w, r)
}
