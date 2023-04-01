package main

import (
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

	// Establish database connection
	database.CreateDbConnection()
}

func cleanup() {
	// Close the database connection
	database.CloseDbConnection()
}

func main() {
	defer cleanup()

	// Define routes
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/pages/:key", controllers.GetPage(database.DbPool))
	apiV1.GET("/lists/:key", controllers.GetListHead(database.DbPool))
	apiV1.POST("/lists/:key", controllers.SetList(database.DbPool))

	// Start the server
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
