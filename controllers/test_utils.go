package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func assertResponseStatus(w *httptest.ResponseRecorder, data int, t *testing.T) {
	assert.Equal(t, w.Code, data)
}

func assertResponseBody(w *httptest.ResponseRecorder, data any, t *testing.T) {
	actual := w.Body.Bytes()
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	if !bytes.Equal(expected, actual) {
		t.Errorf("the expected json: %s is different from actual %s", expected, actual)
	}
}

func setupTest(t *testing.T, params []gin.Param) (pgxmock.PgxConnIface, *gin.Context, *httptest.ResponseRecorder) {
	// Create a mock pool and connection
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}

	// Set up the Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = params

	return mock, c, w
}

func stringPtr(s string) *string {
	return &s
}
