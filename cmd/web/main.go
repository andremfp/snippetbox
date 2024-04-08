package main

import (
	"flag"
	"log"
	"os"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/server"
	"github.com/andremfp/snippetbox/internal/templates"
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

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &server.Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		SnippetStore:  &database.SnippetModel{DB: db},
		TemplateCache: templateCache,
	}

	webserver := server.NewWebserver(*addr, errorLog, app)

	infoLog.Printf("Starting server on %s", *addr)
	err = webserver.ListenAndServe()
	errorLog.Fatal(err)
}
