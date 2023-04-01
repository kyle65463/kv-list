package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/kyle65463/kv-list/database"
)

type ListHeadResponse struct {
	Key         *string `json:"key"`
	NextPageKey *string `json:"nextPageKey"`
}

func GetListHead(c *gin.Context) {
	// Parse request
	key := c.Param("key")

	// Acquire db connection
	conn, err := database.DbPool.Acquire(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Read the list head from db
	var page ListHeadResponse
	err = conn.QueryRow(
		context.Background(),
		"SELECT key, next_page_key FROM lists WHERE key = $1",
		key,
	).Scan(&page.Key, &page.NextPageKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No result found
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No result found",
			})
		} else {
			// Other errors
			c.Status(http.StatusInternalServerError)
		}
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
	key := c.Param("key")

	// Acquire db connection
	conn, err := database.DbPool.Acquire(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Begin the transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	// Insert pages from the last element of the list to the first element
	var nextPageKey *string = nil
	for _, data := range list.Data {
		pageKey := uuid.New() // Generate a random key for the new page
		err = tx.QueryRow(context.Background(), `
            INSERT INTO pages (key, data, next_page_key)
            VALUES ($1, $2, $3) RETURNING key
        `,
			pageKey, data, nextPageKey,
		).Scan(&nextPageKey)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	// Insert the list head
	_, err = tx.Exec(context.Background(), `
	        INSERT INTO lists (key, next_page_key)
	        VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET next_page_key = $2
	    `,
		key, nextPageKey,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Set the list successfully",
	})
}
