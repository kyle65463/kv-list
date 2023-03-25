package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyle65463/kv-list/controllers"
)

func main() {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")

	// Define routes
	apiV1.GET("/pages", controllers.GetPage)
	apiV1.GET("/heads", controllers.GetHead)

	// Handle 404 errors
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 not found",
		})
	})

	// Start server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
