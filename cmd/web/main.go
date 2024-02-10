package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/dinizgab/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
    templateCache map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB("web:123456@/snippetbox?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
    }
	defer db.Close()
   
    templateCache, err := newTemplateCache()
    if err != nil {
        errorLog.Fatal(err) 
    } 

	app := &Application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
        templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     ":4000",
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Println("Starting server on :4000")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(connString string) (*sql.DB, error) {
	db, err := sql.Open(
		"mysql", connString,
	)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
