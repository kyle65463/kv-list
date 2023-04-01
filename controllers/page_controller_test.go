package controllers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShouldGetPage(t *testing.T) {
	// Setup the test
	data := PageResponse{
		Key: "test_key", Data: []byte("dGVzdF9kYXRh"), NextPageKey: nil,
	}
	mock, c, w := setupTest(t, gin.Params{{Key: "key", Value: data.Key}})

	// Mock the query
	rows := mock.NewRows([]string{"key", "data", "next_page_key"}).
		AddRow(data.Key, data.Data, data.NextPageKey)
	mock.ExpectQuery("SELECT key, data, next_page_key FROM pages").WithArgs(data.Key).WillReturnRows(rows)

	// Call the GetPage function
	GetPage(mock)(c)

	// Check the response
	assertResponseStatus(w, http.StatusOK, t)
	assertResponseBody(w, data, t)
}
