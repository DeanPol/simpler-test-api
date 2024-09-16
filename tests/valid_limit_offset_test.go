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
   Set up expectations for mock queries.
   Assert that the total count, limit, and offset are correct.
*/

func TestGetProducts_ValidLimitOffset(t *testing.T) {
    dbConn, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening mock db: %v", err)
    }
    defer dbConn.Close()

    db.DB = dbConn

    mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))
    mock.ExpectQuery("SELECT id, name, description, price, stock FROM products").
        WithArgs(2, 0).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock"}).
            AddRow(1, "Product 1", "Description 1", 100.0, 10).
            AddRow(2, "Product 2", "Description 2", 150.0, 5))

    req, err := http.NewRequest("GET", "/products?limit=2&offset=0", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetProducts)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    // Parse the response body
    var response struct {
        Total    int               `json:"total"`
        Limit    int               `json:"limit"`
        Offset   int               `json:"offset"`
        Products []models.Product `json:"products"`
    }
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)

    assert.Equal(t, 20, response.Total)
    assert.Equal(t, 2, response.Limit)
    assert.Equal(t, 0, response.Offset)
    assert.Len(t, response.Products, 2)
    assert.Equal(t, "Product 1", response.Products[0].Name)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

