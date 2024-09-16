package main

import (
	"log"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

/*
	main() is, as one can easily guess, the application's entry point.
	It is mandatory and does not return anything.
	Typically used to set up the environment, initialize global resources and handle command
	line arguments.
	Here we :
	1) specify our .env file which holds our necessary credentials for the db connection,
	2) initialize our database, and
	3) register our routes/handlers needed to process our product records.
*/
func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    db.InitializeDB()

	r := mux.NewRouter()

	routes.RegisterProductRoutes(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
