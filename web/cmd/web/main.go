package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"encoding/gob"
	"github.com/alexedwards/scs/v2"

	"github.com/mcgigglepop/cloudpilot/web/internal/config"

)


const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	// run function
	db, err := run()

	// error handle
	if err != nil {
		log.Fatal(err)
	}

	// close the db connection
	defer db.SQL.Close()

	// print
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	// serve the application
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// listen and serve
	err = srv.ListenAndServe()
	
	// error log
	log.Fatal(err)
}

// run is the primary run function
func run() (*driver.DB, error) {
	// load into the session
	gob.Register(models.User{})
	gob.Register(map[string]int{})

	// runtime flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "password", "Database password")
	dbPort := flag.String("dbport", "5433", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	// parse flags
	flag.Parse()

	// db flag check
	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	// in production and use cache
	app.InProduction = *inProduction
	app.UseCache = *useCache

	// info log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// error log
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// session info
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	
	// connection string
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	
	// error handle
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	// print
	log.Println("Connected to database!")

	// create template cache
	tc, err := render.CreateTemplateCache()
	
	// error handle
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	// repo pattern
	repo := handlers.NewRepo(&app, db)
	
	// handlers, render, helpers
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// return
	return db, nil
}