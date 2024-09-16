package handlers

import (
	"encoding/json"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/models"

	"github.com/gorilla/mux"
)

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