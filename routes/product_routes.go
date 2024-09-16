package routes

import (
	"simpler-test-api/handlers"

	"github.com/gorilla/mux"
)

/*
   Register our routes and CRUD operations.
   We are using gorilla/mux, a routing system that seems widely popular.
   Apparently, it provides many benefits over the 'vanilla' way of things, such as:
   1) Path variables, as seen with '/products/{id}'
   2) Allows for request processing before they reach the handlers, useful for authentication and logging
   3) Method-specific matching. I think this is being able to handle different HTTP methods for the same path (POST, GET for '/products')
   4) Route matching with regex, as seen in our GET, PUT and DELETE methods
*/

// HandleFunc allows us to register our functions to the HTTP methods with the given endpoint.
func RegisterProductRoutes(r *mux.Router) {
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.GetProductById).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
}
