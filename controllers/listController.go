package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetListHead(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Get Head",
	})
}
