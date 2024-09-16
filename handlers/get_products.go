package handlers

import (
	"encoding/json"
	"net/http"
	"simpler-test-api/db"
	"simpler-test-api/helper"
	"simpler-test-api/models"
	"strconv"
)

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

    limit := helper.GetEnvInt("DEFAULT_PAGINATION_LIMIT", 10)
    offset := helper.GetEnvInt("DEFAULT_PAGINATION_OFFSET", 0)

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

    var total int
    err := db.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&total)
    if err != nil {
        http.Error(w, "Error fetching total count", http.StatusInternalServerError)
        return
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

    response := struct {
        Total int `json:"total"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Products []models.Product `json:"products"`
    }{
        Total: total,
        Limit: limit,
        Offset: offset,
        Products: products,
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

