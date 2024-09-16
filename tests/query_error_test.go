package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"simpler-test-api/db"
	"simpler-test-api/handlers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

/*
   Create an HTTP request with an invalid limit (negative value)
   Check the response status code is 500 Internal Server Error
*/

func TestGetProducts_QueryError(t *testing.T) {
    dbConn, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening mock db: %v", err)
    }
    defer dbConn.Close()

    db.DB = dbConn

    // Mock the COUNT query response (total count of products)
    mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))

    mock.ExpectQuery("SELECT id, name, description, price, stock FROM products").
        WithArgs(10, 0). // Default limit and offset
        WillReturnError(sql.ErrConnDone)

    req, err := http.NewRequest("GET", "/products", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetProducts)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusInternalServerError, rr.Code)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
