package handlers

import (
	"net/http"
	"simpler-test-api/db"

	"github.com/gorilla/mux"
)

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}