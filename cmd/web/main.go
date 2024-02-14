package main

import (
	"log"
	"net/http"
)

func main() {

	webserver := &Webserver{}
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", webserver)
	log.Fatal(err)
}
