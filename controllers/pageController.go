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

	// Insert the page from db
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

func SetPage(c *gin.Context) {
	// Parse request
	type body struct {
		Data        *[]byte `json:"data"`
		NextPageKey *int    `json:"nextPageKey"`
	}
	var p body
	if err := c.ShouldBindJSON(&p); err != nil {
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

	// Insert the value into db
	_, err = conn.Exec(context.Background(), `
            INSERT INTO pages (data, next_page_key)
            VALUES ($1, $2)
        `,
		p.Data, p.NextPageKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func GetHead(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Get Head",
	})
}
