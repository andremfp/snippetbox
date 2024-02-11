package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent / from being catch all
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
	w.WriteHeader(http.StatusOK)
}

func snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
	w.WriteHeader(http.StatusOK)
}

func snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
	w.WriteHeader(http.StatusOK)
}
