package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func errorInternalServer(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}