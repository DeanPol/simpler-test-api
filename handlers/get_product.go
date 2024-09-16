package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/models"

	"github.com/gorilla/mux"
)


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