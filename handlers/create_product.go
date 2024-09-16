package handlers

import (
	"encoding/json"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/models"
)

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
