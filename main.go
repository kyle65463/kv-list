package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kyle65463/kv-list/controllers"
	"github.com/kyle65463/kv-list/database"
)

func init() {
	// Initialize environment variables
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// Connect to the database
	database.CreateDbConnection()
}

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
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
