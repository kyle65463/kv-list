package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/kyle65463/kv-list/database"
)

type PageResponse struct {
	Key         *string `json:"key"`
	Data        []byte  `json:"data"`
	NextPageKey *string `json:"nextPageKey"`
}

func GetPage(c *gin.Context) {
	// Parse request
	key := c.Param("key")

	// Acquire db connection
	conn, err := database.DbPool.Acquire(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Read the page from db
	var page PageResponse
	err = conn.QueryRow(
		context.Background(),
		"SELECT key, data, next_page_key FROM pages WHERE key = $1",
		key,
	).Scan(&page.Key, &page.Data, &page.NextPageKey)
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
