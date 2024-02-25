package main

import (
	"flag"
	"log"
	"os"

	"github.com/andremfp/snippetbox/internal/database"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:snippetbox_dev@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := database.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	webserver := NewWebserver(*addr, errorLog, app)

	// svr := &http.Server{}
	infoLog.Printf("Starting server on %s", *addr)
	err = webserver.ListenAndServe()
	errorLog.Fatal(err)
}
