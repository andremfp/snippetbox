package main

import "net/http"

type Webserver struct{}

// server.go
func (p *Webserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)

	mux.ServeHTTP(w, r)
}
