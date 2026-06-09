package main

import (
	"database/sql" //New Import
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hqr999/Snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore" //New import
	"github.com/alexedwards/scs/v2"         //New import
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

//Add a new sessionManager field to the application struct
type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder 
	sessionMangaer *scs.SessionManager
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

	cachePageTmpl, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	formDecoder := form.NewDecoder()

	//Use the scs.New() function to initialize a new session manager. Then we 
	//configure it to use our MySQL database as the session store, and set a 
	//lifetime of 12 hours (so that sessions automatically expire 12 hours 
	//after first being created). 
	sessionMan := scs.New()
	sessionMan.Store = mysqlstore.New(db)
	sessionMan.Lifetime = 12 * time.Hour




	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
		templateCache: cachePageTmpl,
		formDecoder: formDecoder,
		sessionMangaer: sessionMan,
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
