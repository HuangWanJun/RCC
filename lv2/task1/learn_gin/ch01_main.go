package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func pong(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "43.34",
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run(":8083")
}
