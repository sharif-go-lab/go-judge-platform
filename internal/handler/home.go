package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler handles the home page
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Online Judge System",
	})
}
