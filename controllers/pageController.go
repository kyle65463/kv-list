package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyle65463/kv-list/database"
)

func GetPage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Get Page",
	})
}

func SetPage(c *gin.Context) {
	// Parse request body
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
