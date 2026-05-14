package main

import (
	"database/sql" //New Import
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/hqr999/Snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql" //New Import
)

// Add a snippets field to the application struct. This will
// allow us to make the application the SnippetModel available to our handlers.
type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:gintoki@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	//Initialize a new template cache... 
	cachePageTmpl, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//And add it to the application dependencies 
	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
		templateCache: cachePageTmpl,
	}

	logger.Info("starting server", "addr", *addr)

	
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
