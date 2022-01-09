package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// simple trick to clear the cache and force an update
const version = "1.0.0"
const cssVersion = "1.0"

// configuration information for aplication, write once, use many times
type config struct {
	port int
	env  string // environment
	api  string
	pi   string
	db   struct { // database
		dsn string // data source name
	}
	stripe struct { // optional, for stripe management
		secret string // secret api key
		key    string // public api key
	}
}

// application data
type application struct {
	config        config      // the config struct write above
	infoLog       *log.Logger // * means return value of the pointer
	errorLog      *log.Logger
	templateCache map[string]*template.Template // map, string as key and return value of the pointer
	version       string
}

func (app *application) serve() error { // the asterisk is because you want to modify the information that is in memory, pass by reference
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server in %s on port %d", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config

	// adding specific attributes to the application, such as port, type of env, api...
	flag.IntVar(&cfg.port, "port", 8080, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "app.environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:8081/", "url to api")
	flag.StringVar(&cfg.pi, "pi", "http://localhost:8081/api/payment-intent", "url to pi")
	flag.Parse()

	// prevents the key from being accesible, hidden in environment os
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// template cache
	tc := make(map[string]*template.Template)

	app := &application{ // ampersand is used because you want to modify the application parameters and not get a copy of application
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

}
