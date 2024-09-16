package main

import (
	"log"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/routes"

	"github.com/gorilla/handlers"
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
	2) initialize our database,
	3) - since 0.0.2 - register our routes/handlers needed to process our product records.
	4) - since 0.0.4 - include CORS middleware
*/
func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    db.InitializeDB()

	r := mux.NewRouter()

	routes.RegisterProductRoutes(r)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	log.Println("Server starting on port 8080...")
	
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
