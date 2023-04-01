package controllers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShouldGetListHeadWithoutNextPageKey(t *testing.T) {
	// Setup the test
	data := ListHeadResponse{
		Key:         "bc8f16c3-521b-4d76-8d2f-61147454c430",
		NextPageKey: nil,
	}
	mock, c, w := setupTest(t, gin.Params{{Key: "key", Value: data.Key}})

	// Mock the query
	rows := mock.NewRows([]string{"key", "next_page_key"}).
		AddRow(data.Key, data.NextPageKey)
	mock.ExpectQuery("SELECT key, next_page_key FROM lists").WithArgs(data.Key).WillReturnRows(rows)

	// Call the GetPage function
	GetListHead(mock)(c)

	// Check the response
	assertResponseStatus(w, http.StatusOK, t)
	assertResponseBody(w, data, t)
}

func TestShouldGetListHeadWithNextPageKey(t *testing.T) {
	// Setup the test
	data := ListHeadResponse{
		Key:         "bc8f16c3-521b-4d76-8d2f-61147454c430",
		NextPageKey: stringPtr("94c7a446-149c-436a-ab6f-9aa67fe4d09d"),
	}
	mock, c, w := setupTest(t, gin.Params{{Key: "key", Value: data.Key}})

	// Mock the query
	rows := mock.NewRows([]string{"key", "next_page_key"}).
		AddRow(data.Key, data.NextPageKey)
	mock.ExpectQuery("SELECT key, next_page_key FROM lists").WithArgs(data.Key).WillReturnRows(rows)

	// Call the GetPage function
	GetListHead(mock)(c)

	// Check the response
	assertResponseStatus(w, http.StatusOK, t)
	assertResponseBody(w, data, t)
}

func TestShouldRespondBadRequestWhenNoListHeadFound(t *testing.T) {
	// Setup the test
	listKey := "bc8f16c3-521b-4d76-8d2f-61147454c430"
	mock, c, w := setupTest(t, gin.Params{{Key: "key", Value: listKey}})

	// Mock the empty query
	rows := mock.NewRows([]string{"key", "next_page_key"})
	mock.ExpectQuery("SELECT key, next_page_key FROM lists").WithArgs(listKey).WillReturnRows(rows)

	// Call the GetPage function
	GetListHead(mock)(c)

	// Check the response
	assertResponseStatus(w, http.StatusBadRequest, t)
}
