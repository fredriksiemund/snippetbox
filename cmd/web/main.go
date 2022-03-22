package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"fredriksiemund/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// The application struct containing all of our dependencies
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetRepository
	templateCache map[string]*template.Template
}

func main() {
	// Parsing the runtime configuration settings for the application
	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "web:password@/snippetbox?parseTime=true", "MySQL connection string")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session cookies secret key")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDb(*connStr)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Establishing the dependencies for the handlers (depenency injection)
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetRepository{DB: db},
		templateCache: templateCache,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want
	// the server to use. In this case the only thing that we're changing is the
	// curve preferences value, so that only elliptic curves with assembly
	// implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Running the HTTP server
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDb(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	// Since connections to the database are established lazily,
	// we can verify that everything is set up correctly by calling db.Ping()
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
