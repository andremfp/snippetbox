package main

import (
	"flag"
	"log"
	"os"

	"github.com/andremfp/snippetbox/internal/html/config"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &config.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	webserver := NewWebserver(*addr, errorLog, app)

	// svr := &http.Server{}
	infoLog.Printf("Starting server on %s", *addr)
	err := webserver.ListenAndServe()
	errorLog.Fatal(err)
}
