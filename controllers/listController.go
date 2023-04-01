package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyle65463/kv-list/database"
)

type ListHeadResponse struct {
	ID          *int   `json:"id"`
	NextPageKey *int   `json:"nextPageKey"`
}

func GetListHead(c *gin.Context) {
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

	// Read the list head from db
	var page ListHeadResponse
	err = conn.QueryRow(
		context.Background(),
		"SELECT id, next_page_key FROM lists WHERE id = $1",
		id,
	).Scan(&page.ID, &page.NextPageKey)
	if err != nil {
		// TODO: Handle no result error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to query database: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, page)
}
