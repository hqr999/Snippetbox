package main

import (
	"crypto/tls" //New import 
	"database/sql" 
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

// Add a new sessionManager field to the application struct
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
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
	sessionMan.Cookie.Secure = true

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  cachePageTmpl,
		formDecoder:    formDecoder,
		sessionMangaer: sessionMan,
	}
	//Initialize a tls.Config struct to hold the non-default TLS settings we
	//want the server to use. In this case the only thing that we´re changing
	//is the curve preference value, so that only elliptic curves with 
	//assembly implementations are used. 
	tlsConf := &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//See the servers TLSConfig field to use the tlsConfig variable we just 
	//created.
	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConf,
	}

	logger.Info("starting server", "addr", server.Addr)

	//Use the ListenAndServeTLS() method to start the HTTPS server. We
	//pass in the paths to the TLS certificates and corresponding private key as
	//the two parameters.
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
