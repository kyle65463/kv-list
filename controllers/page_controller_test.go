package controllers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShouldGetPageWithoutNextPageKey(t *testing.T) {
	// Setup the test
	data := PageResponse{
		Key:         "957e0240-5e40-4241-8c8f-c27cf21bce0a",
		Data:        []byte("dGVzdF9kYXRh"),
		NextPageKey: nil,
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

func TestShouldGetPageWithNextPageKey(t *testing.T) {
	// Setup the test
	data := PageResponse{
		Key:         "957e0240-5e40-4241-8c8f-c27cf21bce0a",
		Data:        []byte("dGVzdF9kYXRh"),
		NextPageKey: stringPtr("c9fc1f0a-8b56-4167-a638-38d04841be1f"),
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

func TestShouldRespondBadRequestWhenNoPageFound(t *testing.T) {
	// Setup the test
	pageKey := "957e0240-5e40-4241-8c8f-c27cf21bce0a"
	mock, c, w := setupTest(t, gin.Params{{Key: "key", Value: pageKey}})

	// Mock the empty query
	rows := mock.NewRows([]string{"key", "data", "next_page_key"})
	mock.ExpectQuery("SELECT key, data, next_page_key FROM pages").WithArgs(pageKey).WillReturnRows(rows)

	// Call the GetPage function
	GetPage(mock)(c)

	// Check the response
	assertResponseStatus(w, http.StatusBadRequest, t)
}
