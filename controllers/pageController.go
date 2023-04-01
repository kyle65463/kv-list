package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyle65463/kv-list/database"
)

type PageResponse struct {
	ID          *int   `json:"id"`
	Data        []byte `json:"data"`
	NextPageKey *int   `json:"nextPageKey"`
}

func GetPage(c *gin.Context) {
	// Parse request
	id := c.Param("id")

	// Acquire db connection
	conn, err := database.DbPool.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to acquire connection"),
		})
		return
	}
	defer conn.Release()

	// Read the page from db
	var page PageResponse
	err = conn.QueryRow(
		context.Background(),
		"SELECT id, data, next_page_key FROM pages WHERE id = $1",
		id,
	).Scan(&page.ID, &page.Data, &page.NextPageKey)
	if err != nil {
		// TODO: Handle no result error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to query database: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, page)
}
