package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	webserver := &Webserver{}
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, webserver)
	log.Fatal(err)
}
