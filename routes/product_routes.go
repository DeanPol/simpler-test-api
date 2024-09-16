package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/models"

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
	r.HandleFunc("/products", CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", GetProduct).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", ListProducts).Methods("GET")
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
    /* 
        json.NewDecoder creates a json.Decoder which allows for streaming JSON decoding. THis is useful for 
        when dealing with large JSON payloads because it doesn't require the entire payload to be loaded into 
        memory at once, rather, it instead reads and decodes the JSON as it streams in from the request body.
        Apparently, Decode returns an error if the JSON provided is malformed or if there is a type mismatch or missing field.
    */
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Some data entry validation, just to see how it works.
    if product.Name == "" || product.Price < 0 {
        http.Error(w, "Invalid product data", http.StatusBadRequest)
        return
    }

	_, err := db.DB.Exec("INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4)", product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var product models.Product
	row := db.DB.QueryRow("SELECT id, name, description, price, stock FROM products WHERE id = $1", id)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4 WHERE id = $5", product.Name, product.Description, product.Price, product.Stock, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, description, price, stock FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}
