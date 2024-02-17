package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andremfp/snippetbox/internal/html"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent / from being catch all
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	htmlFiles := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	err := html.RenderTemplate(w, htmlFiles)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	w.WriteHeader(http.StatusOK)
}

func snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
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
