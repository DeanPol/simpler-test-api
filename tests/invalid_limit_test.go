package tests

import (
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
*/
func TestGetProducts_InvalidLimit(t *testing.T) {
    dbConn, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening mock db: %v", err)
    }
    defer dbConn.Close()

    db.DB = dbConn

    req, err := http.NewRequest("GET", "/products?limit=-1&offset=0", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetProducts)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusBadRequest, rr.Code)

    assert.Equal(t, "Invalid limit value\n", rr.Body.String())

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
