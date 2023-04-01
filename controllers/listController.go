package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyle65463/kv-list/database"
)

type ListHeadResponse struct {
	ID          *int `json:"id"`
	NextPageKey *int `json:"nextPageKey"`
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

func SetList(c *gin.Context) {
	// Parse request
	type body struct {
		Data [][]byte `json:"data"`
	}
	var list body
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acquire db connection
	conn, err := database.DbPool.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to acquire connection"),
		})
		return
	}
	defer conn.Release()

	// Begin the transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to begin the transaction"),
		})
		return
	}
	defer tx.Rollback(context.Background())

	// Insert pages from the last element of the list to the first element
	var nextPageKey *int = nil
	for _, data := range list.Data {
		err = tx.QueryRow(context.Background(), `
            INSERT INTO pages (data, next_page_key)
            VALUES ($1, $2) RETURNING id
        `,
			data, nextPageKey,
		).Scan(&nextPageKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Insert the list head
	err = tx.QueryRow(context.Background(), `
	        INSERT INTO lists (next_page_key)
	        VALUES ($1) RETURNING id
	    `,
		nextPageKey,
	).Scan(nextPageKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	tx.Commit(context.Background())

	c.JSON(http.StatusOK, list)
}
