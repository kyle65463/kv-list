package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kyle65463/kv-list/controllers"
	"github.com/kyle65463/kv-list/database"
)

var dbPool database.PgxInterface

func init() {
	// Initialize environment variables
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// Establish database connection
	dbPool = database.CreateDbConnection()
}

func main() {
	// Define routes
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/pages/:key", controllers.GetPage(dbPool))
	apiV1.DELETE("/pages", controllers.DeletePages(dbPool))
	apiV1.GET("/lists/:key", controllers.GetListHead(dbPool))
	apiV1.POST("/lists/:key", controllers.SetList(dbPool))

	// Start the server
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
