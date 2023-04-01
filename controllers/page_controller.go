package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyle65463/kv-list/database"
)

type PageResponse struct {
	Key         string  `json:"key"`
	Data        []byte  `json:"data"`
	NextPageKey *string `json:"nextPageKey"`
}

func GetPage(db database.PgxInterface) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse request
		key := c.Param("key")

		// Read the page from db
		rows, err := db.Query(context.Background(),
			"SELECT key, data, next_page_key FROM pages WHERE key = $1",
			key,
		)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Check if the result is empty
		if !rows.Next() {
			// No result found
			c.JSON(http.StatusBadRequest, gin.H{"error": "No result found"})
			return
		}

		// Parse the row
		var page PageResponse
		err = rows.Scan(&page.Key, &page.Data, &page.NextPageKey)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, page)
	}
}
