package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Get Page",
	})
}

func GetHead(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Get Head",
	})
}
