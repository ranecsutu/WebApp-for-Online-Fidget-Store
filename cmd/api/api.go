package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ranecsutu/fidget/internal/driver"
	"github.com/ranecsutu/fidget/internal/models"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	secretkey string
	frontend  string
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting Back End server in %s mode on port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.db.dsn, "dsn", "", "URL to dsn")
	flag.StringVar(&cfg.smtp.host, "smtphost", "sandbox.smtp.mailtrap.io", "smtp host")
	flag.StringVar(&cfg.smtp.username, "smtpuser", "1bf61a987cb9f4", "smtp user")
	flag.StringVar(&cfg.smtp.password, "smtppassword", "40161f4f0c8bb9", "smtp password")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "smtp port")
	flag.StringVar(&cfg.secretkey, "secret", "01234567890123456789012345678901", "secret key") //secret key needs to be exactly 32 letters long
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to frontend")

	/* Read stripe key&secret from flags (if any) */
	flag.StringVar(&cfg.stripe.key, "STRIPE_KEY", "", "Stripe key")
	flag.StringVar(&cfg.stripe.secret, "STRIPE_SECRET", "", "Stripe secret")
	flag.Parse()

	/* Read stripe key&secret from env (if they weren't already declared in the CLI flags) */
	if cfg.stripe.key == "" {
		cfg.stripe.key = os.Getenv("STRIPE_KEY")
	}
	if cfg.stripe.secret == "" {
		cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
