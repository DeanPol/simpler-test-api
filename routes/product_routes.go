package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/models"
	"strconv" // convert a string to int

	"os"

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
	r.HandleFunc("/products/{id:[0-9]+}", GetProductById).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", GetProducts).Methods("GET")
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

func GetProductById(w http.ResponseWriter, r *http.Request) {
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
/*
    w.WriteHeader() sends an HTTP response status code to the client. The status code indicates the result of the HTTP request.
    In this example, http.StatusNoContent corresponds to the status code 204 (request successful but no content to send in the response body)
*/
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

/*
    In the following method we allow for pagination by retrieving the 'limit' and 'offset' url params.
    If none are provided, we use our environment variables. If something goes wrong, default values are used instead.
*/

/*
    rows.Close(): Closes the database result set. Not closing the rows can lead to memory leaks or a situation
    where too many open connections to the database accumulate.

    defer: keyword that schedules the execution of the following function to occur after the surrounding function has finished execution.
    No matter how the function exits, the rows.Close() will be called. (Similar to 'finaly' in 'try/catch')
*/

func GetProducts(w http.ResponseWriter, r *http.Request) {
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")

    limit := getEnvInt("DEFAULT_PAGINATION_LIMIT", 10)
    offset := getEnvInt("DEFAULT_PAGINATION_OFFSET", 0)

    if limitStr != "" {
        parsedLimit, err := strconv.Atoi(limitStr)
        if err != nil || parsedLimit <= 0 {
            http.Error(w, "Invalid limit value", http.StatusBadRequest)
            return
        }
        limit = parsedLimit
    }

    if offsetStr != "" {
        parsedOffset, err := strconv.Atoi(offsetStr)
        if err != nil || parsedOffset < 0 {
            http.Error(w, "Invalid offset value", http.StatusBadRequest)
            return
        }
        offset = parsedOffset
    }

    query := `SELECT id, name, description, price, stock FROM products LIMIT $1 OFFSET $2`
	rows, err := db.DB.Query(query, limit, offset)
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

// Helper function to fetch .env variable with fallback value
// TODO: If more helper functions, I should create a new directory
func getEnvInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}
