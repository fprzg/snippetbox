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

	"snippetbox.fepg.org/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	debugMode      bool
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type config struct {
	addr      string
	dsn       string
	TLSCert   string
	TLSKey    string
	debugMode bool
}

var cfg config
var app *application

func main() {
	//addr := flag.String("addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.TLSCert, "tls-cert", "./tls/cert.pem", "TLS cert.pem directory")
	flag.StringVar(&cfg.TLSKey, "tls-key", "./tls/key.pem", "TLS key.pem directory")
	flag.BoolVar(&cfg.debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	app = &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	// log.Llongfile -> long file names
	// log.LUTC -> utc time zone

	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLog.Fatal("Error parsing templates.")
	} else {
		app.infoLog.Printf("Templates parsed successfully.")
	}

	db, err := openDB(cfg.dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app.snippets = &models.SnippetModel{DB: db}
	app.users = &models.UserModel{DB: db}
	app.templateCache = templateCache
	app.formDecoder = formDecoder
	app.sessionManager = sessionManager

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:           cfg.addr,
		ErrorLog:       app.errorLog,
		Handler:        app.routes(),
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	app.infoLog.Printf("Starting server on %s", cfg.addr)
	//err := http.ListenAndServe(cfg.addr, mux)
	err = srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
