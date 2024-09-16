package main

import (
	"log"
	"simpler-test-api/db"

	"github.com/joho/godotenv"
)

/*
	main() is, as one can easily guess, the application entry point.
	It is mandatory and does not return anything.
	Typically used to set up the environment, initialize global resources and handle command
	line arguments.
	Here we just specify our .env file which holds our necessary credentials for the db connection, and initialize our database.
*/
func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    db.InitializeDB()
}
