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

	"github.com/alexedwards/scs/mysqlstore" 	
	"github.com/alexedwards/scs/v2"         
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

	//Initialize a new http.Server struct. We set the Addr and Handler fields type so 
	//that the server uses the same network address and routes as before

	server := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		//Create a *log.Logger from our structured logger handler, which writes 
		//log entries at the Error level, and assign it to the ErrorLog field. If 
		//you would prefer to log on the server errors at Warn level instead, you 
		//could pass slog.LevelWarn as the final paramter. 
		ErrorLog: slog.NewLogLogger(logger.Handler(),slog.LevelError),
			
	}

	logger.Info("starting server", "addr", server.Addr)

	
	//Call the ListenAndServe() method on our new http.Server struct to start
	//the server 
	err = server.ListenAndServe()
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
