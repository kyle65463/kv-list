package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func parseInt(s string) *int {
	if i, err := strconv.Atoi(s); err == nil {
		return &i
	}
	return nil
}

func DeletePages(db database.PgxInterface) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse request
		rawInterval := c.Query("interval")
		limit := parseInt(c.Query("limit"))

		// Read the page from db
		interval, err := time.ParseDuration(rawInterval)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong interval format"})
			return
		}

		// Delete pages older than the interval
		timeThreshold := time.Now().Add(-interval).UTC()
		if limit != nil {
			// Delete with limit
			fmt.Println(*limit)
			_, err = db.Exec(context.Background(), `
				DELETE FROM pages WHERE ctid IN (
					SELECT ctid FROM pages WHERE created_at < $1 LIMIT $2
				)
			`,
				timeThreshold,
				*limit,
			)
		} else {
			// Delete without limit
			_, err = db.Exec(context.Background(),
				"DELETE FROM pages WHERE created_at < $1",
				timeThreshold,
			)
		}
		if err != nil {
			fmt.Print(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Delete outdated pages successfully",
		})
	}
}
