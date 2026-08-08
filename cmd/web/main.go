package main

import (
	"crypto/tls" //New import
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hqr999/Snippetbox/internal/models"
	"github.com/pressly/goose/v3"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// Add a new sessionManager field to the application struct
type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface //Use our new interface type.
	users          models.UserModelInterface    //Use our new interfaceype.
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Connection to MySQL running on Docker.
	dsn := flag.String("dsn", "web:gintoki@tcp(localhost:3308)/snippetbox?parseTime=true", "MySQL data source name")
	migrate := flag.String("migrate", "up", "Command for our Database migration tool (Goose)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	err = migrateDbCommand(*migrate, *dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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
		users:          &models.UserModel{DB: db},
		templateCache:  cachePageTmpl,
		formDecoder:    formDecoder,
		sessionManager: sessionMan,
	}

	tlsConf := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConf,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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

func migrateDbCommand(command, dsn string) error {
	err := goose.SetDialect("mysql")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}
	switch command {
	case "up":
		err = goose.Up(db, "migrations")
	case "status":
		err = goose.Status(db, "migrations")
	case "down":
		err = goose.Down(db, "migrations")
	case "fix":
		err = goose.Fix("migrations")
	case "redo":
		err = goose.Redo(db, "migrations")
	case "version":
		err = goose.Version(db, "migrations")
	default:
		return fmt.Errorf("Invalid goose commando: %s", command)
	}

	if err != nil {
		db.Close()
		return err
	}
	return nil
}
