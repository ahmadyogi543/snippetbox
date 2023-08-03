package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ahmadyogi543/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	// get flags from the command line and parse it.
	addr := flag.String("addr", ":3000", "HTTP network address")
	dsn := flag.String("dsn", "web:123456@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// custom loggers (info and error) to log info and error to stdout and stderr.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// open DB connection and close it before the main function exists with defer.
	db, err := openDB("mysql", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// session manager setup.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// initialize app struct for containing the dependencies.
	app := &App{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	// initialize a new http.Server so that we can use custom error log.
	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// listen and serve TLS with address from flags addr and with the servemux above.
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
