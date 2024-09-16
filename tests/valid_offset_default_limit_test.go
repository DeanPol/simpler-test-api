package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simpler-test-api/db"
	"simpler-test-api/handlers"
	"simpler-test-api/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

/*
   Mock the product fetching query response (using default limit 10, offset=5)
   Assert that the total count, limit, and offset are correct
*/
func TestGetProducts_ValidOffsetDefaultLimit(t *testing.T) {
    dbConn, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening mock db: %v", err)
    }
    defer dbConn.Close()

    db.DB = dbConn

    mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))

    mock.ExpectQuery("SELECT id, name, description, price, stock FROM products").
        WithArgs(10, 5).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock"}).
            AddRow(6, "Product 6", "Description 6", 200.0, 3).
            AddRow(7, "Product 7", "Description 7", 250.0, 2))

    req, err := http.NewRequest("GET", "/products?offset=5", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetProducts)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var response struct {
        Total int `json:"total"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Products []models.Product `json:"products"`
    }
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)

    assert.Equal(t, 20, response.Total)
    assert.Equal(t, 10, response.Limit)
    assert.Equal(t, 5, response.Offset)
    assert.Len(t, response.Products, 2)
    assert.Equal(t, "Product 6", response.Products[0].Name)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
