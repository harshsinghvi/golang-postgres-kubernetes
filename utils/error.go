package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func InternalServerError(c *gin.Context, msg string, err error) {
	log.Printf("%s %v\n", msg, err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Something went wrong",
	})
}
