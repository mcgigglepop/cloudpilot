package main

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